package db

import (
	"strings"
	"time"
)

type Tb_etcd_backup struct {
	Id             int64  `orm:"pk;size(11)"`
	Key            string `orm:"size(255)"`
	Parent_key     string `orm:"size(255)"`
	Value          string `orm:"size(255)"`
	Dir            int    `orm:"size(1)"`
	Ttl            int64  `orm:"size(11)"`
	Modified_index uint64 `orm:"size(11)"`
	Created_index  uint64 `orm:"size(11)"`
}

func (this *Tb_etcd_backup) Update() (int64, error) {
	return Orm.Update(this)
}

// server detection  and delete expire key
func (this *Tb_etcd_backup) DeleteExpire() (bool, int64, error) {
	// To be realized
	return true, 0, nil
}

func (this *Tb_etcd_backup) InertOrUpdate() (bool, int64, error) {
	index := strings.LastIndex(this.Key, "/")
	this.Parent_key = string([]rune(this.Key)[0:index])
	if this.Parent_key == "" {
		this.Parent_key = "/"
	}

	newObj := &Tb_etcd_backup{}
	newObj.Key = this.Key
	newObj.Value = this.Value
	newObj.Dir = this.Dir
	newObj.Ttl = this.Ttl
	newObj.Modified_index = this.Modified_index
	newObj.Created_index = this.Created_index
	newObj.Parent_key = this.Parent_key

	create, id, err := Orm.ReadOrCreate(this, "Key")
	if err == nil && !create {
		if newObj.Modified_index > this.Modified_index {
			newObj.Id = this.Id
			id, err = Orm.Update(newObj)
		}
	}
	return create, id, err
}

func (this *Tb_etcd_backup) Delete() (int64, error) {
	return Orm.Delete(this)
}

func (this *Tb_etcd_backup) Read() error {
	return Orm.Read(this)
}
func (this *Tb_etcd_backup) Insert() (int64, error) {
	return Orm.Insert(this)
}
func (this *Tb_etcd_backup) Exist() bool {
	return Orm.QueryTable("Tb_etcd_backup").Filter("Key", this.Key).Exist()
}
func (this *Tb_etcd_backup) IsDirectory() bool {
	return this.Dir == 1
}

func (this *Tb_etcd_backup) IsExpired() bool {
	if this.Ttl > 0 {
		if time.Now().Unix()-this.Ttl > 0 {
			return true
		}
	}
	return false
}
