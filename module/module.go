package module

import (
	"github.com/containerops/vessel/models"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

const (
	//MaxIdleConnNum : MaxIdleConnNum
	MaxIdleConnNum = 0
	// MaxOpenConnNum : MaxOpenConnNum
	MaxOpenConnNum = 100
)

var db *gorm.DB

func init() {

	var err error

	//打开数据库获连接
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/k8sci?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panic("erro when conn to db:", err)
	}
	//设置数据库名称单数
	db.SingularTable(true)

	//设置池子大小
	db.DB().SetMaxIdleConns(int(MaxIdleConnNum))
	db.DB().SetMaxOpenConns(int(MaxOpenConnNum))

	db.LogMode(false)
}

// GetProjectInfoByName :
func GetProjectInfoByName(projectName string) (project *models.Project) {
	project = new(models.Project)
	db.Model(&models.Project{}).Where("name = ?", projectName).First(project)
	return project
}

// CreatePipeline :
func CreatePipeline(pipeline *models.Pipeline) error {
	return db.Model(&models.Pipeline{}).Save(pipeline).Error
}

// CreatePoint :
func CreatePoint(point *models.Point) error {
	return db.Model(&models.Point{}).Save(point).Error
}

// CreateStage :
func CreateStage(stage *models.Stage) error {
	return db.Model(&models.Stage{}).Save(stage).Error
}
