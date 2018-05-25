package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/montanaflynn/stats"
	"golang.org/x/net/context"
)

// A message processes url and returns the result on responseChan.
// ctx is places in a struct, but this is ok to do.

var (
	ctx context.Context
)

const (
	// HTTPfolder - the folder where the reports will be dumped
	HTTPfolder = "./data/"
)

// the structure that will be passed to channels
type message struct {
	responseChan chan<- *message
	worker       Worker
	ctx          context.Context
}

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

// Worker details, needed for returning the output and build the report
type Worker struct {
	Request  int     `json:"request"`
	Status   int     `json:"status"` // json:"status,omitempty"
	Thread   int     `json:"thread"`
	url      string  // should use net.url
	resolve  string  // ip:port
	insecure bool    // insecure request, does not check the certificate
	Duration float64 `json:"duration"`
	Error    string  `json:"error"` //`json:"error,omitempty"`
}

// Report is the report structure, object
// @todo calculate percentile 99, 95, 75, 50
type Report struct {
	// id uuid
	// timestamp
	URL       string    `json:"url"`
	Requests  int       `json:"requests"`
	Resolve   string    `json:"resolve"`
	TimeStamp time.Time `json:"timestamp"`
	UUID      uuid.UUID `json:"uuid"`
	Stats     struct {
		Median          float64 `json:"median"`
		PercentileA     float64 `json:"50_percentile"`
		PercentileB     float64 `json:"75_percentile"`
		PercentileC     float64 `json:"95_percentile"`
		PercentileD     float64 `json:"99_percentile"`
		ErrorPercentage float64 `json:"error_percentage"`
	} `json:"stats"`

	Duration float64   `json:"durationTotal"`
	Workers  []*Worker `json:"data"`
}

func processMessages(id int, work <-chan *message) {
	for job := range work {
		select {
		// If the context is finished, don't bother processing the
		// message.
		case <-job.ctx.Done():
			continue
		default:
		}

		job.worker.doWork(id)

		select {
		case <-job.ctx.Done():
		case job.responseChan <- job:
		}
	}
}

// doWork method for the worker
func (wrk *Worker) doWork(id int) *Worker {
	// every new worker has a new http client.

	start := time.Now()
	wrk.Thread = id

	// curl -v -k --resolve "idam-pp.metrosystems.net:443:10.29.30.8"  'https://idam-pp.metrosystems.net:443/.well-known/openid-configuration' --insecure

	tr := &http.Transport{}

	// insecure request
	if wrk.insecure {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	dialer := &net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 0 * time.Second,
		DualStack: true,
	}

	log.Println(wrk.resolve)
	if wrk.resolve != "" {
		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, wrk.resolve)
		}
	}

	client := http.Client{
		Timeout:   time.Duration(2 * time.Second),
		Transport: tr,
	}

	//httpResponse, error := client.Head(url)
	httpResponse, error := client.Get(wrk.url)
	// log.Println(httpResponse.Header, error)
	if error != nil {
		// fmt.Println(error.Error())
		wrk.Error = error.Error()
		// wrk.Status = 000 // status in case of timeout
	} else {
		wrk.Status = httpResponse.StatusCode
	}

	wrk.Duration = time.Since(start).Seconds()
	// log.Printf("Worker Reporting: %+v", *wrk)
	return wrk
}

func newRequest(ctx context.Context, worker Worker, q chan<- *message, report *Report) {
	r := make(chan *message)
	select {
	// If the context finishes before we can send msg onto q,
	// exit early
	case <-ctx.Done():
		fmt.Println("Context ended before q could see message")
		return
	case q <- &message{
		responseChan: r,
		worker:       worker,
		// We are placing a context in a struct.  This is ok since it
		// is only stored as a passed message and we want q to know
		// when it can discard this message
		ctx: ctx,
	}:
	}

	select {
	case out := <-r:
		// fmt.Printf("%v\n", out)
		report.Workers = append(report.Workers, &out.worker)
	// If the context finishes before we could get the result, exit early
	case <-ctx.Done():
		fmt.Println("Context ended before q could process message")
	}
}

func (rep *Report) calcStats() *Report {
	var requestDurations []float64
	var numberOfErrors int
	var err error
	for _, value := range rep.Workers {
		// ignore errors
		if value.Status != 0 {
			requestDurations = append(requestDurations, value.Duration)
		} else {
			numberOfErrors++
		}
	}
	log.Printf("%-v", numberOfErrors)
	log.Printf("%-v", requestDurations)
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
	rep.Stats.ErrorPercentage = float64((numberOfErrors / rep.Requests) * 100)
	// log.Printf("%-v", rep)
	return rep
}

// NewURLStressReport probes an endpoint and generates a new report
func (mk *MonkeyConfig) NewURLStressReport() ([]byte, error) {
	// @toDo refactor vars
	url := mk.URL
	requests := mk.Requests
	threads := mk.Threads

	start := time.Now()
	report := Report{URL: url, Resolve: mk.Resolve, TimeStamp: time.Now(), UUID: uuid.New(), Requests: mk.Requests}

	q := make(chan *message, threads)
	// start number of threads
	for i := 1; i <= threads; i++ {
		go processMessages(i, q)
	}

	// send requests to q
	for k := 1; k <= requests; k++ {
		ctx := context.Background()
		wrk := Worker{url: url, Request: k, resolve: mk.Resolve, insecure: mk.Insecure} // fuck ! all this logic is stupid
		newRequest(ctx, wrk, q, &report)
	}
	close(q)
	report.calcStats()
	report.Duration = time.Since(start).Seconds()
	b, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	fileWriter := NewFile(fmt.Sprintf("%s.json", report.UUID))
	fmt.Fprintf(fileWriter, string(b))
	return b, nil
}

// NewFile returns a new file to write data to
func NewFile(filename string) *os.File {

	CreateDirIfNotExist(HTTPfolder)
	f, err := os.Create(HTTPfolder + filename)
	f, err = os.OpenFile(HTTPfolder+filename, os.O_RDWR|os.O_APPEND, 0766) // For read access.
	if err != nil {
		fmt.Println(err.Error())
	}

	return f
}

// CreateDirIfNotExist the function name says it all
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
