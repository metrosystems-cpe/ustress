package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sync"
	"time"

	log "git.metrosystems.net/reliability-engineering/rest-monkey/log"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

var (
	msgCounter int
	wsConns    = make(map[*websocket.Conn]interface{})
	mutex      = sync.Mutex{}
)

func addConn(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	wsConns[conn] = nil
	log.LogWithFields.Info("client connected")
}

func rmConn(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(wsConns, conn)
	log.LogWithFields.Info("client disconnected")
}

func writeAll(msg string) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range wsConns {
		err := websocket.JSON.Send(conn, msg)
		// err := websocket.Message.Send(conn, t)
		if err != nil {
			log.LogWithFields.Error(err.Error())
		}
	}
}

// ValidateConfig ...
func (mk *MonkeyConfig) ValidateConfig() error {
	_, err := url.ParseRequestURI(mk.URL)
	if err != nil {
		return fmt.Errorf("param: URL is not a valid url")
	}
	if reflect.TypeOf(mk.Requests).Kind() != reflect.Int {
		return fmt.Errorf("param: requests is of wrong type, must be int")
	}
	if mk.Requests <= 0 {
		return fmt.Errorf("param: requests <= 0")
	}
	if reflect.TypeOf(mk.Threads).Kind() != reflect.Int {
		return fmt.Errorf("param: workers is of wrong type, must be int")
	}
	if mk.Threads <= 0 {
		return fmt.Errorf("param: workers <= 0 ")
	}
	if mk.Requests < mk.Threads {
		mk.Threads = mk.Requests
	}
	return nil
}

// WsServer ...
func WsServer(ws *websocket.Conn) {
	addConn(ws)
	for {
		// var data map[string]interface{}
		mkcfg := MonkeyConfig{} // TODO validate mkcfg

		err := websocket.JSON.Receive(ws, &mkcfg)
		if err != nil {
			fmt.Println(err.Error())
			if err == io.EOF {
				rmConn(ws)
				return
			}
		}

		fmt.Printf("Got message: %#v\n", mkcfg)
		err = mkcfg.ValidateConfig()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("work accepted")
			fmt.Printf("Got message: %#v\n", mkcfg)
		}
		// b, _ := mkcfg.NewWebsocketStressReport()
		// writeAll(string(b))
	}
}

// NewWebsocketStressReport ...
func (mk *MonkeyConfig) NewWebsocketStressReport() ([]byte, error) {
	start := time.Now()
	report := Report{TimeStamp: time.Now(), UUID: uuid.New(), MonkeyConfig: *mk}

	q := make(chan *message, mk.Threads)
	// start number of threads
	for i := 1; i <= mk.Threads; i++ {
		go processMessages(i, q)
	}

	// send requests to q
	for k := 1; k <= mk.Requests; k++ {
		ctx := context.Background()
		worker := Worker{Request: k, mkcfg: mk}

		r := make(chan *message)
		select {
		case <-ctx.Done():
			fmt.Println("Context ended before q could see message")
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
			// this is CPU cycles consuming because it calculates stats after each worker finishes the work but it is cool :)
			report.calcStats()
			report.Duration = time.Since(start).Seconds()
			b, err := json.Marshal(report)
			if err != nil {
				return nil, err
			}
			writeAll(string(b))
		// If the context finishes before we could get the result, exit early
		case <-ctx.Done():
			fmt.Println("Context ended before q could process message")
		}

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

// func newReportSimulation() {
// 	start := time.Now()
// 	Report := Report{
// 		TimeStamp: time.Now(),
// 		UUID:      uuid.New(),
// 		Workers:   []*Worker{},
// 		MonkeyConfig: MonkeyConfig{
// 			Requests: 10,
// 		},
// 	}

// 	for i := 0; i < 100; i++ {
// 		worker := Worker{Request: i, Status: 200, Thread: 1, Duration: 0.002153429, Error: ""}
// 		Report.Workers = append(Report.Workers, &worker)
// 		fmt.Printf("append worker: %#v\n", worker)
// 		Report.calcStats()
// 		Report.Duration = time.Since(start).Seconds()
// 		b, err := json.Marshal(Report)
// 		if err != nil {
// 			log.LogWithFields.Error(err.Error())
// 		}

// 		writeAll(string(b))
// 		time.Sleep(200 * time.Millisecond)
// 	}

// }
