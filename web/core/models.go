package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
	"git.metrosystems.net/reliability-engineering/ustress/ustress"
	"github.com/gocql/gocql"
)

const StressTestTableName = "stress_test"

// ScheduledTest is an incoming feature
const ScheduledTestTableName = "scheduled_test"

var Tables = map[string]map[string]string{
	StressTestTableName: map[string]string{
		"id":     "UUID PRIMARY KEY",
		"report": "text",
		"meta":   "WITH default_time_to_live = " + "604800", // one week
	},
	// ScheduledTestTableName: map[string]string{
	// 	"id":            "UUID PRIMARY KEY",
	// 	"repeat_period": "INTEGER",
	// 	"report":        "TEXT",
	// },
}

func Select(tableName string, pk interface{}, val interface{}) string {
	return fmt.Sprintf(
		"SELECT * FROM %s WHERE %s = %v",
		tableName,
		pk,
		val)
}

func Insert(tableName string, keys []string, vals ...interface{}) string {
	q := fmt.Sprintf("INSERT INTO %s ( ", tableName)
	q += strings.Join(keys, ", ")

	q += " ) VALUES ( "
	for _, v := range vals {
		q += fmt.Sprintf("%v, ", v)
	}
	q = strings.TrimRight(q, " , ")
	q += " )"

	return q
}

// type ScheduledTest struct {
// 	RepeatPeriod int // 0 no repeat
// 	ScheduleTime time.Time
// 	Config       ustress.MonkeyConfig
// }

type StressTest struct {
	ID     uuid.UUID       `gocql:"id"`
	Report *ustress.Report `gocql:"report"`
}

func (test *StressTest) Get(sess *gocql.Session) error {
	q := Select(StressTestTableName, "id", gocql.UUID(test.ID))
	mapString := map[string]interface{}{}

	err := sess.Query(q).MapScan(mapString)
	if mapString["report"] != nil {
		v, _ := mapString["report"].(string)
		json.Unmarshal([]byte(v), test.Report)
		return err
	}
	return err
}

func (test *StressTest) Save(sess *gocql.Session) error {
	q := Insert(
		StressTestTableName,
		[]string{"id", "report"},
		gocql.UUID(test.ID),
		fmt.Sprintf("'%s'", string(test.Report.JSON())),
	)
	err := sess.Query(q).Exec()
	log.LogWithFields.Infof("Inserting row into table %s", StressTestTableName)
	return err
}

func (test *StressTest) All(sess *gocql.Session) ([]map[string]interface{}, error) {
	q := fmt.Sprintf("SELECT * FROM %s", StressTestTableName)
	return sess.Query(q).Iter().SliceMap()

}
