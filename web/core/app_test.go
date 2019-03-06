package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	reCassandra "git.metrosystems.net/reliability-engineering/reliability-incubator/reutils/cassandra"
	"github.com/stretchr/testify/assert"
)

func newConfig(ip string) *reCassandra.Config {
	return &reCassandra.Config{
		Hosts:    []string{ip},
		Port:     9042,
		Keyspace: "ustress",
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
	a := NewApp("0.0.1", Config)

	// Do
	err := a.InitSession()

	// Validate
	assert.NotNil(t, err)
	if err == nil {
		t.Error("Init session should have returned an error")

	}
	Config = newConfig("127.0.0.1")
	a = NewApp("0.0.1", Config)

	err = a.InitSession()
	assert.Nil(t, err)

}
func TestMiddleware(t *testing.T) {

	Config := newConfig("127.0.0.1")
	a := NewApp("0.0.1", Config)
	mux := http.NewServeMux()

	handler := func(a *App, w http.ResponseWriter, r *http.Request) (interface{}, error) {
		assert.NotNil(t, a)
		if assert.NotNil(t, a.Version) {
			assert.Equal(t, "0.0.1", a.Version)
		}
		return nil, nil
	}
	mux.HandleFunc("/", Middleware(a, handler))

	req, _ := http.NewRequest("GET", "/", nil)
	makeReq(req, mux)

}
