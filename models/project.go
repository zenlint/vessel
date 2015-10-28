package models

import (
	"time"
)

type Project struct {
	Id          int64     `json:"id"`
	WorkspaceId int64     `json:"workspaceId"`
	Name        string    `json:"name" orm:"varchar(255)"`
	Description string    `json:"description" orm:"null;type(text)"`
	Actived     bool      `json:"actived" orm:"null;default(0)"`
	Locked      bool      `json:"locked" orm:"null;default(0)"`
	Created     time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `json:"updated" orm:"auto_now;type(datetime)"`
	Memo        string    `json:"memo" orm:"null;type(text)"`
}
