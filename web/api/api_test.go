package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.metrosystems.net/reliability-engineering/ustress/web/core"
)

func getApp() *core.App {
	return core.NewApp("0.0.1", core.LocalCassandraConfig())

}

func makeReq(req *http.Request, router *http.ServeMux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestURLStress(t *testing.T) {
	// TODO data folder should have a absolute path

	mux := http.NewServeMux()
	handler := core.Middleware(getApp(), URLStress)
	mux.HandleFunc("/ustress/api/v1/probe", handler)

	// Init
	params := "?url=http://localhost:8080/ustress/api/v1/test&requests=10&workers=2&method=GET"
	req, err := http.NewRequest("GET", "/ustress/api/v1/probe"+params, nil)

	assert.Nil(t, err)

	// Do
	rr := makeReq(req, mux)
	resbody := rr.Result()
	result, err := ioutil.ReadAll(resbody.Body)
	assert.Nil(t, err)
	fmt.Println(result)

	// Validate
	assert.Equal(t, 200, resbody.StatusCode)
	assert.Equal(t, "application/json", resbody.Header.Get("content-type"))

	var jr core.JSONResponse
	err = json.Unmarshal(result, &jr)
	entries := jr["entries"].(map[string]interface{})
	fmt.Println(entries)

	assert.Nil(t, err)
	assert.NotNil(t, entries)
	assert.NotNil(t, entries["timestamp"])
	assert.NotNil(t, entries["config"])
	assert.NotNil(t, entries["stats"])
	assert.NotNil(t, entries["data"])
	assert.NotNil(t, entries["uuid"])

	//TODO entries['config']['requests'] etc...
	// if assert.NotNil(t, r.Config.Method) {
	// 	assert.Equal(t, "GET", r.Config.Method)

	// }

	// if assert.NotNil(t, r.Config.Requests) {
	// 	assert.Equal(t, 10, r.Config.Requests)

	// }

	// if assert.NotNil(t, r.Config.Threads) {
	// 	assert.Equal(t, 2, r.Config.Threads)

	// }
	// Todo assert many more things
	// Check multiple cases such as incomplete params
	// Bad data type for the params etc...

}
