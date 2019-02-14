package core

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	ustress "git.metrosystems.net/reliability-engineering/ustress/ustress"
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

func WriteAllWebsockets(msg string) {
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

func InjectWsContext(app *App, fn func(app *App, ws *websocket.Conn)) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		fn(app, ws)
	}
}

// WsServer handles the ws connections
func WsServer(app *App, ws *websocket.Conn) {
	addConn(ws)
	for {
		// var data map[string]interface{}
		monkeyConfig := &ustress.MonkeyConfig{} // TODO validate mkcfg

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
			log.LogError(err)
		} else {
			// @todo, you can do better
			monkeyConfig = ustress.NewConfig(
				monkeyConfig.URL, monkeyConfig.Requests, monkeyConfig.Threads,
				monkeyConfig.Resolve, monkeyConfig.Insecure, monkeyConfig.Method,
				monkeyConfig.Payload, monkeyConfig.Headers, monkeyConfig.WithResponse)
			b, _ := NewWebsocketStressReport(app, monkeyConfig)
			WriteAllWebsockets(string(b))
		}

	}
}

// NewWebsocketStressReport will returin a report via websocket
// It is configured from the websocket handler
func NewWebsocketStressReport(a *App, monkeyConfig *ustress.MonkeyConfig) ([]byte, error) {
	fmt.Println(monkeyConfig)
	start := time.Now()
	report := ustress.Report{TimeStamp: time.Now(), UUID: uuid.New(), Config: monkeyConfig}

	requests := make(chan ustress.WorkerData, monkeyConfig.Requests)
	response := make(chan ustress.WorkerData, monkeyConfig.Requests)

	drainRequests := make(chan bool) // true to drain requests queue
	quitWorkers := make(chan bool)   // true to kill go routines

	// start number of threads
	for w := 1; w <= monkeyConfig.Threads; w++ {
		go ustress.Worker(w, quitWorkers, requests, response)
	}

	// send work to request channel
	log.LogWithFields.Infof("Work Accepted: %#v\n", monkeyConfig)
	go func() {
		for req := 1; req <= monkeyConfig.Requests; req++ {
			requests <- ustress.WorkerData{Request: req, MonkeyConfig: monkeyConfig}
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
				report.CalcStats()
				report.Duration = time.Since(start).Seconds()
				b := report.JSON()
				if len(b) == 0 {
					return
				}
				WriteAllWebsockets(string(b))
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
		if !monkeyConfig.WithResponse {
			wrkConf.ResponseBody = ""
		}
		report.Data = append(report.Data, wrkConf)
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
	report.CalcStats() //TODO - > return errors
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
		ustress.Tr.CloseIdleConnections()
	}()

	fileWriter := ustress.NewFile(fmt.Sprintf("%s.json", report.UUID))
	defer fileWriter.Close()

	if a.Session != nil {
		stressTestReport := StressTest{ID: report.UUID, Report: &report}
		err = stressTestReport.Save(a.Session)
	}

	jsonReport := report.JSON()
	fmt.Fprintf(fileWriter, string(jsonReport))

	return b, err
}
