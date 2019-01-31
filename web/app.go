package web

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"

	"io/ioutil"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	"github.com/gocql/gocql"
)

var session *gocql.Session

// AppConfigPath default path
const AppConfigPath = "git.metrosystems.net/reliability-engineering/ustress/configuration.yaml"

// APIendpoint ...
type APIendpoint func(a *App, w http.ResponseWriter, r *http.Request)

// Configuration will store app config
type Configuration struct {
	Cluster  []string // List of IPs
	Keyspace string
	Username string
	Password string
}

// App will store app state, and other metadata alongside with utility functions
type App struct {
	Version       string
	Configuration *Configuration
	Session       *gocql.Session
}

// InitSession initializes a cassandra session and attaches to the app struct
func (a *App) InitSession() error {
	var err error
	cluster := gocql.NewCluster(a.Configuration.Cluster...)
	auth := gocql.PasswordAuthenticator{
		Username: a.Configuration.Username,
		Password: a.Configuration.Password,
	}
	cluster.Keyspace = a.Configuration.Keyspace

	//Defaults
	cluster.Consistency = gocql.One //Write on at least one node
	cluster.Authenticator = auth

	a.Session, err = cluster.CreateSession()
	return err
}

// NewAppFromYAML inits the app from a yaml file
func NewAppFromYAML(configpath string) *App {
	var a App
	a.load(configpath)
	err := a.InitSession()
	log.LogError(err)
	return &a
}

// NewApp inits the app
func NewApp(version string, c *Configuration) *App {
	a := &App{Version: version, Configuration: c}
	err := a.InitSession()
	log.LogError(err)
	return a
}

// InjectContext will inject app state into each API handler
func InjectContext(a *App, endpoint APIendpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint(a, w, r)
	}
}

func (a *App) load(configpath string) {
	yamlFile, err := ioutil.ReadFile(configpath)
	log.LogError(err)
	err = yaml.Unmarshal(yamlFile, a)
	log.LogError(err)
}
