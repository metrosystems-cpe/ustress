package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"net"
	"net/http"
	"os"
	"time"

	log "git.metrosystems.net/reliability-engineering/rest-monkey/log"

	"github.com/google/uuid"

	"github.com/montanaflynn/stats"
	"golang.org/x/net/context"
)

// A message processes url and returns the result on responseChan.
// ctx is places in a struct, but this is ok to do.

var (
	ctx context.Context
	tr  = &http.Transport{
		MaxIdleConns:        30, // this should be set as the number of go routines
		MaxIdleConnsPerHost: 30,
	}
	client = &http.Client{}
)

const (
	// HTTPfolder - the folder where the reports will be dumped
	HTTPfolder = "./data/"
)

// MonkeyConfig structure
type MonkeyConfig struct {
	// URL to probe
	URL string
	// Number of request to be sent to the victim
	Requests int
	// Ho many treads to be used (dependent on the image resources)
	Threads int
	// similar to curl --resolve Force resolve of HOST:PORT to ADDRESS
	Resolve string
	// insecure
	Insecure bool
}

// WorkerConfig structure is used to track worker work
type WorkerConfig struct {
	Request  int          `json:"request"`
	Status   int          `json:"status"` // json:"status,omitempty"
	Thread   int          `json:"thread"`
	Duration float64      `json:"duration"`
	Error    string       `json:"error"` //`json:"error,omitempty"`
	mkcfg    MonkeyConfig // monkey cfg
}

// Report is the report structure, object
// @todo calculate percentile 99, 95, 75, 50
type Report struct {
	UUID         uuid.UUID    `json:"uuid"`
	TimeStamp    time.Time    `json:"timestamp"`
	MonkeyConfig MonkeyConfig `json:"config"`
	Stats        struct {
		Median          float64 `json:"median"`
		PercentileA     float64 `json:"50_percentile"`
		PercentileB     float64 `json:"75_percentile"`
		PercentileC     float64 `json:"95_percentile"`
		PercentileD     float64 `json:"99_percentile"`
		ErrorPercentage float64 `json:"error_percentage"`
	} `json:"stats"`

	Duration float64        `json:"durationTotal"`
	Workers  []WorkerConfig `json:"data"`
}

// function for the concurent goroutines
// request channel it is used to receive work
// response channel is used to send back the work done
// id is the routine id, named thread for easy uderstanding
func work(thread int, request <-chan WorkerConfig, response chan<- WorkerConfig) {
	for {
		start := time.Now()
		wrk := <-request
		// fmt.Printf("%+v\n", wrk)
		wrk.Thread = thread

		// insecure request
		if wrk.mkcfg.Insecure {
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}

		dialer := &net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 0 * time.Second,
			DualStack: true,
		}

		// resolve ip
		if wrk.mkcfg.Resolve != "" {
			tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, wrk.mkcfg.Resolve)
			}
		}

		client = &http.Client{
			Timeout:   time.Duration(2 * time.Second),
			Transport: tr,
		}
		// fmt.Println(wrk.mkcfg.URL)
		httpResponse, error := client.Get(wrk.mkcfg.URL)
		// log.Println(httpResponse.Header, error)
		if error != nil {
			wrk.Error = error.Error()
		} else {
			wrk.Status = httpResponse.StatusCode
		}
		wrk.Duration = time.Since(start).Seconds()
		log.LogWithFields.Debugf("Worker Reporting: %+v", wrk) // @todo add report uuid to worker
		// fmt.Printf("%+v\n", *wrk)
		response <- wrk
	}
}

func (rep *Report) calcStats() *Report {
	// fmt.Printf("%-v", rep)
	var requestDurations []float64
	var numberOfErrors int
	// var errorPercantege float64
	var err error
	for _, value := range rep.Workers {
		// ignore http codes 100s to 500s
		if value.Status > 100 && value.Status < 600 {
			requestDurations = append(requestDurations, value.Duration)
		} else {
			numberOfErrors++
		}
	}
	if rep.Stats.PercentileA, err = stats.Percentile(requestDurations, 50); err != nil {
		rep.Stats.PercentileA = 0
	}
	if rep.Stats.PercentileB, _ = stats.Percentile(requestDurations, 75); err != nil {
		rep.Stats.PercentileB = 0
	}
	if rep.Stats.PercentileC, _ = stats.Percentile(requestDurations, 95); err != nil {
		rep.Stats.PercentileC = 0
	}
	if rep.Stats.PercentileD, _ = stats.Percentile(requestDurations, 99); err != nil {
		rep.Stats.PercentileD = 0
	}
	if rep.Stats.Median, _ = stats.Median(requestDurations); err != nil {
		rep.Stats.Median = 0
	}

	rep.Stats.ErrorPercentage = float64(numberOfErrors) / float64(rep.MonkeyConfig.Requests) * 100
	log.LogWithFields.Debugf("%-v", rep)
	return rep
}

// NewRESTStressReport probes an endpoint and generates a new report
func (mkcfg *MonkeyConfig) NewRESTStressReport() ([]byte, error) {
	start := time.Now()
	report := Report{TimeStamp: time.Now(), UUID: uuid.New(), MonkeyConfig: *mkcfg}

	requests := make(chan WorkerConfig, mkcfg.Requests)
	response := make(chan WorkerConfig, mkcfg.Requests)
	// start number of threads
	for w := 1; w <= mkcfg.Threads; w++ {
		go work(w, requests, response)
	}

	// send requests to q
	go func() {
		for req := 1; req <= mkcfg.Requests; req++ {
			wrk := WorkerConfig{Request: req, mkcfg: *mkcfg}
			requests <- wrk
		}
		// close(requests)
	}()

	for res := 1; res <= mkcfg.Requests; res++ {
		report.Workers = append(report.Workers, <-response)
	}

	report.calcStats()
	report.Duration = time.Since(start).Seconds()
	b, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	fileWriter := newFile(fmt.Sprintf("%s.json", report.UUID))
	fmt.Fprintf(fileWriter, string(b))
	return b, nil
}

// NewFile returns a new file to write data to
func newFile(filename string) *os.File {

	createDirIfNotExist(HTTPfolder)
	f, err := os.Create(HTTPfolder + filename)
	f, err = os.OpenFile(HTTPfolder+filename, os.O_RDWR|os.O_APPEND, 0766) // For read access.
	if err != nil {
		log.LogWithFields.Errorln(err.Error())
	}

	return f
}

// CreateDirIfNotExist the function name says it all
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.LogWithFields.Errorln(err.Error())
		}
	}
}
