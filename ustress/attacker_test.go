package ustress

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)


func getCfg() (*StressConfig, error) {
	return NewStressConfig(
		NewOption("URL", "http://localhost:8080/ustress/api/v1/test"),
		NewOption("Requests", 1),
		NewOption("Threads", 1),
	)

}



func TestWorker(t *testing.T) {
	var wg sync.WaitGroup
	request, response := make(chan WorkerData, 1), make(chan WorkerData, 1)
	cfg, err := getCfg()
	wg.Add(1)
	go Worker(1, request, response, &wg)
	request <- WorkerData{StressConfig:cfg, Request:1}
	res := <- response

	assert.Nil(t, err)
	assert.Equal(t, 1, res.Thread)
	assert.Equal(t, 1, res.Request)
	assert.Equal(t, 200, res.Status)

}

func TestAttack(t *testing.T) {

	cfg, _ := getCfg()
	ch := Attack(cfg)
	res := <-ch

	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Thread)
	assert.Equal(t, 1, res.Request)
	assert.Equal(t, 200, res.Status)


}

func TestNewReport(t *testing.T) {
	cfg, _ := getCfg()
	r, e := NewReport(cfg, nil, 0)
	assert.Nil(t, e)
	assert.Equal(t, len(r.Data), 1)
	assert.NotNil(t, r.Stats)

}



