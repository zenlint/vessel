package pipeline

import (
	"errors"
	"strings"

	"github.com/containerops/vessel/models"
)

var (
	DEFAULT_PIPELINE_ETCD_PATH        = "/containerops/vessel/ws-%d/pj-%d/pl-%d/stage/"
	DEFAULT_PIPELINEVERSION_ETCD_PATH = "/containerops/vessel/ws-%d/pj-%d/pl-%d/version/plv-%d"
)

// RunPipeline : run pipeline generate pipelineVersion
func RunPipeline(pl *models.Pipeline) (*models.Pipeline, error) {
	// first test is pipeline legal if not return err
	relationMap, err := isPipelineLegal(pl)
	if err != nil {
		return nil, err
	}

	// save pipeline info to db
	id, err := pl.Save()
	if err != nil {
		return nil, err
	}
	pl.Id = id

	// save stage infos to db
	for _, stage := range pl.Stages {
		if relationMap[stage.Name][0] != "" {
			stage.From = strings.Split(relationMap[stage.Name][0], ",")
		}
		if relationMap[stage.Name][1] != "" {
			stage.To = strings.Split(relationMap[stage.Name][1], ",")
		}
		stage.Save()
	}

	return pl, nil
}

// test is the given pipeline is legal ,if legal return pipeline's stage relationMap if not return error
func isPipelineLegal(pipeline *models.Pipeline) (map[string][]string, error) {
	stageMap := make(map[string]*models.Stage, 0)
	dependenceCount := make(map[string]int, 0)
	stageRelationMap := make(map[string][]string, 0)

	// regist all stage,and check repeat/nil stage name
	for _, stage := range pipeline.Stages {
		if stage.Name == "" {
			return nil, errors.New("stage has a nil name")
		}
		if _, exist := stageMap[stage.Name]; !exist {
			stageMap[stage.Name] = stage

			// init stage dependence count
			dependenceCount[stage.Name] = 0

			// count stage dependence
			for _, from := range stage.From {
				dependenceCount[from]++
			}
		} else {
			// has a repeat stage name ,return
			return nil, errors.New("stage has repeat name:" + stage.Name)
		}
	}

	// check DAG
	//if AnnulusTag == nowReleaseStageCount or nowReleaseStageCount == len(dependenceCount) then exit for,if nowReleaseStageCount == len(dependenceCount) then isDAG,else isNotDAG
	nowReleaseStageCount := 0
	for true {

		annulusTag := 0
		for stageName, stage := range stageMap {
			if dependenceCount[stageName] == 0 {
				nowReleaseStageCount++
				for _, from := range stage.From {
					dependenceCount[from]--
				}

				dependenceCount[stage.Name] = -1
			} else if dependenceCount[stageName] == -1 {
				annulusTag++
			}
		}

		if annulusTag == nowReleaseStageCount || nowReleaseStageCount == len(dependenceCount) {
			break
		}
	}

	if nowReleaseStageCount != len(dependenceCount) {
		return nil, errors.New("given pipeline's stage can't create a DAG")
	}

	// generate stage relationMap
	// stageRelationMap := map[stageName]{"stage.From","stage.To"}
	for stageName, stage := range stageMap {
		if _, exist := stageRelationMap[stageName]; !exist {
			stageRelationMap[stageName] = make([]string, 2)
		}
		stageRelationMap[stageName][0] = strings.Join(stage.From, ",")

		for _, from := range stage.From {
			if _, exist := stageRelationMap[from]; !exist {
				stageRelationMap[from] = make([]string, 2)
			}
			if len(stageRelationMap[from][1]) == 0 {
				stageRelationMap[from][1] = stageName
			} else {
				stageRelationMap[from][1] = strings.Join(append(strings.Split(stageRelationMap[from][1], ","), stageName), ",")
			}
		}
	}
	return stageRelationMap, nil
}
