package models

import (
	"time"
)

const (
	PIPESUCCESS = iota
	PIPEEERROR
	PIPEEXCEPT
)

type Pipeline struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"projectId"`
	Name        string    `json:"name" orm:"varchar(255)"`
	Description string    `json:"description" orm:"null;type(text)"`
	Actived     bool      `json:"actived" orm:"null;default(0)"`
	Locked      bool      `json:"locked" orm:"null;default(0)"`
	Created     time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `json:"updated" orm:"auto_now;type(datetime)"`
	Memo        string    `json:"memo" orm:"null;type(text)"`
}

type Status struct {
	Id         int64     `json:"id"`
	PipelineId int64     `json:"pipelineId"`
	Started    time.time `json:"started" orm:"type(datetime)"`
	Ended      time.time `json:"ended" orm:"type(datetime)"`
	Log        string    `json:"log" orm:"type(text)"`
	Result     int64     `json:"result" orm:"null"`
	Actived    bool      `json:"actived" orm:"null;default(0)"`
	Locked     bool      `json:"locked" orm:"null;default(0)"`
	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`
	Memo       string    `json:"memo" orm:"null;type(text)"`
}
