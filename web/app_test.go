package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Config = Configuration{
	Cluster:  []string{"127.0.0.1"},
	Keyspace: "test",
	Username: "cassandra",
	Password: "cassandra",
}

func makeReq(req *http.Request, router *http.ServeMux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestInitSession(t *testing.T) {
	// Init
	Config.Cluster = []string{"127.0.0.2"}
	a := NewApp("0.0.1", &Config)

	// Do
	err := a.InitSession()

	// Validate
	assert.NotNil(t, err)
	if err == nil {
		t.Error("Init session should have returned an error")

	}
	// Redo
	Config.Cluster = []string{"127.0.0.1"}
	a = NewApp("0.0.1", &Config)

	err = a.InitSession()
	assert.Nil(t, err)

}

func TestNewAppFromYAML(t *testing.T) {
	configpath := fmt.Sprintf("%s/%s/%s", os.Getenv("GOPATH"), "src", AppConfigPath)
	a := NewAppFromYAML(configpath)

	assert.NotNil(t, a.Version)
	assert.NotNil(t, a.Configuration)
	assert.NotNil(t, a.Configuration.Cluster)
	assert.NotNil(t, a.Configuration.Keyspace)
	assert.NotNil(t, a.Configuration.Username)
	assert.NotNil(t, a.Configuration.Password)
	assert.NotNil(t, a.Session)
}

func TestInjectContext(t *testing.T) {

	a := NewApp("0.0.1", &Config)
	mux := http.NewServeMux()

	handler := func(a *App, w http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, a)
		if assert.NotNil(t, a.Version) {
			assert.Equal(t, "0.0.1", a.Version)
			assert.Equal(t, "test", a.Configuration.Keyspace)

		}
	}

	mux.HandleFunc("/", InjectContext(a, handler))

	req, _ := http.NewRequest("GET", "/", nil)
	makeReq(req, mux)

}
