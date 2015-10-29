package models

import (
	"time"
)

type Param struct {
	Id       int64     `json:"id"`                                        //
	ActionId string    `json:"actionId" orm:"varchar(255)"`               //Point.UUID; Stage.UUID; Pipeline.UUID
	Key      string    `json:"key" orm:"varchar(255)"`                    //
	Value    string    `json:"value" orm:"type(text)"`                    //
	Actived  bool      `json:"actived" orm:"null;default(0)"`             //
	Created  time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
	Updated  time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
	Memo     string    `json:"memo" orm:"null;type(text)"`                //
}

func (param *Param) Create(uuid, key, value string) (error, string) {
	return nil, ""
}

func (param *Param) Copy(uuid string) (error, string) {
	return nil, ""
}
