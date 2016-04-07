package models

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
	Id          int64    `json:"id"`
	WorkspaceId int64    `json:"workspaceId"`
	ProjectId   int64    `json:"projectId"`
	PipelineId  int64    `json:"pipelineId"`
	Created     int64    `json:"created"`
	Updated     int64    `json:"updated"`
	Name        string   `json:"name"`
	Detail      string   `json:"detail"`
	Dependences []string `json:"dependences"`
}

type StageVersion struct {
	Id                int64    `json:"id"`
	PipelineId        int64    `json:"pipelineId"`
	PipelineVersionId int64    `json:"pipelineVersionId"`
	StageId           int64    `json:"stageId"`
	Created           int64    `json:"created"`
	Updated           int64    `json:"updated"`
	Name              string   `json:"name"`
	Detail            string   `json:"detail"`
	Dependences       []string `json:"dependences"`
	State             int64    `json:"state"` // 0 not start    1 working    2 success    3 failed
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
