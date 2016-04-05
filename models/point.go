package models

const (
	POINTSTART = iota
	POINTEND
	POINTCHECK
)

// type Point struct {
// 	Id         int64     `json:"id"`                                        //
// 	PipelineId int64     `json:"pipelineId"`                                //
// 	Type       int       `json:"type"`                                      //POINTSTART, POINTEND, POINTCHECK
// 	UUID       string    `json:"uuid" orm:"unique;varchar(255)"`            //
// 	Name       string    `json:"name" orm:"varchar(255)"`                   //
// 	From       string    `json:"from" orm:"null;type(text)"`                //mutil Point.UUID or Stage.UUID
// 	To         string    `json:"to" orm:"null;type(text)"`                  //mutil Point.UUID or Stage.UUID
// 	Actived    bool      `json:"actived" orm:"null;default(0)"`             //
// 	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
// 	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
// 	Memo       string    `json:"memo" orm:"null;type(text)"`                //
// }

type Point struct {
	Id         int64  `json:"id"`
	PipelineId int64  `json:"pipelineId"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	Name       string `json:"name"`
	From       string `json:"from"`
	To         string `json:"to"`
}

type PointVersion struct {
	Id         int64  `json:"id"`
	PipelineId int64  `json:"pipelineId"`
	PointId    int64  `json:"pointId"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	Name       string `json:"name"`
	From       string `json:"from"`
	To         string `json:"to"`
	State      int64  `json:"state"` // 0 waitting    1 finish
}

//
// func (point *Point) Create(pipelineId int64, pointType int, name string) (error, string) {
// 	o := orm.NewOrm()
// 	p := Point{PipelineId: pipelineId, Type: pointType, UUID: uuid.NewV4().String(), Name: name, Actived: true}
//
// 	if err := o.Begin(); err != nil {
// 		log.Errorf("Transcation error: %s", err.Error())
//
// 		return err, ""
// 	} else {
// 		if id, err := o.Insert(&p); err != nil {
// 			log.Errorf("Create point error: %s", err.Error())
//
// 			o.Rollback()
// 			return err, ""
// 		} else {
// 			log.Infof("Create point successfully, id is: %d", id)
//
// 			o.Commit()
// 			return nil, p.UUID
// 		}
// 	}
// }
//
// func (point *Point) AddFrom(uuid string, from ...string) error {
// 	o := orm.NewOrm()
// 	p := Point{UUID: uuid, Actived: true}
//
// 	if err := o.Read(&p, "UUID", "Actived"); err != nil {
// 		log.Errorf("Get point %s error: %s", uuid, err.Error())
//
// 		return err
// 	} else {
// 		if err := o.Begin(); err != nil {
// 			log.Errorf("Transcation error: %s", err.Error())
//
// 			o.Rollback()
// 			return err
// 		} else {
// 			new_from, _ := underscore.Uniq(strings.Split(p.From+";"+strings.Join(from, ";"), ";"), nil)
// 			p.From = strings.Join(new_from.([]string), ";")
//
// 			if _, err := o.Update(&p, "From"); err != nil {
// 				log.Errorf("Put point from %s error %s:", uuid, err.Error())
//
// 				o.Rollback()
// 				return err
// 			} else {
// 				log.Infof("Put point from %s successfully", uuid)
//
// 				o.Commit()
// 				return nil
// 			}
// 		}
// 	}
// }
//
// func (point *Point) RemoveFrom(uuid string, from string) error {
// 	return nil
// }
//
// func (point *Point) AddTo(uuid string, to ...string) error {
// 	o := orm.NewOrm()
// 	p := Point{UUID: uuid, Actived: true}
//
// 	if err := o.Read(&p, "UUID", "Actived"); err != nil {
// 		log.Errorf("Get point %s error: %s", uuid, err.Error())
//
// 		return err
// 	} else {
// 		if err := o.Begin(); err != nil {
// 			log.Errorf("Transcation error: %s", err.Error())
//
// 			o.Rollback()
// 			return err
// 		} else {
// 			new_to, _ := underscore.Uniq(strings.Split(p.To+";"+strings.Join(to, ";"), ";"), nil)
// 			p.To = strings.Join(new_to.([]string), ";")
//
// 			if _, err := o.Update(&p, "To"); err != nil {
// 				log.Errorf("Put point end %s error %s:", uuid, err.Error())
//
// 				o.Rollback()
// 				return err
// 			} else {
// 				log.Infof("Put point end %s successfully", uuid)
//
// 				o.Commit()
// 				return nil
// 			}
// 		}
// 	}
// }
//
// func (point *Point) RemoveTo(uuid string, to string) error {
// 	return nil
// }
//
// func (point *Point) Run(uuid string) (error, string) {
// 	return nil, ""
// }
//
// func (point *Point) Copy(uuid string) (error, string) {
// 	return nil, ""
// }
