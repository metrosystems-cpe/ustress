package server

import (
	"io"
	"sync"

	log "git.metrosystems.net/reliability-engineering/rest-monkey/log"
	"git.metrosystems.net/reliability-engineering/rest-monkey/restmonkey"
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
		monkeyConfig := &restmonkey.MonkeyConfig{} // TODO validate mkcfg

		err := websocket.JSON.Receive(ws, &monkeyConfig)
		if err != nil {
			log.LogWithFields.Error(err.Error())
			if err == io.EOF {
				rmConn(ws)
				return
			}
		}

		err = monkeyConfig.ValidateConfig()
		if err != nil {
			log.LogWithFields.Error(err.Error())
		} else {
			// @todo, you can do better
			monkeyConfig = restmonkey.NewConfig(monkeyConfig.URL, monkeyConfig.Requests, monkeyConfig.Threads, monkeyConfig.Resolve, monkeyConfig.Insecure)

			report, err := restmonkey.NewReport(monkeyConfig)
			if err != nil {
				log.LogWithFields.Error(err.Error())
			}
			writeAll(string(report))
		}

	}
}
