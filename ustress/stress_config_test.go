package ustress

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


var cfg *StressConfig
var err error


func TestOptions(t *testing.T) {
	cfg, err = NewStressConfig(NewOption("Threads", "2"))

	assert.NotNil(t, err)

	cfg, err = NewStressConfig(NewOption("Threads", 2), NewOption("Requests", "10"))

	assert.NotNil(t, err)
	assert.Equal(t, err, InvalidOptionType)
	assert.Equal(t, cfg.Threads, 2)


	cfg, err = NewStressConfig(NewOption("threads", 2))
	assert.NotNil(t, err)
	assert.Equal(t, err, InvalidOptionName)

}
func TestStressConfig(t *testing.T) {

	headers := map[string]string{"content-type":"application/json"}

	cfg, err = NewStressConfig(
		NewOption("URL", "http://localhost"),
		NewOption("Requests", 10),
		NewOption("Threads", 2),
		NewOption("Method", "POST"),
		NewOption("Headers", headers),
		NewOption("Duration", 60),
		NewOption("Frequency", 100),
	)

	assert.Nil(t, err)
	assert.NotNil(t, cfg.client)
	assert.Equal(t, cfg.Requests, 10)
	assert.Equal(t, cfg.Threads, 2)
	assert.Equal(t, cfg.Method, "POST")
	assert.Equal(t, cfg.Headers["content-type"], headers["content-type"])
	assert.Equal(t, cfg.Duration, 60)
	assert.Equal(t, cfg.Frequency, 100)

}

