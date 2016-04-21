package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

var (
	Orm orm.Ormer
)

func Connect(mysqlport string, mysqladdr string, dbname string, account string, pwd string) error {
	glog.Infoln(mysqladdr, mysqlport, dbname, account, pwd)
	orm, err := connect("default", mysqladdr, mysqlport, account, pwd, dbname)
	Orm = orm
	return err
}
func init() {
	orm.RegisterModel(new(Tb_etcd_backup))
}

//connect to mysql database
func connect(name string, ip string, port string, account string, pwd string, dbname string) (orm.Ormer, error) {
	err := orm.RegisterDataBase(name, "mysql",
		account+":"+pwd+"@tcp("+ip+":"+port+")/"+dbname+"?charset=utf8&loc=Local", 30)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	err = o.Using(name)
	if err != nil {
		return nil, err
	}
	return o, nil
}
func LoadConfig() {

}
