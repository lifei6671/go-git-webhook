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

func NewRelation() *Relation {
	return &Relation{}
}

func (m *Relation) Save () error {

	o := orm.NewOrm()
	if o.QueryTable(m.TableName()).Filter("web_hook_id",m.WebHookId).Filter("server_id",m.ServerId).Exist() {
		return ErrServerAlreadyExist
	}
	var err error

	if m.RelationId > 0 {
		if m.WebHookId <= 0 || m.ServerId <= 0 {
			return errors.New("Data format error")
		}
		_,err = o.Update(m)
	}else{
		_,err =o.Insert(m)
	}
	return err
}

func (m *Relation) Delete()error {
	o := orm.NewOrm()
	_,err := o.Delete(m)

	return err
}

func (m *Relation) Find(id int) error {
	o := orm.NewOrm()

	m.RelationId = id

	if err := o.Read(m) ;err != nil {
		return err
	}
	return nil;
}

type ServerRelation struct {
	ServerId int
	RelationId int
	WebHookId int
	Status int
	Name string
	IpAddress string
	Port int
	Type string
	CreateTime time.Time
	CreateAt int
}

func (m *Relation) QueryByWebHookId (webHookId int,memberId int) ( []*ServerRelation ,error){
	o := orm.NewOrm()

	var res []*ServerRelation

	sql := "SELECT servers.*,relations.create_time,web_hook_id,relation_id FROM relations LEFT JOIN servers ON relations.server_id = servers.server_id WHERE web_hook_id = ? AND create_at = ? ORDER BY relation_id DESC "

	_,err := o.Raw(sql,webHookId,memberId).QueryRows(&res)


	return res,err
}