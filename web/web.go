package web

import (
	"net/http"
	"net/http/pprof"

	"golang.org/x/net/websocket"

	ustress "git.metrosystems.net/reliability-engineering/rest-monkey/ustress"
	api "git.metrosystems.net/reliability-engineering/rest-monkey/web/api"
)

// MuxHandlers ...
func MuxHandlers() *http.ServeMux {

	// var addr = flag.String("addr", ":8080", "The addr of the application.")
	// flag.Parse()

	mux := http.NewServeMux()

	// redirect to ui
	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		http.Redirect(writer, req, "/ustress/ui/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/ustress", func(writer http.ResponseWriter, req *http.Request) {
		http.Redirect(writer, req, "/ustress/ui/", http.StatusMovedPermanently)
	})

	mux.Handle("/ustress/ui/", http.StripPrefix("/ustress/ui/", http.FileServer(http.Dir("web/ui"))))
	mux.Handle("/ustress/data/", http.StripPrefix("/ustress/data/", http.FileServer(http.Dir("data"))))

	mux.Handle("/ustress/api/v1/ws", websocket.Handler(ustress.WsServer))
	mux.HandleFunc("/ustress/api/v1/reports", reports)

	mux.HandleFunc("/ustress/api/v1/probe", api.URLStress)
	mux.HandleFunc("/ustress/api/v1/test", testHandler)

	mux.HandleFunc("/.well-known/ready", healthHandler)
	mux.HandleFunc("/.well-known/live", healthHandler)
	mux.HandleFunc("/.well-known/metrics", prometheusHandler)

	// Register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux

	// log.LogWithFields.Infof("Starting proxy server on: %v", *addr)
	// if err := http.ListenAndServe(*addr, mux); err != nil {
	// 	log.LogWithFields.Fatalf("ListenAndServe: %v", err.Error())
	// }
}
