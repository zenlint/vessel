package dependence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
)

func ParsePipelineTemplate(template *models.PipelineSpecTemplate) (map[string][]*stage.Stage, error) {
	pipelineInfo := template.MetaData

	stageSpec := template.Spec
	stageDagMap := make(map[string][]*stage.Stage, 0)
	stageMap := make(map[string]*stage.Stage, 0)

	pipelineInfo.Stages = make([]string, 0, len(stageSpec))
	for _, stageInfo := range stageSpec {
		if stageInfo.Name == "" {
			return nil, errors.New("Stage has an empty name")
		}
		_, ok := stageMap[stageInfo.Name]
		if ok {
			return nil, errors.New(fmt.Sprintf("Stage has repeat name: %v", stageInfo.Name))
		}

		currStage := createStage(pipelineInfo.Namespace, stageInfo)
		stageMap[stageInfo.Name] = currStage
		for _, dependenceItem := range currStage.Dependence {
			stageList, ok := stageDagMap[dependenceItem];
			if !ok {
				stageList = make([]*stage.Stage, 0, 10)
			}
			stageList = append(stageList, currStage)
			stageDagMap[dependenceItem] = stageList
		}
		pipelineInfo.Stages = append(pipelineInfo.Stages, stageInfo.Name)
	}
	return stageDagMap, checkDependenceValidity(stageDagMap, stageMap)
}

func createStage(namespace string, stageInfo *models.Stage) *stage.Stage {
	stage := &stage.Stage{}
	stageInfo.Namespace = namespace
	stage.SetInfo(stageInfo)
	return stage
}

func checkDependenceValidity(stageDagMap map[string][]*stage.Stage, stageMap map[string]*stage.Stage) error {
	if len(stageDagMap[""]) == 0 {
		return errors.New("The first start stage list is empty")
	}

	//Check dependence stage name is exist
	for dependenceName, _ := range stageDagMap {
		if dependenceName == "" {
			continue
		}
		_, ok := stageMap[dependenceName]
		if !ok {
			return errors.New(fmt.Sprintf("Dependence stage name: %v is not exist", dependenceName))
		}
	}

	//Check dependence directed acyclic graph
	return checkEndlessChain(stageDagMap, make([]string, 0, 10), nil)
}

func checkEndlessChain(stageMap map[string][]*stage.Stage, chain []string, stage *stage.Stage) error {
	var stageName string
	if stage == nil {
		stageName = ""
	} else {
		stageName = stage.GetName()
		for _, chainItem := range chain {
			if chainItem == stageName {
				return errors.New(fmt.Sprintf("Dependence chain [%v,%v] is endless chain",
					strings.Join(chain, ","), stageName))
			}
		}
	}
	stageList, ok := stageMap[stageName]
	if ok {
		for _, nextStage := range stageList {
			chain = append(chain, stageName)
			err := checkEndlessChain(stageMap, chain, nextStage)
			if err != nil {
				return err
			}
			chain = chain[0:len(chain) - 1]
		}
	}
	return nil
}