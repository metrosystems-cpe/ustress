package core

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	Config := newConfig("127.0.0.1")
	a := NewApp("0.0.1", Config)
	a.Init()
	s := StressTest{}
	d, e := s.All(a.Session)
	fmt.Println(d, e)

}
