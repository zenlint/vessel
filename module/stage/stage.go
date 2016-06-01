package stage

import (
	"github.com/containerops/vessel/models"
	"strings"
	"time"
)

type Stage struct {
	readies    map[string]bool
	info       *models.Stage
	result     *models.StageResult
	Dependence []string
}

func (self *Stage)Start(finishChan chan models.Result, endTime time.Time) {

}

func (self *Stage)IsReady(dependenceName string) bool {
	if _, ok := self.readies[dependenceName]; ok {
		self.readies[dependenceName] = true
	}
	for _, item := range self.readies {
		if !item {
			return false
		}
	}
	return true
}

func (self *Stage) GetResult() *models.StageResult {
	return self.result
}

func (self *Stage) GetName() string {
	return self.info.Name
}

func (self *Stage) SetInfo(info *models.Stage) {
	self.info = info
	list := strings.Split(info.Dependence, ",")
	self.readies = make(map[string]bool)
	for _, item := range list {
		self.readies[item] = false
	}
	self.Dependence = list
}

func (self *Stage)Stop(finishChan chan models.Result, endTime time.Time) {

}