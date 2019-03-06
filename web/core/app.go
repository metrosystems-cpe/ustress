package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"

	"io/ioutil"

	"time"

	reCassandra "git.metrosystems.net/reliability-engineering/reliability-incubator/reutils/cassandra"
	log "git.metrosystems.net/reliability-engineering/ustress/log"
	"github.com/gocql/gocql"
)

var session *gocql.Session

var NoDBConn = errors.New("No database connection")

// AppConfigPath default path
const AppConfigPath = "git.metrosystems.net/reliability-engineering/ustress/configuration.yaml"

// APIendpoint ...
type APIendpoint func(a *App, w http.ResponseWriter, r *http.Request) (interface{}, error)

type JSONResponse map[string]interface{}

// App will store app state, and other metadata alongside with utility functions
type App struct {
	Version       string
	Configuration *reCassandra.Config
	Session       *gocql.Session
}

// InitSession initializes a cassandra session and attaches to the app struct
func (a *App) InitSession() error {
	var err error
	cluster := gocql.NewCluster(a.Configuration.Hosts...)
	cluster.Port = a.Configuration.Port
	auth := gocql.PasswordAuthenticator{
		Username: a.Configuration.Username,
		Password: a.Configuration.Password,
	}

	// This will set the keyspace for all of our foregoing queries
	cluster.Keyspace = a.Configuration.Keyspace

	//Defaults
	cluster.Consistency = gocql.Quorum //Write on at least one node
	cluster.Authenticator = auth
	d, _ := time.ParseDuration("1m")
	cluster.Timeout = d

	a.Session, err = cluster.CreateSession()
	return err
}

// Init initializes app
func (a *App) Init() {
	log.LogWithFields.Info("Initializing cassandra session")
	err := a.InitSession()
	log.LogError(err)
	if err != nil {
		log.LogWithFields.Error("Failed to initialize cassandra connection")
		return
	}
	log.LogWithFields.Info("Creating cassandra schema")
	a.CreateSchema()
	// defer a.Session.Close()

}

// CreateSchema generates the required tables
func (a *App) CreateSchema() {
	var err error
	for tableName, tableCols := range Tables {
		var metaq string
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s", a.Configuration.Keyspace, tableName)
		cols := "( "
		for cName, cAttribute := range tableCols {
			if cName == "meta" {
				metaq = cAttribute
				continue

			}
			col := fmt.Sprintf("%s %s, ", cName, cAttribute)
			cols += col
		}
		cols += ") "
		query += cols
		query += metaq

		log.LogWithFields.Infof("Executing DB query: %s", query)
		err = a.Session.Query(query).Exec()
		log.LogError(err)

	}

}

// NewAppFromYAML inits the app from a yaml file
func NewAppFromYAML(configpath string) *App {

	var a App
	a.load(configpath)
	a.Init()
	return &a
}

// NewAppFromEnv Gets configuration from env
func NewAppFromEnv() (*App, error) {
	var a App
	var err error
	a.Configuration, err = reCassandra.ParseConnectionString()
	if err == nil {
		a.Init()
		return &a, nil
	}
	return nil, err
}

// NewApp inits the app
func NewApp(version string, c *reCassandra.Config) *App {
	a := &App{Version: version, Configuration: c}
	a.Init()
	return a
}

// Middleware sets custom headers, and writes json response
func Middleware(a *App, endpoint APIendpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if a.Session == nil {
			res := newResponse(nil, NoDBConn)
			writeResponse(w, res)
			return
		}

		d, e := endpoint(a, w, r)
		res := newResponse(d, e)

		writeResponse(w, res)

	}
}

// LocalCassandraConfig Used for development
func LocalCassandraConfig() *reCassandra.Config {
	return &reCassandra.Config{
		Hosts:    []string{"127.0.0.1"},
		Keyspace: "ustress",
		Port:     9042,
	}
}

func (a *App) load(configpath string) {
	yamlFile, err := ioutil.ReadFile(configpath)
	log.LogError(err)
	err = yaml.Unmarshal(yamlFile, a)
	log.LogError(err)
}

func writeResponse(w http.ResponseWriter, response JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if response["error"] != "" {
		w.WriteHeader(400)
		return
	}
	jsonBytes, err := json.Marshal(response)
	log.LogError(err)
	w.Write(jsonBytes)
}

func newResponse(data interface{}, err error) map[string]interface{} {
	res := map[string]interface{}{
		"entries": data,
		"error":   "",
	}
	if err != nil {
		res["error"] = err.Error()

	}
	return res
}
