package slackNotifier

import (
	"encoding/json"
	"fmt"
	"strconv"

	log "git.metrosystems.net/reliability-engineering/traffic-monkey/log"
	"github.com/ashwanthkumar/slack-go-webhook"
)

// TODO: Find a better way to handle host :(
const (
	webhook = "https://hooks.slack.com/services/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"
	host    = "http://localhost:9090"
)

// RawParams keeps necessay data to send notifications on slack
type RawParams struct {
	Link       string
	NrRequests int
	NrThreads  int
	Result     []byte
}

// FieldsList is a list of *slack.Field
type FieldsList []*slack.Field

// ComputeAttachmentFields will return a list of custom fields for attachment
func (data RawParams) ComputeAttachmentFields() FieldsList {
	list := FieldsList{
		&slack.Field{Title: "Stressed URL", Value: data.Link, Short: true},
		&slack.Field{Title: "Number of Requests", Value: strconv.Itoa(data.NrRequests), Short: true},
		&slack.Field{Title: "Number of Threads", Value: strconv.Itoa(data.NrThreads), Short: true},
	}

	return list
}

// ReportLink compute full url for full report json file
func (data RawParams) ReportLink() string {
	var out map[string]interface{}
	json.Unmarshal(data.Result, &out)
	return fmt.Sprintf("%s/data/%s.json", host, out["uuid"].(string))
}

// DeliverReport is used to deliver stres test report as slack notification
func DeliverReport(params RawParams) {

	attachment := slack.Attachment{Fields: params.ComputeAttachmentFields()}
	payload := slack.Payload{
		Text:        fmt.Sprintf("New Stress Test was performed. Full report can be found here: %s", params.ReportLink()),
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.Send(webhook, "", payload)

	if len(err) > 0 {
		log.LogWithFields.Debugf("%+v", err)
	} else {
		log.LogWithFields.Println("Success delivering Slack notification")
	}
}
