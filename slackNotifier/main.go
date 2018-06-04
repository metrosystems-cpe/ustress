package slackNotifier

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	log "git.metrosystems.net/reliability-engineering/traffic-monkey/log"
	"github.com/ashwanthkumar/slack-go-webhook"
)

// TODO: Find a better way to handle host :(
// For testing https://hooks.slack.com/services/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"
// For pp: https://hooks.slack.com/services/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
const (
	webhook = "https://hooks.slack.com/services/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"
	deliver = true // set to true to deliver slack notifications
)

var (
	authorName  = "Traffic Monkey"
	authorImage = "https://i.pinimg.com/736x/29/61/55/29615560c6387dd576b3076eae0b760d--cartoon-monkey-family-guy.jpg"
)

func host() string {
	var h string
	if h = os.Getenv("HTTP_INGRESS"); h == "" {
		return "http://localhost:8080/"
	}

	return h
}

// RawParams keeps necessay data to send notifications on slack
type RawParams struct {
	Link       string
	NrRequests int
	NrThreads  int
	Result     []byte
}

// FieldsList is a list of *slack.Field
type FieldsList []*slack.Field

func floatToString(inputNum interface{}) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum.(float64), 'f', 6, 64)
}

// ComputeAttachmentFields will return a list of custom fields for attachment
func (data RawParams) ComputeAttachmentFields() FieldsList {
	var out map[string]interface{}
	json.Unmarshal(data.Result, &out)
	stats := out["stats"].(map[string]interface{})

	list := FieldsList{
		&slack.Field{Title: "Stressed URL", Value: data.Link},
		&slack.Field{Title: "Number of Requests", Value: strconv.Itoa(data.NrRequests), Short: true},
		&slack.Field{Title: "Number of Threads", Value: strconv.Itoa(data.NrThreads), Short: true},
		&slack.Field{Title: "Median", Value: floatToString(stats["median"]), Short: true},
		&slack.Field{Title: "50 Percentile", Value: floatToString(stats["50_percentile"]), Short: true},
		&slack.Field{Title: "75 Percentile", Value: floatToString(stats["75_percentile"]), Short: true},
		&slack.Field{Title: "95 Percentile", Value: floatToString(stats["95_percentile"]), Short: true},
		&slack.Field{Title: "99 Percentile", Value: floatToString(stats["99_percentile"]), Short: true},
		&slack.Field{Title: "Error Percentage", Value: floatToString(stats["error_percentage"]), Short: true},
	}

	return list
}

// ReportLink compute full url for full report json file
func (data RawParams) ReportLink() string {
	var out map[string]interface{}
	json.Unmarshal(data.Result, &out)
	return fmt.Sprintf("%s/ui/?report_id=%s", host(), out["uuid"].(string))
}

// DeliverReport is used to deliver stres test report as slack notification
func DeliverReport(params RawParams) {
	attachment := slack.Attachment{Fields: params.ComputeAttachmentFields()}
	payload := slack.Payload{
		Text:        fmt.Sprintf("New Stress Test was performed. Full report can be found here: %s", params.ReportLink()),
		Attachments: []slack.Attachment{attachment},
	}
	attachment.AuthorName = &authorName
	attachment.AuthorIcon = &authorImage

	if deliver {
		err := slack.Send(webhook, "", payload)
		if len(err) > 0 {
			log.LogWithFields.Debugf("%+v", err)
		} else {
			log.LogWithFields.Println("Success delivering Slack notification")
		}
	} else {
		log.LogWithFields.Info("Message prepared but deliver constatnt is set to false :(")
	}
}
