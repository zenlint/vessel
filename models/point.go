package models

import (
	"time"
)

const (
	POINTSTART = iota
	POINTEND
	POINTCHECK
)

type Point struct {
	Id         int64     `json:"id"`                                        //
	PipelineId int64     `json:"pipelineId"`                                //
	Type       int       `json:"type"`                                      //POINTSTART, POINTEND, POINTCHECK
	UUID       string    `json:"uuid" orm:"unique;varchar(255)"`            //
	Name       string    `json:"name" orm:"varchar(255)"`                   //
	From       string    `json:"from" orm:"null;type(text)"`                //mutil Point.UUID or Stage.UUID
	To         string    `json:"to" orm:"null;type(text)"`                  //mutil Point.UUID or Stage.UUID
	Actived    bool      `json:"actived" orm:"null;default(0)"`             //
	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
	Memo       string    `json:"memo" orm:"null;type(text)"`                //
}

func (point *Point) Create(pipelineId int64, pointType int, name string) (error, string) {
	return nil, ""
}

func (point *Point) AddFrom(uuid string, from ...string) error {
	return nil
}

func (point *Point) AddEnd(uuid string, end ...string) error {
	return nil
}

func (point *Point) Run(uuid string) (error, string) {
	return nil, ""
}

func (point *Point) Copy(uuid string) (error, string) {
	return nil, ""
}
