package ustress

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	// "errors"
	"sync"
	"time"

	"github.com/google/uuid"
	// "time"
)

// ChannelTimeout is a last resort method of killing goroutines
const ChannelTimeout = 5

// WorkerData structure is used to track worker work
type WorkerData struct {
	Request      int           `json:"request"`
	Status       int           `json:"status"` // json:"status,omitempty"
	Thread       int           `json:"thread"`
	Duration     float64       `json:"duration"`
	Error        string        `json:"error"` //`json:"error,omitempty"`
	ResponseBody string        `json:"response"`
	StressConfig *StressConfig `json:"-"` // monkey cfg
}

// Worker is a goroutine that receives from request channel configuration
// for probing URL and sends its response via response chan
func Worker(
	thread int,
	request <-chan WorkerData,
	response chan<- WorkerData,
	group *sync.WaitGroup,
) {
	defer group.Done()
	defer func() {
		if recover() != nil {
			return
		}
	}()

	for work := range request {
		if work.StressConfig == nil {
			continue
		}
		work.Thread = thread
		stress(&work)
		response <- work
	}
}

// Attack main API for probing URL
func Attack(cfg *StressConfig) chan WorkerData {
	var group sync.WaitGroup

	requests := make(chan WorkerData, cfg.Requests)
	response := make(chan WorkerData, cfg.Requests)

	for w := 1; w <= cfg.Threads; w++ {
		group.Add(1)
		go Worker(w, requests, response, &group)
	}

	go func() {
		defer group.Wait()
		group.Add(1)
		go attack(cfg, requests, &group)
	}()

	return response

}

// NewReport probes an endpoint and generates a new report
// saveFunc will be called every [tickerSave]ms meant for handling the report in real time
func NewReport(cfg *StressConfig, saveFunc OutputSaver, tickerSave int) (*Report, error) {
	var wg sync.WaitGroup
	start := time.Now()
	report := &Report{TimeStamp: time.Now(), UUID: uuid.New(), Config: cfg, Completed: false}
	report.Stats.CodesCount = map[int]int{}
	quit := make(chan bool, 100)
	numberOfErrors := 0

	if tickerSave != 0 {
		wg.Add(1)
		go streamOutput(report, saveFunc, tickerSave, quit, &wg)
	}

	response := Attack(cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()

	loop:
		for {
			select {
			case data := <-response:

				report.Data = append(report.Data, data)
				report.Stats.CodesCount[data.Status]++

				if data.Status == 0 {
					numberOfErrors++
				}

				if len(report.Data) >= cfg.Requests {
					quit <- true
					break loop
				}

				// If error rate is greater stop
				if data.Request == cfg.Requests*30/100 {
					if float64(numberOfErrors)/float64(data.Request)*100 == float64(100) {
						quit <- true
						break loop
					}
				}

			case <-time.After(time.Duration(ChannelTimeout) * time.Second):
				quit <- true
				break loop

			case <-quit:
				break loop
			}
		}

	}()

	wg.Wait()

	report.CalcStats()
	report.Duration = time.Since(start).Seconds()
	report.Completed = true

	if saveFunc == nil {
		SaveFileReport(report)
		return report, nil
	}

	saveFunc(report)

	go func() {
		defer close(response)
		defer cfg.StopAttack()
		defer close(quit)
		Tr.CloseIdleConnections()
	}()

	return report, nil

}

// Goroutine that will stream output
func streamOutput(
	r *Report,
	saveFunc OutputSaver,
	interval int,
	quit chan bool,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	start := time.Now()
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Duration(interval) * time.Millisecond):
			r.CalcStats()
			r.Duration = time.Since(start).Seconds()
			err := saveFunc(r)
			if err != nil {
				break loop
			}
		}
	}
}

// Goroutine that will distribute the requests to workers
func attack(cfg *StressConfig, requests chan WorkerData, group *sync.WaitGroup) {
	defer group.Done()

	start := time.Now()
loop:
	for req := 1; req <= cfg.Requests; req++ {

		if cfg.Duration != 0 && time.Since(start).Seconds() >= float64(cfg.Duration) {
			break loop
		}

		select {
		case requests <- WorkerData{Request: req, StressConfig: cfg}:
			if cfg.Frequency != 0 {
				time.Sleep(time.Duration(cfg.Frequency) * time.Millisecond)
			}
		case <-time.After(time.Duration(1) * time.Second):
			break loop
		case <-cfg.stopCh:
			break loop
		}
	}
}

// Used for hitting a endpoint
func stress(wd *WorkerData) {

	var payload io.Reader

	start := time.Now()
	switch wd.StressConfig.Method {
	case "GET", "DELETE", "OPTIONS":
		payload = nil
	default:
		payload = bytes.NewBuffer([]byte(wd.StressConfig.Payload))

	}

	httpClient = wd.StressConfig.client
	httpClient.Timeout = 3 * time.Second
	httpRequest, err := http.NewRequest(wd.StressConfig.Method, wd.StressConfig.URL, payload)
	if err != nil {
		return
	}
	setHeaders(httpRequest, wd.StressConfig.Headers)

	httpResponse, error := httpClient.Do(httpRequest)

	if error != nil {
		wd.Error = error.Error()
	} else {
		defer httpResponse.Body.Close()
		wd.Status = httpResponse.StatusCode
		if wd.StressConfig.WithResponse {
			bodyBytes, _ := ioutil.ReadAll(httpResponse.Body)
			wd.ResponseBody = string(bodyBytes)
		}

		httpRequest.Close = true
	}

	wd.Duration = time.Since(start).Seconds()
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
