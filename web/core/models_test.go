package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestInsert(t *testing.T) {
	//init
	u := uuid.New()
	q1 := Insert(StressTestTableName, []string{"id", "report"}, 1, "'2'")
	q2 := Insert(StressTestTableName, []string{"id", "report"}, u, "'uuid'")
	q3 := Insert(StressTestTableName, []string{"id", "report"}, "'1'", "2")
	expected1 := fmt.Sprintf("INSERT INTO %s ( id, report ) VALUES ( 1, '2' )", StressTestTableName)
	expected2 := fmt.Sprintf("INSERT INTO %s ( id, report ) VALUES ( %v, 'uuid' )", StressTestTableName, u)
	expected3 := fmt.Sprintf("INSERT INTO %s ( id, report ) VALUES ( '1', 2 )", StressTestTableName)
	assert.Equal(t, expected1, q1)
	assert.Equal(t, expected2, q2)
	assert.Equal(t, expected3, q3)

}

func TestSelect(t *testing.T) {
	u := uuid.New()
	q1 := Select(StressTestTableName, "id", u)
	expected1 := fmt.Sprintf("SELECT * FROM %s WHERE id = %v", StressTestTableName, u)
	assert.Equal(t, expected1, q1)

}

func TestAll(t *testing.T) {
	Config := newConfig("127.0.0.1")
	a := NewApp("0.0.1", Config)
	a.Init()
	s := StressTest{}
	d, e := s.All(a.Session)
	fmt.Println(d, e)

}
