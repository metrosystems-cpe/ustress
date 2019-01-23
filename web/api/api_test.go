package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.metrosystems.net/reliability-engineering/ustress/web"
)

func makeReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app := web.MuxHandlers()
	app.ServeHTTP(rr, req)
	return rr
}

func TestURLStress(t *testing.T) {
	params := "?url=/ustress/api/v1/test&requests=10&workers=2"
	req, _ := http.NewRequest("GET", "/ustress/api/v1/probe"+params, nil)
	makeReq(req)

}
