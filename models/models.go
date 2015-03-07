package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

var (
	LedisDB *ledis.DB

	fakeValue = []byte("0")
)

var (
	ErrObjectNotExist = errors.New("Object does not exist")
)

type SetType string

const (
	SET_TYPE_FLOW     SetType = "GLOBAL_FLOW"
	SET_TYPE_PIPELINE SetType = "GLOBAL_PIPELINE"
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

func getSetName(obj interface{}) (SetType, error) {
	var setName SetType
	switch tp := obj.(type) {
	case *Flow:
		setName = SET_TYPE_FLOW
	case *Pipeline:
		setName = SET_TYPE_PIPELINE
	default:
		return "", fmt.Errorf("unknown type: %v", tp)
	}
	return setName, nil
}

// Save saves an object with given UUID.
func Save(uuid string, obj interface{}) error {
	setName, err := getSetName(obj)
	if err != nil {
		return err
	}

	value, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("encoding JSON: %v", err)
	}

	if _, err = LedisDB.HSet([]byte(setName), []byte(uuid), fakeValue); err != nil {
		return fmt.Errorf("HSet: %v", err)
	} else if err = LedisDB.Set([]byte(uuid), value); err != nil {
		return fmt.Errorf("Set: %v", err)
	}

	return nil
}

// Retrieve reads and returns an object with given UUID.
func Retrieve(uuid string, obj interface{}) error {
	value, err := LedisDB.Get([]byte(uuid))
	if err != nil {
		return err
	} else if len(value) == 0 {
		return ErrObjectNotExist
	}

	return json.Unmarshal(value, obj)
}
