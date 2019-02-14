package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Generated by http://quicktype.io

type LogStructure struct {
	Hostname       string `json:"@hostname"`
	Timestamp      string `json:"@timestamp"`
	Vertical       string `json:"@vertical"`
	Level          string `json:"level"`
	Msg            string `json:"msg"`
	Retention      string `json:"retention"`
	ServiceName    string `json:"service-name"`
	ServiceVersion string `json:"service-version"`
	Type           string `json:"type"`
}

func TestLogWithFields(t *testing.T) {
	//new log with fields
	log := LogWithFields

	// a buffer to store the output from logrus
	buffer := bytes.NewBuffer(make([]byte, 0, 20))
	logrus.SetOutput(buffer)

	message := "test"
	log.Info(message)

	//unmarshal logrus message
	var logmsg LogStructure
	err := json.Unmarshal(buffer.Bytes(), &logmsg)
	if err != nil {
		t.Error("fail")
	}
	// if logmsg.Msg != message {
	//     t.Errorf("fail expected %s to be the same as %s", message, logmsg.Msg)
	//

	assert.NotNil(t, logmsg.Hostname)
	assert.NotNil(t, logmsg.Level)
	assert.NotNil(t, logmsg.Retention)
	assert.NotNil(t, logmsg.ServiceName)
	assert.NotNil(t, logmsg.ServiceVersion)
	assert.NotNil(t, logmsg.Timestamp)
	assert.NotNil(t, logmsg.Type)
	assert.NotNil(t, logmsg.Vertical)

	// assert for not nil (good when you expect something)
	if assert.NotNil(t, logmsg.Msg) {
		assert.Equal(t, message, logmsg.Msg)
	}

}

func TestLogError(t *testing.T) {
	// Init
	var logmsg LogStructure
	testError := errors.New("Test error")
	buffer := bytes.NewBuffer(make([]byte, 0, 20))
	logrus.SetOutput(buffer)

	// Do
	LogError(testError)
	err := json.Unmarshal(buffer.Bytes(), &logmsg)

	// Validate
	if err != nil {
		t.Error(err)
	}

	if assert.NotNil(t, logmsg.Msg) {
		assert.Equal(t, testError.Error(), logmsg.Msg)
	}

}
