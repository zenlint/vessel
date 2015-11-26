package models

import (
	"fmt"

	"github.com/huawei-openlab/newdb/orm"

	"github.com/ngaut/log"
	_ "github.com/pingcap/tidb"
)

//Init Database
func init() {
	log.SetLevelByString("info")

	orm.RegisterModel(new(Workspace), new(Project), new(Pipeline), new(Status), new(Param), new(Point), new(Stage))
	orm.RegisterDataBase("default", "tidb", "goleveldb:///tmp/tidb/vessel", 0, 0)
}

//Sync Database
func Sync(force, verbose bool) error {
	if err := orm.RunSyncdb("default", force, verbose); err != nil {
		return fmt.Errorf("Sync Database Error, ", err.Error())
	}

	return nil
}
