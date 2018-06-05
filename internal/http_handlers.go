package internal

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "git.metrosystems.net/reliability-engineering/rest-monkey/log"
	"git.metrosystems.net/reliability-engineering/rest-monkey/slackNotifier"
)

var (
	uParam         string
	rParam, wParam int
)

type returnError struct {
	RequestedURL string
	Error        string
	ExampleCall  string
	OptionalArgs []string
}

func errorHandler(wr http.ResponseWriter, req *http.Request, err string) {
	wr.WriteHeader(http.StatusBadRequest)
	returnError := returnError{
		RequestedURL: req.URL.String(),
		ExampleCall:  req.URL.String() + "?url=http://localhost:9090&requests=20&workers=4",
		OptionalArgs: []string{"resolve=IP:PORT", "insecure=true"},
		Error:        err,
	}
	log.LogWithFields.Debugf("%+v", returnError)
	z, _ := json.Marshal(returnError)
	wr.Write(z)
}

// URLStress is the handler for monkey probe
func URLStress(wr http.ResponseWriter, req *http.Request) {
	// exampleCall := "?url=http://localhost:9090&requests=20&workers=4"
	// http://localhost:9090/probe?resolve=10.29.30.8:443&url=https://idam-pp.metrosystems.net/.well-known/openid-configuration&requests=10&workers=4
	wr.Header().Set("Content-Type", "application/json")

	urlPath := req.URL.Query()

	if uParam = urlPath.Get("url"); uParam == "" {
		errorHandler(wr, req, "missing url parameter")
		return
	}

	insecure, _ := strconv.ParseBool(urlPath.Get("insecure"))
	// fmt.Println(insecure)

	rParam, _ = strconv.Atoi(urlPath.Get("requests"))
	if rParam <= 0 {
		errorHandler(wr, req, "missing requests parameter")
		return
	}

	wParam, _ = strconv.Atoi(urlPath.Get("workers"))
	if wParam <= 0 {
		errorHandler(wr, req, "missing workers parameter")
		return
	}

	resolve := urlPath.Get("resolve") // @todo validate ip:port

	// limit the number of requests and number of threads.
	if rParam > 1000 {
		rParam = 1000
	}
	if wParam > 20 {
		wParam = 20
	}

	mk := MonkeyConfig{
		URL:      uParam,
		Requests: rParam,
		Threads:  wParam,
		Resolve:  resolve,
		Insecure: insecure,
	}

	log.LogWithFields.Debugf("%+v", mk)

	messages, _ := mk.NewURLStressReport()
	log.LogWithFields.Debugf(string(messages))
	slackNotifier.DeliverReport(
		slackNotifier.RawParams{
			Link:       mk.URL,
			NrRequests: mk.Requests,
			NrThreads:  mk.Threads,
			Result:     messages,
		})
	wr.Write(messages)
}
