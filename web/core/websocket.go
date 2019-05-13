package core

import (
	"io"
	"runtime"
	"sync"

	"git.metrosystems.net/reliability-engineering/ustress/log"
	"git.metrosystems.net/reliability-engineering/ustress/ustress"
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

func WriteAllWebsockets(report *ustress.Report) error {
	for conn := range wsConns {
		err := websocket.JSON.Send(conn, string(report.JSON()))
		if err != nil {
			log.LogWithFields.Error(err.Error())
			rmConn(conn)
			return err
		}
	}
	return nil
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
		cfg := &ustress.StressConfig{} // TODO validate mkcfg

		err := websocket.JSON.Receive(ws, &cfg)
		if err != nil {
			if err == io.EOF {
				rmConn(ws)
				return
			}
		}

		cfg, err = ustress.NewStressConfig(
			ustress.NewOption("URL", cfg.URL),
			ustress.NewOption("Requests", cfg.Requests),
			ustress.NewOption("Threads", cfg.Threads),
			ustress.NewOption("Resolve", cfg.Resolve),
			ustress.NewOption("Insecure", cfg.Insecure),
			ustress.NewOption("Method", cfg.Method),
			ustress.NewOption("Payload", cfg.Payload),
			ustress.NewOption("Headers", cfg.Headers),
			ustress.NewOption("WithResponse", cfg.WithResponse),
			ustress.NewOption("Duration", cfg.Duration),
			ustress.NewOption("Frequency", cfg.Frequency),
		)
		log.LogError(err)

		r, err := ustress.NewReport(cfg, WriteAllWebsockets, 500)
		log.LogError(err)
		if err != nil {
			runtime.Goexit()
		}
		s := NewStressTest(r)
		if app.Session != nil {
			err = s.Save(app.Session)

		}
		log.LogError(err)
	}
}
