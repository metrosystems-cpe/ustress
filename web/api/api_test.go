package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	ustress "git.metrosystems.net/reliability-engineering/ustress/ustress"
	"git.metrosystems.net/reliability-engineering/ustress/web/api"
)

func makeReq(req *http.Request, router *http.ServeMux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestURLStress(t *testing.T) {
	// TODO data folder should have a absolute path
	mux := http.NewServeMux()
	handler := http.HandlerFunc(api.URLStress)
	mux.HandleFunc("/ustress/api/v1/probe", handler)

	// Init
	params := "?url=http://localhost:8080/ustress/api/v1/test&requests=10&workers=2&method=GET"
	req, err := http.NewRequest("GET", "/ustress/api/v1/probe"+params, nil)

	assert.Nil(t, err)

	// Do
	rr := makeReq(req, mux)
	resbody := rr.Result()
	result, _ := ioutil.ReadAll(resbody.Body)

	// Validate
	assert.Equal(t, 200, resbody.StatusCode)
	assert.Equal(t, "application/json", resbody.Header.Get("Content-Type"))

	var r ustress.Report
	err = json.Unmarshal(result, &r)

	assert.Nil(t, err)

	assert.NotNil(t, r.UUID)
	assert.NotNil(t, r.Duration)
	assert.NotNil(t, r.TimeStamp)
	assert.NotNil(t, r.Stats)
	assert.NotNil(t, r.Data)
	assert.NotNil(t, r.Config)

	if assert.NotNil(t, r.Config.Method) {
		assert.Equal(t, "GET", r.Config.Method)

	}

	if assert.NotNil(t, r.Config.Requests) {
		assert.Equal(t, 10, r.Config.Requests)

	}

	if assert.NotNil(t, r.Config.Threads) {
		assert.Equal(t, 2, r.Config.Threads)

	}
	// Todo assert many more things
	// Check multiple cases such as incomplete params
	// Bad data type for the params etc...

}
