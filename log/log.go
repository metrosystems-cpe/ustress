package log

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

var hostName string

// LogWithFields is a logrus WithFields preset as aggreed by peng
var LogWithFields = logrus.WithFields(*logrusFieldsConfig())

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap:        logrus.FieldMap{logrus.FieldKeyTime: "@timestamp"},
		TimestampFormat: time.RFC3339Nano})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

func logrusFieldsConfig() *logrus.Fields {
	if hostName == "" {
		hostName, _ = os.Hostname()
	}
	logrusFields := logrus.Fields{
		"@hostname":       hostName,
		"@vertical":       "reliability",
		"service-name":    "rest-monkey",
		"service-version": "x",
		"type":            "service",
		"retention":       "technical",
	}
	return &logrusFields
}

func LogError(e error) {
	if e != nil {
		LogWithFields.Error(e)
	}
}
