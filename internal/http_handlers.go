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
	// http://localhost:9090/probe?resolve=10.29.30.8:443&url=https://idam-pp.metrosystems.net/.well-known/openid-configuration&requests=10&workers=4

	urlPath := req.URL.Query()

	if uParam = urlPath.Get("url"); uParam == "" {
		wr.WriteHeader(http.StatusBadRequest)
		// wr.Write([]byte("missing url parameter\n eg: %s%s", urlPath, exampleCall))
		wr.Write([]byte("missing url parameter"))
		return
	}

	insecure, _ := strconv.ParseBool(urlPath.Get("insecure"))
	// fmt.Println(insecure)

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

	resolve := urlPath.Get("resolve") // @todo validate'it for god sake

	// @todo handle error
	mk := MonkeyConfig{
		URL:      uParam,
		Requests: rParam,
		Threads:  wParam,
		Resolve:  resolve,
		Insecure: insecure,
	}

	log.Printf("%+v", mk)

	messages, _ := mk.NewURLStressReport()
	// os.Stdout.Write(messages)

	fmt.Println(string(messages))
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(messages)
}
