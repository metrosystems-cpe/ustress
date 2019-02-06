package ustress

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime"

	"time"

	"github.com/google/uuid"
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

type OutputSaver func(Report)

type Headers map[string]string

// MonkeyConfig structure
type MonkeyConfig struct {
	// URL to probe
	URL    string
	Method string
	// Number of request to be sent to the victim
	Requests int
	// Ho many treads to be used (dependent on the image resources)
	Threads int
	// similar to curl --resolve Force resolve of HOST:PORT to ADDRESS
	Resolve string
	// insecure
	Insecure bool

	// payload
	Payload string

	// Headers
	Headers Headers

	// client instantiate a new http client
	client *http.Client // `json:"-"`

	// If each worker should capture response
	WithResponse bool `json:"withResponse"`
}

// WorkerData structure is used to track worker work
type WorkerData struct {
	Request      int           `json:"request"`
	Status       int           `json:"status"` // json:"status,omitempty"
	Thread       int           `json:"thread"`
	Duration     float64       `json:"duration"`
	Error        string        `json:"error"` //`json:"error,omitempty"`
	ResponseBody string        `json:"response"`
	MonkeyConfig *MonkeyConfig `json:"-"` // monkey cfg
}

// Helper function to attach headers to request
func setHeaders(r *http.Request, h Headers) {
	if h == nil {
		return
	}
	for k, v := range h {
		r.Header.Set(k, v)
	}
}

// function for the concurent goroutines
// request channel it is used to receive work
// response channel is used to send back the work done
// id is the routine id, named thread for easy uderstanding
func Worker(thread int, quitWorkers <-chan bool, request <-chan WorkerData, response chan<- WorkerData) {
	for {
		select {
		case <-quitWorkers:
			runtime.Goexit()
			return
		case work := <-request:
			var payload io.Reader
			start := time.Now()
			work.Thread = thread

			switch work.MonkeyConfig.Method {
			case "GET", "DELETE", "OPTIONS":
				payload = nil
			default:
				payload = bytes.NewBuffer([]byte(work.MonkeyConfig.Payload))

			}

			httpClient = work.MonkeyConfig.client
			httpClient.Timeout = 3 * time.Second
			httpRequest, err := http.NewRequest(work.MonkeyConfig.Method, work.MonkeyConfig.URL, payload)
			if err != nil {
				return
			}
			setHeaders(httpRequest, work.MonkeyConfig.Headers)

			httpResponse, error := httpClient.Do(httpRequest)
			if error != nil {
				work.Error = error.Error()
			} else {
				work.Status = httpResponse.StatusCode
				defer httpResponse.Body.Close()

				bodyBytes, _ := ioutil.ReadAll(httpResponse.Body)
				work.ResponseBody = string(bodyBytes)

				httpRequest.Close = true

			}
			work.Duration = time.Since(start).Seconds()
			response <- work
		}
	}
}

// ValidateConfig ...
func (monkeyConfig *MonkeyConfig) ValidateConfig() error {
	_, err := url.ParseRequestURI(monkeyConfig.URL)
	if err != nil {
		return fmt.Errorf("param: URL is not a valid url")
	}
	if reflect.TypeOf(monkeyConfig.Requests).Kind() != reflect.Int {
		return fmt.Errorf("param: requests is of wrong type, must be int")
	}
	if monkeyConfig.Requests <= 0 {
		return fmt.Errorf("param: requests <= 0")
	}
	if reflect.TypeOf(monkeyConfig.Threads).Kind() != reflect.Int {
		return fmt.Errorf("param: workers is of wrong type, must be int")
	}
	if monkeyConfig.Threads <= 0 {
		return fmt.Errorf("param: workers <= 0 ")
	}
	if monkeyConfig.Requests < monkeyConfig.Threads {
		monkeyConfig.Threads = monkeyConfig.Requests
	}
	return nil
}

// NewConfig ...
func NewConfig(
	url string, requests int, threads int,
	resolve string, insecure bool, method string,
	payload string, headers Headers, withResponse bool) *MonkeyConfig {
	monkeyConfig := &MonkeyConfig{
		URL:          url,
		Requests:     requests,
		Threads:      threads,
		Resolve:      resolve,
		Insecure:     insecure,
		Method:       method,
		Payload:      payload,
		Headers:      headers,
		WithResponse: withResponse,
	}
	monkeyConfig.client = monkeyConfig.newHTTPClient()
	return monkeyConfig
}

// NewReport probes an endpoint and generates a new report
func NewReport(monkeyConfig *MonkeyConfig) (Report, error) {

	start := time.Now()
	report := Report{TimeStamp: time.Now(), UUID: uuid.New(), Config: monkeyConfig}

	requests := make(chan WorkerData, monkeyConfig.Requests)
	response := make(chan WorkerData, monkeyConfig.Requests)

	drainRequests := make(chan bool) // true to drain requests queue
	quitWorkers := make(chan bool)   // true to kill go routines

	// start number of threads
	for w := 1; w <= monkeyConfig.Threads; w++ {
		go Worker(w, quitWorkers, requests, response)
	}

	// send requests to q
	go func() {
		for req := 1; req <= monkeyConfig.Requests; req++ {
			wrk := WorkerData{Request: req, MonkeyConfig: monkeyConfig}
			requests <- wrk
		}
	}()

	go func() {
		for {
			select {
			case <-drainRequests:
				for range requests {
					// drain request channel
				}
				runtime.Goexit()
				return
			default:
				// dono
			}
		}
	}()

	numberOfErrors := 0
	for res := 1; res <= monkeyConfig.Requests; res++ {
		data := <-response
		if !monkeyConfig.WithResponse {
			data.ResponseBody = ""
		}
		report.Data = append(report.Data, data)

		if data.Status == 0 {
			numberOfErrors++
		}

		if res == monkeyConfig.Requests*30/100 {
			if float64(numberOfErrors)/float64(res)*100 == float64(100) {
				break
			}
		}
	}
	// send either way to kill the go routine
	drainRequests <- true

	report.CalcStats()
	report.Duration = time.Since(start).Seconds()

	// outputSaver(report)

	go func() {
		for i := 1; i <= monkeyConfig.Threads; i++ {
			quitWorkers <- true
		}
		// close channels
		close(requests)
		close(response)
		close(quitWorkers)
		close(drainRequests)
		// close transporter
		// @todo : on hit and run (repetitive request ) some transporter routines still
		Tr.CloseIdleConnections()
	}()

	fileWriter := NewFile(fmt.Sprintf("%s.json", report.UUID))
	defer fileWriter.Close()
	jsonReport := report.JSON()
	fmt.Fprintf(fileWriter, string(jsonReport))

	return report, nil

}
