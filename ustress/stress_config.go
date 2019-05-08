package ustress

import (
	// "bytes"
	"fmt"
	// "io"
	// "io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"errors"

	"golang.org/x/net/context"
)

// A message processes url and returns the result on responseChan.
// ctx is places in a struct, but this is ok to do.

var (
	ctx context.Context
)

const (
	// HTTPfolder - the folder where the reports will be dumped
	HTTPfolder = "./data/" // Would be nice to generate absolute path via code || gopath
)

var (
	InvalidOptionType = errors.New("Invalid option value has been given")
	InvalidOptionName = errors.New("Option with that name does not exist")
)

// OutputSaver is a callback with report data
type OutputSaver func(*Report, chan bool)

// Headers map
type Headers map[string]string

// MonkeyConfig structure
type StressConfig struct {
	// URL to probe
	URL    string `json:"url"`
	Method string `json:"method"`
	// Number of request to be sent to the victim
	Requests int `json:"requests"`

	// Ho many treads to be used (dependent on the image resources)
	Threads int `json:"threads"`
	// similar to curl --resolve Force resolve of HOST:PORT to ADDRESS
	Resolve string `json:"resolve"`
	// insecure
	Insecure bool `json:"insecure"`

	// payload
	Payload string `json:"payload"`

	// Headers
	Headers Headers `json:"headers"`

	Duration  int `json:"duration"` // Minutes
	Frequency int `json:"frequency"`// Miliseconds

	// client instantiate a new http client
	client *http.Client // `json:"-"`

	// If each worker should capture response
	WithResponse bool      `json:"withResponse"`
	StopCh       chan bool `json:"-"`
}

// ValidateConfig ...
func (cfg *StressConfig) ValidateConfig() error {
	_, err := url.ParseRequestURI(cfg.URL)
	if err != nil {
		return fmt.Errorf("param: URL is not a valid url")
	}
	if reflect.TypeOf(cfg.Requests).Kind() != reflect.Int {
		return fmt.Errorf("param: requests is of wrong type, must be int")
	}
	if cfg.Requests <= 0 && cfg.Duration <= 0 {
		return fmt.Errorf("param: requests <= 0")
	}
	if reflect.TypeOf(cfg.Threads).Kind() != reflect.Int {
		return fmt.Errorf("param: workers is of wrong type, must be int")
	}
	if cfg.Threads <= 0 {
		return fmt.Errorf("param: workers <= 0 ")
	}
	//if cfg.Requests < cfg.Threads {
	//	cfg.Threads = cfg.Requests
	//}
	return nil
}

type Option func(*StressConfig) error

func NewOption(name string, val interface{}) Option {
	return func(s *StressConfig) error {
		elem := reflect.ValueOf(s).Elem()
		field := elem.FieldByName(name)
		if field.IsValid() && field.CanSet() {
			valT := reflect.TypeOf(val).Kind().String()
			fieldT := field.Kind().String()
			if valT != fieldT {
				return InvalidOptionType
			}
			field.Set(reflect.ValueOf(val))
		} else {
			return InvalidOptionName
		}
		return nil
	}

}

func NewStressConfig(opts ...Option) (*StressConfig, error) {
	cfg := &StressConfig{}
	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return cfg, err
		}
	}
	if cfg.Requests == 0 {
		r := time.Duration(cfg.Frequency) * time.Millisecond
		d := time.Duration(cfg.Duration) * time.Second
		hits := d.Seconds() / r.Seconds() // Estimating chan buffer
		cfg.Requests = int(hits)+ 1
	}
	cfg.client = cfg.newHTTPClient()

	return cfg, cfg.ValidateConfig()
}
