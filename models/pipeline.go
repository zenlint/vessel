package models

import (
	"time"
)

const (
	PIPESUCCESS = iota
	PIPEERROR
	PIPERUNNING
	PIPEPENDDING
	PIPEEXCEPT
)

var (
	DEFAULT_PIPELINE_ETCD_PATH        = "/containerops/vessel/ws-%d/pj-%d/pl-%d/stage/"
	DEFAULT_PIPELINEVERSION_ETCD_PATH = "/containerops/vessel/ws-%d/pj-%d/pl-%d/version/plv-%d"
)

type Pipeline struct {
	// gorm.Model
	Id          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time       `sql:"index"`
	WorkspaceId int64            `json:"workspaceId"`
	ProjectId   int64            `json:"projectId"`
	Name        string           `json:"name" gorm:"type:varchar(255)"`
	SelfLink    string           `json:"selfLink" gorm:"type:varchar(255)"`
	Labels      string           `json:"labels"`
	Annotations string           `json:"annotations"`
	Detail      string           `json:"detail" gorm:"type:text"`
	Stages      []*Stage         `sql:"-"`
	MetaData    PipelineMetaData `sql:"-"`
	StageSpecs  []StageSpec      `sql:"-"`
}

type PipelineVersion struct {
	// gorm.Model
	Id            int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time       `sql:"index"`
	WorkspaceId   int64            `json:"workspaceId"`
	ProjectId     int64            `json:"projectId"`
	PipelineId    int64            `json:"pipelineId"`
	Namespace     string           `json:"namespace"`
	SelfLink      string           `json:"selfLink" gorm:"type:varchar(255)"`
	Labels        string           `json:"labels"`
	Annotations   string           `json:"annotations"`
	Detail        string           `json:"detail" gorm:"type:text"`
	StageVersions string           `json:"stageVersions"`
	Log           string           `json:"log" gorm:"type:text"`
	Status        string           `json:"state"`
	MetaData      PipelineMetaData `sql:"-"`
	StageSpecs    []StageSpec      `sql:"-"`
}

func (pv *PipelineVersion) GetMetadata() PipelineMetaData {
	return pv.MetaData
}

func (pv *PipelineVersion) GetSpec() []StageSpec {
	return pv.StageSpecs
}

// type Status struct {
// 	Id         int64     `json:"id"`                                        //
// 	PipelineId int64     `json:"pipelineId"`                                //
// 	UUID       string    `json:"uuid" orm:"varchar(255)"`                   //
// 	Resource   string    `json:"resrouce" orm:"null;type(text)"`            //
// 	ActionId   string    `json:"actionId" orm:"unique;varchar(255)"`        // Point.UUID; Stage.UUID
// 	Started    time.Time `json:"started" orm:"type(datetime)"`              //
// 	Ended      time.Time `json:"ended" orm:"type(datetime)"`                //
// 	Log        string    `json:"log" orm:"type(text)"`                      //
// 	Result     int64     `json:"result" orm:"null"`                         // Success: 0; Failure: 1
// 	Actived    bool      `json:"actived" orm:"null;default(0)"`             //
// 	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
// 	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
// 	Memo       string    `json:"memo" orm:"null;type(text)"`                //
// }

// save pipeline info
func (pipe *Pipeline) Save() (int64, error) {
	db, err := GetDb()
	if err != nil {
		return 0, err
	}

	err = db.Create(pipe).Error
	return pipe.Id, err
}

// save pipelineVersion info
func (plv *PipelineVersion) Save() error {
	db, err := GetDb()
	if err != nil {
		return err
	}

	err = db.Create(plv).Error
	return err
}

// get pipeline info by pipeline id
func (pipe *Pipeline) GetPipelineInfoByPipelineId(id int64) (*Pipeline, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	result := new(Pipeline)
	err = db.Where("id = ?", id).First(result).Error

	return result, err
}

func (plv *PipelineVersion) Done() error {
	db, err := GetDb()
	if err != nil {
		return err
	}
	return db.Model(plv).Updates(map[string]interface{}{"Status": "Done"}).Error
}

//Create Pipeline
func (pipe *Pipeline) Create(projectId int64, name string) (error, int64) {
	return nil, 0
}

//List all run history with pipelineId
func (pipe *Pipeline) History(pipelineId int64) (error, []string) {
	return nil, []string{}
}

//Return all point and stage status with status uuid
func (pipe *Pipeline) Status(uuid string) (error, []int64) {
	return nil, []int64{}
}

//Stop run
func (pipe *Pipeline) Stop(uuid string) error {
	return nil
}

//Clean run resources
func (pipe *Pipeline) Clean(uuid string) error {
	return nil
}

func (pipe *Pipeline) Copy(uuid string) (error, string) {
	return nil, ""
}
