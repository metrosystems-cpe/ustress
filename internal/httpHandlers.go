package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	uParam         string
	rParam, wParam int
)

// URLStress write description
func URLStress(wr http.ResponseWriter, req *http.Request) {
	// exampleCall := "?url=http://localhost:9090&requests=20&workers=4"

	urlPath := req.URL.Query()

	if uParam = urlPath.Get("url"); uParam == "" {
		wr.WriteHeader(http.StatusBadRequest)
		// wr.Write([]byte("missing url parameter\n eg: %s%s", urlPath, exampleCall))
		wr.Write([]byte("missing url parameter"))
		return
	}

	rParam, _ = strconv.Atoi(urlPath.Get("requests"))
	if rParam <= 0 {
		wr.WriteHeader(http.StatusBadRequest)
		wr.Write([]byte("missing nr of requests parameter"))
		return
	}

	wParam, _ = strconv.Atoi(urlPath.Get("workers"))
	if wParam <= 0 {
		wr.WriteHeader(http.StatusBadRequest)
		wr.Write([]byte("missing nr of workers parameter"))
		return
	}

	log.Println(uParam, rParam, wParam)

	// @todo handle error
	messages, _ := NewURLStressReport(uParam, rParam, wParam)
	// os.Stdout.Write(messages)

	fmt.Println(string(messages))
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(messages)
}
