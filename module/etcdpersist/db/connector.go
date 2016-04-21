package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Orm orm.Ormer
)

func init() {
	orm.RegisterModel(new(Tb_etcd_backup))
}

func Syncdb(db_user, db_pass, db_host, db_port, db_name string) error {
	err := Create(db_user, db_pass, db_host, db_port, db_name)
	if err != nil {
		return err
	}
	err = Connect(db_user, db_pass, db_host, db_port, db_name)
	if err != nil {
		return err
	}
	name := "default"
	force := true
	verbose := true
	err = orm.RunSyncdb(name, force, verbose)
	log.Println("database init is complete.\nPlease restart the application")
	return err
}

func Connect(db_user, db_pass, db_host, db_port, db_name string) error {
	var dns string
	orm.RegisterDriver("mysql", orm.DRMySQL)
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
	err := orm.RegisterDataBase("default", "mysql", dns)
	if err != nil {
		return err
	}
	Orm = orm.NewOrm()
	err = Orm.Using("default")
	return err
}

func Create(db_user, db_pass, db_host, db_port, db_name string) error {
	var dns string
	var sqlstring string
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
	sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err.Error())
	}
	r, err := db.Exec(sqlstring)
	if err != nil {
		log.Println(err, r)
	} else {
		log.Println("Database ", db_name, " created")
	}
	defer db.Close()
	return err
}
