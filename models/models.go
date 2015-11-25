package models

import (
	"github.com/huawei-openlab/newdb/orm"
	_ "github.com/pingcap/tidb"
)

func init() {
	orm.RegisterModel(new(Workspace), new(Project), new(Pipeline), new(Status), new(Param), new(Point), new(Stage))
	orm.RegisterDataBase("default", "tidb", "goleveldb:///tmp/tidb/vessel", 0, 0)
}
