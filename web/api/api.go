package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.metrosystems.net/reliability-engineering/ustress/ustress"

	"github.com/google/uuid"

	"git.metrosystems.net/reliability-engineering/ustress/log"
	// ustress "git.metrosystems.net/reliability-engineering/ustress/ustress"
	"git.metrosystems.net/reliability-engineering/ustress/web/core"
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
func URLStress(a *core.App, wr http.ResponseWriter, req *http.Request) (interface{}, error) {
	// exampleCall := "?url=http://localhost:9090&requests=20&workers=4"
	// http://localhost:9090/probe?resolve=10.29.30.8:443&url=https://idam-pp.metrosystems.net/.well-known/openid-configuration&requests=10&workers=4

	urlPath := req.URL.Query()
	uParam = urlPath.Get("url")
	rParam, _ = strconv.Atoi(urlPath.Get("requests"))
	wParam, _ = strconv.Atoi(urlPath.Get("workers"))


	restMK, err := ustress.NewStressConfig(
		ustress.NewOption("URL", uParam),
		ustress.NewOption("Requests", rParam),
		ustress.NewOption("Threads", wParam),
		// ustress.NewOption("Resolve", resolve),
		//insecure,
		//method,
		//"",
		//nil,
		//false
	)
	if err != nil {
		return nil, err
	}
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
