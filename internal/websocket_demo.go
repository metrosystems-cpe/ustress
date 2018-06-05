package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	log "git.metrosystems.net/reliability-engineering/traffic-monkey/log"
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
		}
	}
}

// WsServer ...
func WsServer(ws *websocket.Conn) {
	addConn(ws)
	for {
		// var data map[string]interface{}
		data := MonkeyConfig{}

		err := websocket.JSON.Receive(ws, &data)
		if err != nil {
			fmt.Println(err.Error())
			if err == io.EOF {
				rmConn(ws)
				return
			}
		}
		fmt.Printf("Got message: %#v\n", data)
		newReportSimulation()

	}
}

func newReportSimulation() {
	start := time.Now()
	Report := Report{
		TimeStamp: time.Now(),
		UUID:      uuid.New(),
		Workers:   []*Worker{},
		MonkeyConfig: MonkeyConfig{
			Requests: 10,
		},
	}

	for i := 0; i < 100; i++ {
		worker := Worker{Request: i, Status: 200, Thread: 1, url: "https://idam.metrosystems.net/.well-known/openid-configuration", insecure: false, Duration: 0.002153429}
		Report.Workers = append(Report.Workers, &worker)
		fmt.Printf("append worker: %#v\n", worker)
		Report.calcStats()
		Report.Duration = time.Since(start).Seconds()
		b, err := json.Marshal(Report)
		if err != nil {
			log.LogWithFields.Error(err.Error())
		}

		writeAll(string(b))
		time.Sleep(200 * time.Millisecond)
	}

}
