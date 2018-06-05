package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"

	"golang.org/x/net/websocket"

	"regexp"
	"time"

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

func reports(wr http.ResponseWriter, req *http.Request) {
	log.LogWithFields.Debug(req.URL.RawPath)
	log.LogWithFields.Debug(req.URL.RawQuery)

	if file := req.URL.Query().Get("file"); file != "" {
		// fmt.Println(file)
		if match, _ := regexp.MatchString("^[a-z-0-9]+.json$", file); match == true {
			fileData, err := ioutil.ReadFile("data/" + file)
			if err != nil {
				log.LogWithFields.Error(err.Error())
				return
			}

			// unamrshal and marshal again because of a stupid witespace somewere
			var dat map[string]interface{}
			if err := json.Unmarshal(fileData, &dat); err != nil {
				panic(err)
			}

			data, err := json.Marshal(dat)
			if err != nil {
				log.LogWithFields.Error(err.Error())
				return
			}

			wr.Header().Set("Content-Type", "application/json")
			wr.Write(data)
			return
		}
	}

	files, err := ioutil.ReadDir("data")
	if err != nil {
		log.LogWithFields.Error(err.Error())
		return
	}

	type fileInfo struct {
		File string    `json:"file"`
		Time time.Time `json:"time"`
	}
	var filesInfo []fileInfo

	for _, file := range files {
		filesInfo = append(filesInfo, fileInfo{File: file.Name(), Time: file.ModTime()})
	}

	data, err := json.Marshal(filesInfo)
	if err != nil {
		log.LogWithFields.Error(err.Error())
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(data)
	// wr.WriteHeader(http.StatusOK)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))
	mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))
	mux.Handle("/ws", websocket.Handler(internal.WsServer))

	mux.HandleFunc("/probe", internal.URLStress)

	mux.HandleFunc("/api/v1/reports", reports)

	mux.HandleFunc("/.well-known/ready", healthHandler)
	mux.HandleFunc("/.well-known/live", healthHandler)
	mux.HandleFunc("/.well-known/metrics", prometheusHandler)

	// Register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	log.LogWithFields.Infof("Starting proxy server on: %v", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.LogWithFields.Fatalf("ListenAndServe: %v", err.Error())
	}
}
