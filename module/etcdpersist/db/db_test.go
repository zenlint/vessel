package db

import (
	"strings"
	"testing"
)

func Test_db(t *testing.T) {
	s := "/dir"
	index := strings.LastIndex(s, "/")
	r := []rune(s)
	s = string(r[0:index])
	t.Log(s == "")
}
