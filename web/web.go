package web

import (
	"net/http"
	"net/http/pprof"
	"os"
	"path"

	"golang.org/x/net/websocket"

	api "git.metrosystems.net/reliability-engineering/ustress/web/api"
	"git.metrosystems.net/reliability-engineering/ustress/web/core"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	requestedAsset := path.Join(r.URL.Path, "index.html")
	_, err := os.Stat(requestedAsset)
	if err != nil {
		http.ServeFile(w, r, requestedAsset)
		return
	}

}

// MuxHandlers ...
func MuxHandlers(a *core.App) *http.ServeMux {

	// var addr = flag.String("addr", ":8080", "The addr of the application.")
	// flag.Parse()

	mux := http.NewServeMux()

	// redirect to ui
	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {

		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Pargma", "no-cache")
		writer.Header().Set("Expires", "0")
		// The redirect is cached by the browser, thus most of the endpoints endup with unwanted 301
		http.Redirect(writer, req, "/ustress", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/ustress", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Pargma", "no-cache")
		writer.Header().Set("Expires", "0")
		http.Redirect(writer, req, "/ustress/ui/public/", http.StatusMovedPermanently)
	})

	// Serving static files
	mux.Handle("/ustress/static/", http.StripPrefix("/ustress/", http.FileServer(http.Dir("web/ui/build"))))

	// Index file
	mux.Handle("/ustress/ui/public/", http.StripPrefix("/ustress/ui/public/", http.FileServer(http.Dir("web/ui/build/"))))

	mux.Handle("/ustress/data/", http.StripPrefix("/ustress/data/", http.FileServer(http.Dir("data"))))

	mux.Handle("/ustress/api/v1/ws", websocket.Handler(core.InjectWsContext(a, core.WsServer)))
	mux.HandleFunc("/ustress/api/v1/file_reports", core.FileReportsView)
	mux.HandleFunc("/ustress/api/v1/reports", core.Middleware(a, api.GetReports))

	mux.HandleFunc("/ustress/api/v1/probe", core.Middleware(a, api.URLStress))
	mux.HandleFunc("/ustress/api/v1/test", core.TestHandlerView)

	mux.HandleFunc("/.well-known/ready", core.HealthHandlerView)
	mux.HandleFunc("/.well-known/live", core.HealthHandlerView)
	mux.HandleFunc("/.well-known/metrics", core.PrometheusHandlerView)

	// Register pprof handlers
	mux.HandleFunc("/ustress/debug/pprof/", pprof.Index)
	mux.HandleFunc("/ustress/debug/pprof/cmdline/", pprof.Cmdline)
	mux.HandleFunc("/ustress/debug/pprof/profile/", pprof.Profile)
	mux.HandleFunc("/ustress/debug/pprof/symbol/", pprof.Symbol)
	mux.HandleFunc("/ustress/debug/pprof/trace/", pprof.Trace)

	return mux

	// log.LogWithFields.Infof("Starting proxy server on: %v", *addr)
	// if err := http.ListenAndServe(*addr, mux); err != nil {
	// 	log.LogWithFields.Fatalf("ListenAndServe: %v", err.Error())
	// }
}
