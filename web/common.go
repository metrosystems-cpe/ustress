package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
)

// common ? probably

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

func testHandler(wr http.ResponseWriter, req *http.Request) {
	time.Sleep(250 * time.Millisecond)
	switch req.Method {
	case "GET":
		wr.WriteHeader(http.StatusOK)
	case "POST":
		wr.WriteHeader(http.StatusCreated)
	case "PUT":
		wr.WriteHeader(http.StatusAccepted)
	case "DELETE":
		wr.WriteHeader(http.StatusNoContent)
	}
}

func reports(wr http.ResponseWriter, req *http.Request) {
	log.LogWithFields.Debug(req.URL.RawPath)
	log.LogWithFields.Debug(req.URL.RawQuery)

	if file := req.URL.Query().Get("file"); file != "" {
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
