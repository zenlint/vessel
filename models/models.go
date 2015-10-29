package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

var (
	RedisClient *redis.Client
)

func init() {
	orm.RegisterModel(new(Workspace), new(Project), new(Pipeline), new(Status), new(Param), new(Point), new(Stage))
	orm.RegisterDataBase("default", "tidb", "", 0, 0)
}
