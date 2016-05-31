package stage

import (
	"github.com/containerops/vessel/models"
	"strings"
	"time"
)

type Stage struct {
	info       *models.Stage
	dependence map[string]bool
	Result     string
	Err        error
}

func (stage *Stage)DependenceReady(name string) bool {
	_, ok := stage.dependence[name]
	if ok {
		stage.dependence[name] = true
	}
	for _, item := range stage.dependence {
		if !item {
			return false
		}
	}
	return true
}

func (stage *Stage)Info() *models.Stage {
	if stage.info == nil {
		stage.info = &models.Stage{}
	}
	return stage.info
}

func (stage Stage)SetDependence(names string) {
	list := strings.Split(names, ",")
	stage.dependence = make(map[string]bool)
	for _, item := range list {
		if item == "" {
			continue
		}
		stage.dependence[item] = false
	}
	stage.info.Dependence = list
}

func (stage Stage)GetDependence() []string {
	return stage.info.Dependence
}

func (stage *Stage)Start(finishChan chan models.ExecutorRes, endTime time.Time) {

}

func (stage *Stage)Stop(finishChan chan models.ExecutorRes, endTime time.Time) {

}