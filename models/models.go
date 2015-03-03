package models

import (
	"fmt"

	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

var (
	LedisDB *ledis.DB
)

func InitDb() error {
	opt := &config.Config{
		DataDir: "./vessel.db",
	}

	l, err := ledis.Open(opt)
	if err != nil {
		return fmt.Errorf("open Ledis DB: %v", err)
	}

	db := 0
	LedisDB, err = l.Select(db)
	if err != nil {
		return fmt.Errorf("select Ledis DB '%d': %v", db, err)
	}

	return nil
}
