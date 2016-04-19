package models

import (
	// "fmt"
	// "strconv"
	"time"
)

const (
	FromStageType = 1000
	ToStageType   = 2000
)

// type Stage struct {
// 	Id         int64     `json:"id"`                                        //
// 	PipelineId int64     `json:"pipelineId"`                                //
// 	UUID       string    `json:"uuid" orm:"unique;varchar(255)"`            //
// 	Name       string    `json:"name" orm:"varchar(255)"`                   //
// 	From       string    `json:"from" orm:"null;type(text)"`                //mutil Point.UUID or Stage.UUID
// 	To         string    `json:"to" orm:"null;type(text)"`                  //mutil Point.UUID or Stage.UUID
// 	Content    string    `json:"content" orm:"null;type(text)"`             //
// 	Actived    bool      `json:"actived" orm:"null;default(0)"`             //
// 	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
// 	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
// 	Memo       string    `json:"memo" orm:"null;type(text)"`                //
// }

type Stage struct {
	// gorm.Model
	Id          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time       `sql:"index"`
	WorkspaceId int64            `json:"workspaceId"`
	ProjectId   int64            `json:"projectId"`
	PipelineId  int64            `json:"pipelineId"`
	Name        string           `json:"name"`
	Detail      string           `json:"detail"`
	From        []string         `sql:"-"`
	To          []string         `sql:"-"`
	MetaData    PipelineMetaData `sql:"-"`
	StageSpec   StageSpec        `sql:"-"`
}

type StageVersion struct {
	// gorm.Model
	Id                int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time         `sql:"index"`
	WorkspaceId       int64              `json:"workspaceId"`
	ProjectId         int64              `json:"projectId"`
	PipelineId        int64              `json:"pipelineId"`
	PipelineVersionId int64              `json:"pipelineVersionId"`
	StageId           int64              `json:"stageId"`
	Name              string             `json:"name"`
	Detail            string             `json:"detail"`
	State             *StageVersionState `json:"state" sql:"-"`
	MetaData          PipelineMetaData   `sql:"-"`
	StageSpec         StageSpec          `sql:"-"`
}

type StageRelation struct {
	// gorm.Model
	Id                int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
	WorkspaceId       int64
	ProjectId         int64
	PipelineId        int64
	StageId           int64
	RelationType      uint
	RelationStageName string
}

type StageVersionState struct {
	WorkspaceId       int64
	ProjectId         int64
	PipelineId        int64  `json:"pipelineId"`
	PipelineVersionId int64  `json:"pipelineVersionId"`
	StageId           int64  `json:"pipelineVersionId"`
	StageVersionId    int64  `json:"stageVersionId"`
	StageName         string `json:"stageName"`
	//OK, ERR, TIMEOUT
	RunResult string `json:"runResult"`
	//reserved
	Detail string `json:"detail"`
}

func GetPipelineVersion(pvid int64) *PipelineVersion {

	return &PipelineVersion{}
}

func GetPipeline(pid int64) *Pipeline {

	return &Pipeline{}
}

func GetStageVersion(svid int64) *StageVersion {
	return &StageVersion{}
}

func GetStage(sid int64) *Stage {
	return &Stage{}
}

func (stage *Stage) Save() error {
	db, err := GetDb()
	if err != nil {
		return err
	}

	err = db.Create(stage).Error
	if err != nil {
		return err
	}

	// store stage relation
	for _, fromName := range stage.From {
		tempRelation := new(StageRelation)
		tempRelation.StageId = stage.Id
		tempRelation.WorkspaceId = stage.WorkspaceId
		tempRelation.ProjectId = stage.ProjectId
		tempRelation.PipelineId = stage.PipelineId
		tempRelation.RelationType = FromStageType
		tempRelation.RelationStageName = fromName

		err = db.Create(tempRelation).Error
		if err != nil {
			return err
		}
	}

	for _, toName := range stage.To {
		tempRelation := new(StageRelation)
		tempRelation.StageId = stage.Id
		tempRelation.WorkspaceId = stage.WorkspaceId
		tempRelation.ProjectId = stage.ProjectId
		tempRelation.PipelineId = stage.PipelineId
		tempRelation.RelationType = ToStageType
		tempRelation.RelationStageName = toName

		err = db.Create(tempRelation).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (stageVersion *StageVersion) Save() error {
	db, err := GetDb()
	if err != nil {
		return err
	}

	err = db.Create(stageVersion).Error
	if err != nil {
		return err
	}

	err = db.Create(stageVersion.State).Error
	return err
}

func (stageVersionState StageVersionState) ChangeStageVersionState() error {
	db, err := GetDb()
	if err != nil {
		return err
	}

	err = db.Model(&StageVersionState{}).Where("pipeline_id = ?", stageVersionState.PipelineId).Where("pipeline_version_id = ?", stageVersionState.PipelineVersionId).Where("stage_name = ?", stageVersionState.StageName).Updates(map[string]interface{}{"run_result": stageVersionState.RunResult, "detail": stageVersionState.Detail}).Error
	return err
}

func (stage *Stage) GetStagesByPipelineInfo(pipeline *Pipeline) ([]*Stage, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	result := make([]*Stage, 0)
	err = db.Where("workspace_id = ?", pipeline.WorkspaceId).Where("project_id = ?", pipeline.ProjectId).Where("pipeline_id = ?", pipeline.Id).Find(&result).Error
	if err != nil {
		return nil, err
	}

	// get stage relation
	for k, stage := range result {
		// get all  relation of current stage
		relations := make([]StageRelation, 0)
		from := make([]string, 0)
		to := make([]string, 0)
		db.Where("workspace_id = ?", stage.WorkspaceId).Where("project_id = ?", stage.ProjectId).Where("pipeline_id = ?", stage.PipelineId).Where("stage_id = ?", stage.Id).Find(&relations)
		for _, relation := range relations {
			if relation.RelationType == FromStageType {
				from = append(from, relation.RelationStageName)
			} else if relation.RelationType == ToStageType {
				to = append(to, relation.RelationStageName)
			}
		}
		result[k].From = from
		result[k].To = to

	}

	return result, nil
}

func (stage *Stage) Create(pipelineId int64, name string) (error, string) {
	return nil, ""
}

func (stage *Stage) AddFrom(uuid string, from ...string) error {
	return nil
}

func (stage *Stage) AddEnd(uuid string, end ...string) error {
	return nil
}

func (stage *Stage) Run(uuid string) (error, string) {
	return nil, ""
}

func (stage *Stage) Copy(uuid string) (error, string) {
	return nil, ""
}
