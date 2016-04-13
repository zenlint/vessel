package module

// log "github.com/Sirupsen/logrus"
// "github.com/jinzhu/gorm"

// const (
// 	//MaxIdleConnNum : MaxIdleConnNum
// 	MaxIdleConnNum = 0
// 	// MaxOpenConnNum : MaxOpenConnNum
// 	MaxOpenConnNum = 100
// )

// var db *gorm.DB

func InitModule() {

	// var err error
	//
	// //open database
	// db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/k8sci?charset=utf8&parseTime=True&loc=Local")
	// if err != nil {
	// 	log.Panic("erro when conn to db:", err)
	// }
	//
	// db.SingularTable(true)
	//
	// //set db pool
	// db.DB().SetMaxIdleConns(int(MaxIdleConnNum))
	// db.DB().SetMaxOpenConns(int(MaxOpenConnNum))
	//
	// db.LogMode(false)
}
