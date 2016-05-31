package dependence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
)

func ParsePipelineTemplate(template *models.PipelineSpecTemplate) (*models.Pipeline, map[string][]*stage.Stage, error) {
	pipelineInfo := createPipelineInfo(template.MetaData)

	stageSpec := template.Spec
	stageDagMap := make(map[string][]*stage.Stage, 0)
	stageMap := make(map[string]*stage.Stage, 0)
	dependenceCountMap := make(map[string]int, 0)

	pipelineInfo.Stages = make([]string, 0, len(stageSpec))
	for _, spec := range stageSpec {
		if spec.Name == "" {
			return nil, nil, errors.New("Stage has an empty name")
		}

		_, ok := stageMap[spec.Name]
		if ok {
			return nil, nil, errors.New(fmt.Sprintf("Stage has repeat name: %v", spec.Name))
		}
		//init stage dependence count
		if _, ok := stageMap[spec.Name]; !ok {
			dependenceCountMap[spec.Name] = 0
		}
		currStage := createStage(pipelineInfo.Namespace, spec)
		stageMap[spec.Name] = currStage
		dependence := currStage.GetDependence()
		for _, dependenceItem := range dependence {
			stageList, ok := stageDagMap[dependenceItem];
			if !ok {
				stageList = make([]*stage.Stage, 0, 10)
			}
			stageList = append(stageList, currStage)
			stageDagMap[dependenceItem] = stageList
		}
		pipelineInfo.Stages = append(pipelineInfo.Stages, spec.Name)
	}
	return pipelineInfo, stageDagMap, checkDependenceValidity(stageDagMap, stageMap)
}

func createPipelineInfo(metaData *models.PipelineMetaData) *models.Pipeline {
	return &models.Pipeline{
		Name:metaData.Name,
		Namespace:metaData.Namespace,
		TimeoutDuration:int64(metaData.TimeoutDuration),
	}
}

func createStage(namespace string, stageSpec *models.StageSpec) *stage.Stage {
	stage := &stage.Stage{}
	stageInfo := stage.Info()
	stageInfo.Namespace = namespace
	stageInfo.Name = stageSpec.Name
	stageInfo.Port = int64(stageSpec.Port)
	stageInfo.Image = stageSpec.Image
	stageInfo.StatusCheckLink = stageSpec.StatusCheckUrl
	stageInfo.StatusCheckInterval = int64(stageSpec.StatusCheckInterval)
	stageInfo.StatusCheckCount = int64(stageSpec.StatusCheckCount)
	stageInfo.EnvName = stageSpec.EnvName
	stageInfo.EnvValue = stageSpec.EnvValue
	stage.SetDependence(stageSpec.Dependence)
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
		stageName = stage.Info().Name
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