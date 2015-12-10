package models

import (
	"time"

	"github.com/huawei-openlab/newdb/orm"
	"github.com/ngaut/log"
)

type Project struct {
	Id          int64     `json:"id"`
	WorkspaceId int64     `json:"workspaceId"`
	Name        string    `json:"name" orm:"varchar(255)"`
	Description string    `json:"description" orm:"null;type(text)"`
	Actived     bool      `json:"actived" orm:"null;default(0)"`
	Created     time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `json:"updated" orm:"auto_now;type(datetime)"`
	Memo        string    `json:"memo" orm:"null;type(text)"`
}

func (project *Project) Create(wid int64, name, description string) (int64, error) {
	o := orm.NewOrm()
	p := Project{WorkspaceId: wid, Name: name, Description: description, Actived: true}

	if err := o.Begin(); err != nil {
		log.Errorf("Transcation error: %s", err.Error())

		return 0, err
	} else {
		if id, e := o.Insert(&p); e != nil {
			log.Errorf("Create project error: %s", e.Error())

			o.Rollback()
			return 0, err
		} else {
			log.Errorf("Create project successfully, id is: %d", id)

			o.Commit()
			return id, nil
		}
	}
}
