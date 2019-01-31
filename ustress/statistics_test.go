package ustress

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCalcStats(t *testing.T) {

	var (
		monkeyConfig = &MonkeyConfig{
			URL:      "http://localhost:8080/test",
			Requests: 16,
			Threads:  4,
			Resolve:  "",
			Insecure: false,
		}

		data = []WorkerData{
			{Request: 1, Status: 200, Thread: 1, Duration: 0.002153429, MonkeyConfig: monkeyConfig},
			{Request: 2, Status: 200, Thread: 3, Duration: 0.00088149, MonkeyConfig: monkeyConfig},
			{Request: 3, Status: 200, Thread: 4, Duration: 0.000946606, MonkeyConfig: monkeyConfig},
			{Request: 4, Status: 0, Thread: 2, Duration: 0.001074489, Error: "Get https://foo.bar.com/foo/bar: dial tcp: lookup foo.bar.com: no such host", MonkeyConfig: monkeyConfig},
			{Request: 5, Status: 0, Thread: 1, Duration: 0.000819102, Error: "Get  https://foo.bar.com/foo/bar: dial tcp: lookup foo.bar.com: no such host", MonkeyConfig: monkeyConfig},
			{Request: 6, Status: 200, Thread: 3, Duration: 0.000621576, MonkeyConfig: monkeyConfig},
			{Request: 7, Status: 200, Thread: 4, Duration: 0.001068274, MonkeyConfig: monkeyConfig},
			{Request: 8, Status: 200, Thread: 2, Duration: 0.001021386, MonkeyConfig: monkeyConfig},
			{Request: 9, Status: 0, Thread: 1, Duration: 0.001170958, Error: "Get  https://foo.bar.com/foo/bar: dial tcp: lookup foo.bar.com: no such host", MonkeyConfig: monkeyConfig},
			{Request: 10, Status: 0, Thread: 3, Duration: 0.001052171, Error: "Get  https://foo.bar.com/foo/bar: dial tcp: lookup foo.bar.com: no such host", MonkeyConfig: monkeyConfig},
			{Request: 11, Status: 200, Thread: 3, Duration: 0.000621576, MonkeyConfig: monkeyConfig},
			{Request: 12, Status: 200, Thread: 4, Duration: 0.001068274, MonkeyConfig: monkeyConfig},
			{Request: 13, Status: 200, Thread: 2, Duration: 0.001021386, MonkeyConfig: monkeyConfig},
			{Request: 14, Status: 200, Thread: 3, Duration: 0.000621576, MonkeyConfig: monkeyConfig},
			{Request: 15, Status: 200, Thread: 4, Duration: 0.001068274, MonkeyConfig: monkeyConfig},
			{Request: 16, Status: 200, Thread: 2, Duration: 0.001021386, MonkeyConfig: monkeyConfig},
		}

		report = Report{
			TimeStamp: time.Now(),
			UUID:      uuid.New(),
			Data:      data,
			Config:    monkeyConfig,
		}
	)

	expectedErrorPercentage := float64(25)
	report.CalcStats()
	// t.Logf("%v", report.Stats.ErrorPercentage)
	if report.Stats.ErrorPercentage != expectedErrorPercentage {
		t.Errorf("ErrorPercentage calculation failed: expected %6f, got %6f ",
			expectedErrorPercentage, report.Stats.ErrorPercentage)
	}
}
