package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.metrosystems.net/reliability-engineering/ustress/web"
	"git.metrosystems.net/reliability-engineering/ustress/web/core"

	"github.com/google/uuid"
)

func TestInsert(t *testing.T) {
	//init
	u := uuid.New()
	q1 := core.Insert(core.StressTestTableName, []string{"id", "report"}, 1, "'2'")
	q2 := core.Insert(core.StressTestTableName, []string{"id", "report"}, u, "'uuid'")
	q3 := core.Insert(core.StressTestTableName, []string{"id", "report"}, "'1'", "2")
	expected1 := fmt.Sprintf("INSERT INTO %s ( id, report ) VALUES ( 1, '2' )", core.StressTestTableName)
	expected2 := fmt.Sprintf("INSERT INTO %s ( id, report ) VALUES ( %v, 'uuid' )", core.StressTestTableName, u)
	expected3 := fmt.Sprintf("INSERT INTO %s ( id, report ) VALUES ( '1', 2 )", core.StressTestTableName)
	assert.Equal(t, expected1, q1)
	assert.Equal(t, expected2, q2)
	assert.Equal(t, expected3, q3)

}

func TestSelect(t *testing.T) {
	u := uuid.New()
	q1 := core.Select(core.StressTestTableName, "id", u)
	expected1 := fmt.Sprintf("SELECT * FROM %s WHERE id = %v", core.StressTestTableName, u)
	assert.Equal(t, expected1, q1)

}

func TestAll(t *testing.T) {
	Config := newConfig("127.0.0.1")
	a := web.NewApp("0.0.1", Config)
	s := core.StressTest{}
	d, e := s.All(a.Session)
	fmt.Println(d, e)

}
