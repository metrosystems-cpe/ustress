package core

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/gocql/gocql"
	log "github.com/metrosystems-cpe/ustress/log"
	"github.com/metrosystems-cpe/ustress/ustress"
)

const StressTestTableName = "stress_test"

// ScheduledTest is an incoming feature
const ScheduledTestTableName = "scheduled_test"

// Tables metadata that will be used to generate tables
var Tables = map[string]map[string]string{
	StressTestTableName: map[string]string{
		"id":     "UUID PRIMARY KEY",
		"report": "text",
		"meta":   "WITH default_time_to_live = " + "604800", // one week
	},
}

type StressTest struct {
	ID     uuid.UUID       `gocql:"id"`
	Report *ustress.Report `gocql:"report"`
}

func (test *StressTest) Get(sess *gocql.Session) error {
	q := fmt.Sprintf("select from %s where id = ?", StressTestTableName)
	mapString := map[string]interface{}{}
	err := sess.Query(q, gocql.UUID(test.ID)).MapScan(mapString)
	if mapString["report"] != nil {
		v, _ := mapString["report"].(string)
		json.Unmarshal([]byte(v), test.Report)
		return err
	}
	return err
}

func (test *StressTest) Save(sess *gocql.Session) error {
	q := fmt.Sprintf("insert into %s (id, report) values (?, ?)", StressTestTableName)
	err := sess.Query(q, gocql.UUID(test.ID), test.Report.JSON()).Exec()
	log.LogWithFields.Infof("[INSERT] row into table %s", StressTestTableName)
	return err
}

func (test *StressTest) All(sess *gocql.Session) ([]map[string]interface{}, error) {
	q := fmt.Sprintf("SELECT * FROM %s", StressTestTableName)
	return sess.Query(q).Iter().SliceMap()
}

func NewStressTest(report *ustress.Report) *StressTest {
	return &StressTest{ID: report.UUID, Report: report}
}
