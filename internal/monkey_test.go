package internal

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	workers = []Worker{
		{Request: 1, Status: 200, Thread: 1, Duration: 0.002153429},
		{Request: 2, Status: 200, Thread: 3, Duration: 0.00088149},
		{Request: 3, Status: 200, Thread: 4, Duration: 0.000946606},
		{Request: 4, Status: 0, Thread: 2, Duration: 0.001074489, Error: "Get https://idamm.metrosystems.net/.well-known/openid-configuration: dial tcp: lookup idamm.metrosystems.net: no such host"},
		{Request: 5, Status: 0, Thread: 1, Duration: 0.000819102, Error: "Get https://idamm.metrosystems.net/.well-known/openid-configuration: dial tcp: lookup idamm.metrosystems.net: no such host"},
		{Request: 6, Status: 200, Thread: 3, Duration: 0.000621576},
		{Request: 7, Status: 200, Thread: 4, Duration: 0.001068274},
		{Request: 8, Status: 200, Thread: 2, Duration: 0.001021386},
		{Request: 9, Status: 0, Thread: 1, Duration: 0.001170958, Error: "Get https://idamm.metrosystems.net/.well-known/openid-configuration: dial tcp: lookup idamm.metrosystems.net: no such host"},
		{Request: 10, Status: 0, Thread: 3, Duration: 0.001052171, Error: "Get https://idamm.metrosystems.net/.well-known/openid-configuration: dial tcp: lookup idamm.metrosystems.net: no such host"},
	}

	report = Report{
		TimeStamp: time.Now(),
		UUID:      uuid.New(),
		Workers:   workers,
		MonkeyConfig: MonkeyConfig{
			Requests: 10,
		},
	}
)

func TestCalcStats(t *testing.T) {
	report.calcStats()

	t.Logf("%v", report.Stats.ErrorPercentage)

	if report.Stats.ErrorPercentage == 0 {
		t.Error("ErrorPercentage")
	}
}
