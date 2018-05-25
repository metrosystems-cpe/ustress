package main

import (
	"flag"
	"io"
	"net/http"

	"git.metrosystems.net/reliability-engineering/traffic-monkey/internal"
	log "git.metrosystems.net/reliability-engineering/traffic-monkey/log"
)

func healthHandler(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
	wr.Header().Set("Content-Type", "application/json")
	log.LogWithFields.Debug(req.URL.Path)
	io.WriteString(wr, `{"Status": OK}`)
}

func prometheusHandler(wr http.ResponseWriter, req *http.Request) {
	log.LogWithFields.Debug(req.URL.Path)
	wr.WriteHeader(http.StatusOK)
}

func main() {
	var addr = flag.String("addr", ":9090", "The addr of the application.")
	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))

	mux.HandleFunc("/probe", internal.URLStress)
	mux.HandleFunc("/.well-known/ready", healthHandler)
	mux.HandleFunc("/.well-known/live", healthHandler)
	mux.HandleFunc("/.well-known/metrics", prometheusHandler)

	log.LogWithFields.Infof("Starting proxy server on: %v", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.LogWithFields.Fatalf("ListenAndServe: %v", err.Error())
	}
}
