package ustress

import (
	"encoding/json"
	"time"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	"github.com/google/uuid"
	"github.com/montanaflynn/stats"
)

// Stats Structure
type Stats struct {
	Median          float64 `json:"median"`
	PercentileA     float64 `json:"50_percentile"`
	PercentileB     float64 `json:"75_percentile"`
	PercentileC     float64 `json:"95_percentile"`
	PercentileD     float64 `json:"99_percentile"`
	ErrorPercentage float64 `json:"error_percentage"`
}

// Report structure
type Report struct {
	UUID      uuid.UUID     `json:"uuid"`
	TimeStamp time.Time     `json:"timestamp"`
	Config    *MonkeyConfig `json:"config"`
	Stats     Stats         `json:"stats"`
	Duration  float64       `json:"durationTotal"`
	Data      []WorkerData  `json:"data"`
}

func (report *Report) CalcStats() *Report {
	var requestDurations []float64
	var numberOfErrors int
	var err error
	for _, value := range report.Data {
		// ignore http codes 100s to 500s
		if value.Status > 100 && value.Status < 600 {
			requestDurations = append(requestDurations, value.Duration)
		} else {
			numberOfErrors++
		}
	}

	if report.Stats.PercentileA, err = stats.Percentile(requestDurations, 50); err != nil {
		report.Stats.PercentileA = 0
	}
	if report.Stats.PercentileB, err = stats.Percentile(requestDurations, 75); err != nil {
		report.Stats.PercentileB = 0
	}
	if report.Stats.PercentileC, err = stats.Percentile(requestDurations, 95); err != nil {
		report.Stats.PercentileC = 0
	}
	if report.Stats.PercentileD, err = stats.Percentile(requestDurations, 99); err != nil {
		report.Stats.PercentileD = 0
	}
	if report.Stats.Median, err = stats.Median(requestDurations); err != nil {
		report.Stats.Median = 0
	}

	report.Stats.ErrorPercentage = float64(numberOfErrors) / float64(report.Config.Requests) * 100
	return report
}

func (r *Report) JSON() []byte {
	b, err := json.Marshal(r)
	log.LogError(err)
	return b
}
