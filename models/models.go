package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(Workspace), new(Project), new(Pipeline), new(Status), new(Param), new(Point), new(Stage))
	orm.RegisterDataBase("default", "tidb", "", 0, 0)
}
