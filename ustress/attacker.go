package ustress

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"

	"fmt"

	// "errors"
	"time"

	"github.com/google/uuid"
	// "time"
)


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

// Helper function to attach headers to request
func setHeaders(r *http.Request, h Headers) {
	if h == nil {
		return
	}
	for k, v := range h {
		r.Header.Set(k, v)
	}
}



// Worker func
// is function for the concurent goroutines
// request channel it is used to receive work
// response channel is used to send back the work done
// id is the routine id, named thread for easy uderstanding
func Worker(thread int, quitWorkers <-chan bool, request <-chan WorkerData, response chan<- WorkerData) {
	defer func() {
		if recover() != nil {
			runtime.Goexit()
		}
	}()
	for {
		select {
		case <-quitWorkers:
			runtime.Goexit()
			return
		case work := <-request:
			if work.StressConfig == nil {
				continue
			}
			var payload io.Reader
			start := time.Now()
			work.Thread = thread

			switch work.StressConfig.Method {
			case "GET", "DELETE", "OPTIONS":
				payload = nil
			default:
				payload = bytes.NewBuffer([]byte(work.StressConfig.Payload))

			}

			httpClient = work.StressConfig.client
			httpClient.Timeout = 3 * time.Second
			httpRequest, err := http.NewRequest(work.StressConfig.Method, work.StressConfig.URL, payload)
			if err != nil {
				return
			}
			setHeaders(httpRequest, work.StressConfig.Headers)

			httpResponse, error := httpClient.Do(httpRequest)
			if error != nil {
				work.Error = error.Error()
			} else {
				work.Status = httpResponse.StatusCode
				if work.StressConfig.WithResponse {
					bodyBytes, _ := ioutil.ReadAll(httpResponse.Body)
					work.ResponseBody = string(bodyBytes)
				}

				httpRequest.Close = true
				httpResponse.Body.Close()

			}
			work.Duration = time.Since(start).Seconds()
			response <- work
		}
	}
}

// goroutine
func attack(cfg *StressConfig, requests chan WorkerData, drainReq chan bool) {
	defer func() {
		if recover() != nil {
			runtime.Goexit()
		}
	}()

	for req := 1; req <= cfg.Requests; req++ {
		select {
		case <-drainReq:
			runtime.Goexit()

		case requests <- WorkerData{Request: req, StressConfig: cfg}:
			if req == cfg.Requests {
				runtime.Goexit()
			}
			if cfg.Frequency != 0 {
				time.Sleep(time.Duration(cfg.Frequency) * time.Millisecond)
			}
		case <- time.After(time.Duration(1) * time.Second):
			runtime.Goexit()
		}
	}
}


func closeChannels(ch ...chan interface{}) {
	for _, c := range ch {
		close(c)
	}


}

// goroutine
func streamOutput(r *Report, saveFunc OutputSaver, interval int, stopStreaming chan bool) {
	start := time.Now()
	for {
		select {
		case <-stopStreaming:
			runtime.Goexit()
		case <-time.After(time.Duration(interval) * time.Millisecond):
			r.CalcStats()
			r.Duration = time.Since(start).Seconds()
			saveFunc(r, stopStreaming)
		}
	}
}

func SaveFileReport(r *Report) {
	jsonReport := r.JSON()
	fileWriter := NewFile(fmt.Sprintf("%s.json", r.UUID))
	defer fileWriter.Close()
	fmt.Fprintf(fileWriter, string(jsonReport))
}



// NewReport probes an endpoint and generates a new report
func NewReport(cfg *StressConfig, saveFunc OutputSaver, tickerSave int) (*Report, error) {

	start := time.Now()

	report := &Report{TimeStamp: time.Now(), UUID: uuid.New(), Config: cfg, Completed: false}
	report.Stats.CodesCount = map[int]int{}

	requests := make(chan WorkerData, cfg.Requests)
	response := make(chan WorkerData, cfg.Requests)
	stop := make(chan bool, 1000)

	drainRequests := make(chan bool, 1000) // true to drain requests queue
	quitWorkers := make(chan bool, cfg.Threads)   // true to kill go routines
	numberOfErrors := 0

	// start number of threads
	for w := 1; w <= cfg.Threads; w++ {
		go Worker(w, quitWorkers, requests, response)
	}


	// Will distribute requests to workers
	go attack(cfg, requests, drainRequests)

	if tickerSave != 0 {
		go streamOutput(report, saveFunc, tickerSave, stop)
	}


	// Counting errors
loop:
	for {
		select {
		case data := <-response:

			report.Data = append(report.Data, data)
			report.Stats.CodesCount[data.Status] += 1

			if data.Status == 0 {
				numberOfErrors++
			}

			// If error rate is greater stop
			if data.Request == cfg.Requests*30/100 {
				if float64(numberOfErrors)/float64(data.Request)*100 == float64(100) {
					break loop
				}
			}
			if len(report.Data) == cfg.Requests {
				break loop

			}

			if cfg.Duration != 0 && time.Since(start).Seconds() >= float64(cfg.Duration) {
				break loop
			}
			// Mainly used for sending data through websocket

		case <-time.After(time.Duration(ChannelTimeout) * time.Second):
			break loop

		case <- stop:

			for i := 1; i <= cfg.Threads; i++ {
				quitWorkers <- true
			}
			defer close(requests)
			defer close(response)
			defer close(quitWorkers)
			defer close(drainRequests)
			Tr.CloseIdleConnections()
			return report, errors.New("Gracefully stopped")

		}
	}

	report.CalcStats()
	report.Duration = time.Since(start).Seconds()
	report.Completed = true
	go func() {
		for i := 1; i <= cfg.Threads; i++ {
			quitWorkers <- true
		}
		// close channels
		defer close(requests)
		defer close(response)
		defer close(quitWorkers)
		defer close(drainRequests)
		defer close(stop)
		// close transporter
		// @todo : on hit and run (repetitive request ) some transporter routines still
		Tr.CloseIdleConnections()
	}()

	if saveFunc == nil {
		SaveFileReport(report)
		return report, nil
	}

	saveFunc(report, drainRequests)

	return report, nil

}
