package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	ustress "git.metrosystems.net/reliability-engineering/ustress/ustress"
	"git.metrosystems.net/reliability-engineering/ustress/web/core"
)

var RequiredParamsMissing = errors.New("Some of the required parameters are missing")

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
func URLStress(a *core.App, wr http.ResponseWriter, req *http.Request) (interface{}, error) {
	// exampleCall := "?url=http://localhost:9090&requests=20&workers=4"
	// http://localhost:9090/probe?resolve=10.29.30.8:443&url=https://idam-pp.metrosystems.net/.well-known/openid-configuration&requests=10&workers=4
	// TODO
	// Isolate validations
	// Extract method, payload, headers
	// More dynamicity
	wr.Header().Set("Content-Type", "application/json")

	urlPath := req.URL.Query()

	if uParam = urlPath.Get("url"); uParam == "" {
		errorHandler(wr, req, "missing url parameter")
		return nil, RequiredParamsMissing
	}

	insecure, _ := strconv.ParseBool(urlPath.Get("insecure"))

	rParam, _ = strconv.Atoi(urlPath.Get("requests"))
	if rParam <= 0 {
		errorHandler(wr, req, "missing requests parameter")
		return nil, RequiredParamsMissing
	}

	wParam, _ = strconv.Atoi(urlPath.Get("workers"))
	if wParam <= 0 {
		errorHandler(wr, req, "missing workers parameter")
		return nil, RequiredParamsMissing
	}

	resolve := urlPath.Get("resolve") // @todo validate ip:port

	method := urlPath.Get("method") // @todo validate ip:port

	// limit the number of requests and number of threads.
	if rParam > 1000 {
		rParam = 1000
	}
	if wParam > 20 {
		wParam = 20
	}

	restMK := ustress.NewConfig(uParam, rParam, wParam, resolve, insecure, method, "", nil, false)
	report, err := ustress.NewReport(restMK, nil, 0)
	if err != nil {
		log.LogWithFields.Error(err.Error())
	}

	stressTestReport := core.StressTest{ID: report.UUID, Report: report}
	err = stressTestReport.Save(a.Session)
	log.LogError(err)
	return report, err

}

func GetReports(a *core.App, wr http.ResponseWriter, req *http.Request) (interface{}, error) {
	r := core.StressTest{}
	urlPath := req.URL.Query()
	if uParam = urlPath.Get("id"); uParam == "" {
		return r.All(a.Session)
	}
	u, _ := uuid.Parse(urlPath.Get("id"))
	r.ID = u
	err := r.Get(a.Session)
	return r, err

}
