package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.metrosystems.net/reliability-engineering/ustress/web"
	"git.metrosystems.net/reliability-engineering/ustress/web/core"
	"github.com/stretchr/testify/assert"
)

func newConfig(ip string) *core.Configuration {
	return &core.Configuration{
		Cluster:  []string{ip},
		Keyspace: "test",
		Username: "cassandra",
		Password: "cassandra",
	}
}

func makeReq(req *http.Request, router *http.ServeMux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestInitSession(t *testing.T) {
	// Init
	Config := newConfig("127.0.0.2")
	a := web.NewApp("0.0.1", Config)

	// Do
	err := a.InitSession()

	// Validate
	assert.NotNil(t, err)
	if err == nil {
		t.Error("Init session should have returned an error")

	}
	Config = newConfig("127.0.0.1")
	a = web.NewApp("0.0.1", Config)

	err = a.InitSession()
	assert.Nil(t, err)

}

func TestNewAppFromYAML(t *testing.T) {
	configpath := fmt.Sprintf("%s/%s/%s", os.Getenv("GOPATH"), "src", core.AppConfigPath)
	a := web.NewAppFromYAML(configpath)

	assert.NotNil(t, a.Version)
	assert.NotNil(t, a.Configuration)
	assert.NotNil(t, a.Configuration.Cluster)
	assert.NotNil(t, a.Configuration.Keyspace)
	assert.NotNil(t, a.Configuration.Username)
	assert.NotNil(t, a.Configuration.Password)
	assert.NotNil(t, a.Session)
}

func TestMiddleware(t *testing.T) {

	Config := newConfig("127.0.0.1")
	a := web.NewApp("0.0.1", Config)
	mux := http.NewServeMux()

	handler := func(a *core.App, w http.ResponseWriter, r *http.Request) (interface{}, error) {
		assert.NotNil(t, a)
		if assert.NotNil(t, a.Version) {
			assert.Equal(t, "0.0.1", a.Version)
			assert.Equal(t, "test", a.Configuration.Keyspace)

		}
		return nil, nil
	}
	mux.HandleFunc("/", core.Middleware(a, handler))

	req, _ := http.NewRequest("GET", "/", nil)
	makeReq(req, mux)

}
