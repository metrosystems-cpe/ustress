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

// ValidateConfig ...
func (mkcfg *MonkeyConfig) ValidateConfig() error {
	_, err := url.ParseRequestURI(mkcfg.URL)
	if err != nil {
		return fmt.Errorf("param: URL is not a valid url")
	}
	if reflect.TypeOf(mkcfg.Requests).Kind() != reflect.Int {
		return fmt.Errorf("param: requests is of wrong type, must be int")
	}
	if mkcfg.Requests <= 0 {
		return fmt.Errorf("param: requests <= 0")
	}
	if reflect.TypeOf(mkcfg.Threads).Kind() != reflect.Int {
		return fmt.Errorf("param: workers is of wrong type, must be int")
	}
	if mkcfg.Threads <= 0 {
		return fmt.Errorf("param: workers <= 0 ")
	}
	if mkcfg.Requests < mkcfg.Threads {
		mkcfg.Threads = mkcfg.Requests
	}
	return nil
}

// WsServer handles the ws connections
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
			log.LogWithFields.Infof("Work Accepted: %#v\n", mkcfg)
			b, _ := mkcfg.NewWebsocketStressReport()
			writeAll(string(b))
		}

	}
}

// NewWebsocketStressReport will returin a report via websocket
// It is configured from the websocket handler
func (mkcfg *MonkeyConfig) NewWebsocketStressReport() ([]byte, error) {
	start := time.Now()
	report := Report{TimeStamp: time.Now(), UUID: uuid.New(), MonkeyConfig: *mkcfg}

	requests := make(chan WorkerConfig, mkcfg.Requests)
	response := make(chan WorkerConfig, mkcfg.Requests)
	// start number of threads
	for w := 1; w <= mkcfg.Threads; w++ {
		go work(w, requests, response)
	}

	// send work to request channel
	// fmt.Println(mkcfg.Requests)
	go func() {
		for req := 1; req <= mkcfg.Requests; req++ {
			requests <- WorkerConfig{Request: req, mkcfg: *mkcfg}
		}
		// close(requests) // daca inchid canalul apar mesaje in plus
		return
	}()

	// a go routine to update the report sent via websocket at a given interval
	// the goroutine it is canceled once the cancelUpdatechannel receives a true statement
	cancelUpdate := make(chan bool)
	cancelReport := make(chan bool)
	go func() {
		for {
			// fmt.Printf("\n%d\n %+v\n\n", len(report.Workers), report)
			select {
			case <-cancelUpdate:
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
			case <-cancelReport:
				for range requests {
					// drain request channel
				}
			default:
				// dono
			}
		}
	}()

	numberOfErrors := 0
	for res := 1; res <= mkcfg.Requests; res++ {
		wrkConf := <-response
		report.Workers = append(report.Workers, wrkConf)
		if wrkConf.Status == 0 {
			numberOfErrors++
		}
		// should be percentage 20%
		if res == mkcfg.Requests*30/100 {
			if float64(numberOfErrors)/float64(res)*100 == float64(100) {
				fmt.Printf("%d : %d : %v\n", numberOfErrors, res, float64(numberOfErrors)/float64(res)*100)
				cancelReport <- true
				break
			}
		}
	}
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
	// save report to file
	fileWriter := newFile(fmt.Sprintf("%s.json", report.UUID))
	fmt.Fprintf(fileWriter, string(b))
	// return json report
	return b, nil
}
