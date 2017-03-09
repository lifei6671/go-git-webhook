package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"errors"
)

type Relation struct {
	RelationId int		`orm:"pk;auto;unique;column(relation_id)" json:"relation_id"`
	WebHookId int		`orm:"type(int);column(web_hook_id)" json:"web_hook_id"`
	ServerId int		`orm:"type(int);column(server_id)" json:"server_id"`
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"` //添加时间
}


func (m *Relation) TableName() string {
	return "relations"
}

func (m *Relation) TableEngine() string {
	return "INNODB"
}

func (m *Relation) Save () error {

	o := orm.NewOrm()

	var err error

	if m.RelationId > 0 {
		if m.WebHookId <= 0 || m.ServerId <= 0 {
			return errors.New("Data format error")
		}
		if o.QueryTable(m.TableName()).Filter("web_hook_id",m.WebHookId).Filter("server_id",m.ServerId).Exist() {
			return nil
		}
		_,err = o.Update(m)
	}else{
		_,err =o.Insert(m)
	}
	return err
}