package models

import (
	"time"

	"github.com/huawei-openlab/newdb/orm"
	"github.com/ngaut/log"
)

type Workspace struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" orm:"unique;varchar(255)"`
	Description string    `json:"description" orm:"null;type(text)"`
	Actived     bool      `json:"actived" orm:"null;default(0)"`
	Created     time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `json:"updated" orm:"auto_now;type(datetime)"`
	Memo        string    `json:"memo" orm:"null;type(text)"`
}

func (workspace *Workspace) Create(name, description string) (int64, error) {
	o := orm.NewOrm()
	w := Workspace{Name: name, Description: description}

	if err := o.Begin(); err != nil {
		log.Errorf("Transcation error: %s", err.Error())
		return 0, err
	} else {
		if id, e := o.Insert(&w); e != nil {
			log.Errorf("Insert workspace data error: %s", e.Error())

			o.Rollback()
			return 0, e
		} else {
			log.Errorf("Insert worksapce successfully, id is: %d", id)

			o.Commit()
			return id, nil
		}
	}

	return 0, nil
}
