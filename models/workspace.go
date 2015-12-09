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

func (ws *Workspace) Create(name, description string) (int64, error) {
	o := orm.NewOrm()
	w := Workspace{Name: name, Description: description, Actived: true}

	if err := o.Begin(); err != nil {
		log.Errorf("Transcation error: %s", err.Error())

		return 0, err
	} else {
		if id, e := o.Insert(&w); e != nil {
			log.Errorf("Create workspace error: %s", e.Error())

			o.Rollback()
			return 0, e
		} else {
			log.Infof("Create worksapce successfully, id is: %d", id)

			o.Commit()
			return id, nil
		}
	}

	return 0, nil
}

func (ws *Workspace) Put(id int64, description string) error {
	o := orm.NewOrm()
	w := Workspace{Id: id, Actived: true}

	if err := o.Read(&w, "Id", "Actived"); err != nil {
		log.Errorf("Get workspace %d error: &s", id, err.Error())

		return err
	} else {
		if err := o.Begin(); err != nil {
			log.Errorf("Transcation error: %s", err.Error())
			return err
		} else {
			w.Description = description

			if _, err := o.Update(&w, "Description"); err != nil {
				log.Errorf("Put workspace %d error: %s", id, err.Error())

				o.Rollback()
				return err
			} else {
				log.Infof("Put workspace successfully: %d", id)

				o.Commit()
				return nil
			}
		}
	}
}

func (ws *Workspace) Get(id int64) (Workspace, error) {
	o := orm.NewOrm()
	w := Workspace{Id: id, Actived: true}

	if err := o.Read(&w, "Id", "Actived"); err != nil {
		log.Errorf("Get workspace %d error: &s", id, err.Error())

		return w, err
	} else {
		return w, nil
	}
}

func (ws *Workspace) Delete(id int64) error {
	o := orm.NewOrm()
	w := Workspace{Id: id}

	if err := o.Read(&w, "Id"); err != nil {
		log.Errorf("Get workspace %d error: %s", id, err.Error())

		return err
	} else {
		if err := o.Begin(); err != nil {
			log.Errorf("Transcation error: %s", err.Error())
			return err
		} else {
			if _, err := o.Delete(&w); err != nil {
				log.Errorf("Delete workspace %d error: %s", id, err.Error())

				o.Rollback()
				return err
			} else {
				log.Infof("Delete workspace %d successfully", id)

				//TODO Delete relate projects.

				o.Commit()
				return nil
			}
		}
	}
}
