package models

import (
	"time"
)

type Stage struct {
	Id         int64     `json:"id"`                                        //
	PipelineId int64     `json:"pipelineId"`                                //
	UUID       string    `json:"uuid" orm:"unique;varchar(255)"`            //
	Name       string    `json:"name" orm:"varchar(255)"`                   //
	From       string    `json:"from" orm:"null;type(text)"`                //mutil Point.UUID or Stage.UUID
	To         string    `json:"to" orm:"null;type(text)"`                  //mutil Point.UUID or Stage.UUID
	Content    string    `json:"content" orm:"null;type(text)"`             //
	Actived    bool      `json:"actived" orm:"null;default(0)"`             //
	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
	Memo       string    `json:"memo" orm:"null;type(text)"`                //
}

func (stage *Stage) Create(pipelineId int64, name string) (error, string) {
	return nil, ""
}

func (stage *Stage) AddFrom(uuid string, from ...string) error {
	return nil
}

func (stage *Stage) AddEnd(uuid string, end ...string) error {
	return nil
}

func (stage *Stage) Run(uuid string) (error, string) {
	return nil, ""
}

func (stage *Stage) Copy(uuid string) (error, string) {
	return nil, ""
}
