package ustress

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	"github.com/google/uuid"
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
			rmConn(conn)
		}
	}
}

// WsServer handles the ws connections
func WsServer(ws *websocket.Conn) {
	addConn(ws)
	for {
		// var data map[string]interface{}
		monkeyConfig := &MonkeyConfig{} // TODO validate mkcfg

		err := websocket.JSON.Receive(ws, &monkeyConfig)
		if err != nil {
			fmt.Println(err.Error())
			if err == io.EOF {
				rmConn(ws)
				return
			}
		}

		err = monkeyConfig.ValidateConfig()
		if err != nil {
			fmt.Println(err)
		} else {
			// @todo, you can do better
			monkeyConfig = NewConfig(monkeyConfig.URL, monkeyConfig.Requests, monkeyConfig.Threads, monkeyConfig.Resolve, monkeyConfig.Insecure, monkeyConfig.Method, monkeyConfig.Payload, monkeyConfig.Headers)
			b, _ := monkeyConfig.NewWebsocketStressReport()
			writeAll(string(b))
		}

	}
}

// NewWebsocketStressReport will returin a report via websocket
// It is configured from the websocket handler
func (monkeyConfig *MonkeyConfig) NewWebsocketStressReport() ([]byte, error) {
	start := time.Now()
	report := Report{TimeStamp: time.Now(), UUID: uuid.New(), MonkeyConfig: *monkeyConfig}

	requests := make(chan WorkerConfig, monkeyConfig.Requests)
	response := make(chan WorkerConfig, monkeyConfig.Requests)

	drainRequests := make(chan bool) // true to drain requests queue
	quitWorkers := make(chan bool)   // true to kill go routines

	// start number of threads
	for w := 1; w <= monkeyConfig.Threads; w++ {
		go worker(w, quitWorkers, requests, response)
	}

	// send work to request channel
	log.LogWithFields.Infof("Work Accepted: %#v\n", monkeyConfig)
	go func() {
		for req := 1; req <= monkeyConfig.Requests; req++ {
			requests <- WorkerConfig{Request: req, monkeyConfig: monkeyConfig}
		}
		return
	}()

	// a go routine to update the report sent via websocket at a given interval
	// the goroutine it is canceled once the cancelUpdatechannel receives a true statement
	cancelUpdate := make(chan bool)
	go func() {
		for {
			// fmt.Printf("\n%d\n %+v\n\n", len(report.Workers), report)
			select {
			case <-cancelUpdate:
				runtime.Goexit()
				return
			default:
				time.Sleep(500 * time.Millisecond)
				// create a snapshot of the current report
				report.calcStats()
				report.Duration = time.Since(start).Seconds()
				b, err := json.Marshal(report)
				if err != nil {
					return
				}
				writeAll(string(b))
			}
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
		wrkConf := <-response
		report.Workers = append(report.Workers, wrkConf)
		if wrkConf.Status == 0 {
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
	// cancelUpdatethe update go rutine
	// from this point fw the websocket is updated from the socket handler
	cancelUpdate <- true
	// close(response)
	report.calcStats() //TODO - > return errors
	report.Duration = time.Since(start).Seconds()
	// marshal the report
	b, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

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
		tr.CloseIdleConnections()
	}()

	// save report to file
	fileWriter := newFile(fmt.Sprintf("%s.json", report.UUID))
	fmt.Fprintf(fileWriter, string(b))
	// return json report
	return b, nil
}
