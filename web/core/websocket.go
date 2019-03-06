package core

import (
	"fmt"
	"io"
	"sync"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	ustress "git.metrosystems.net/reliability-engineering/ustress/ustress"
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

func WriteAllWebsockets(report *ustress.Report) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range wsConns {
		err := websocket.JSON.Send(conn, string(report.JSON()))
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
			r, err := ustress.NewReport(monkeyConfig, WriteAllWebsockets, 500)
			log.LogError(err)
			s := NewStressTest(r)
			if app.Session != nil {
				err = s.Save(app.Session)

			}
			log.LogError(err)
		}
	}
}
