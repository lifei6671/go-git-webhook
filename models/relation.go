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
	MemberId int		`orm:"type(int);column(member_id)" json:"member_id"`
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

// 包含 WebHook 和 Server 信息的关系实体
type RelationDetailed struct {
	RelationId int			`json:"relation_id"`
	MemberId int			`json:"member_id"`
	WebHookId int			`json:"web_hook_id" orm:"column(web_hook_id)"`
	RepositoryName string		`json:"repository_name"`
	BranchName string		`json:"branch_name"`
	ServerId int			`json:"server_id"`
	WebHookTag string		`json:"web_hook_tag"`
	Shell string			`json:"shell"`
	WebHookStatus int		`json:"web_hook_status"`
	Key string			`json:"key"`
	Secure string			`json:"secure"`

	ServerName string		`json:"server_name"`
	ServerType string		`json:"server_type"`
	IpAddress string		`json:"ip_address"`
	Port int			`json:"port"`
	Account string			`json:"account"`
	PrivateKey string		`json:"-"`
	ServerTag string                `json:"server_tag"`
	ServerStatus int                `json:"server_status"`

}

func NewRelationDetailed() *RelationDetailed {
	return &RelationDetailed{}
}

func FindRelationDetailed(relationId int) (RelationDetailed,error) {
	var relationDetailed RelationDetailed
	if relationId <= 0 {
		return  relationDetailed,ErrInvalidParameter
	}

	relation := &Relation{ RelationId: relationId}

	o := orm.NewOrm()
	if err := o.Read(relation);err != nil {
		return relationDetailed,err
	}


	server := &Server{ ServerId: relation.ServerId}

	if err := o.Read(server);err != nil {
		return relationDetailed,err
	}

	hook := &WebHook{WebHookId: relation.WebHookId}

	if err := o.Read(hook);err != nil {
		return relationDetailed,err
	}


	relationDetailed.RelationId 	= relationId
	relationDetailed.MemberId	= relation.MemberId
	relationDetailed.WebHookId 	= relation.WebHookId
	relationDetailed.ServerId 	= relation.ServerId

	relationDetailed.RepositoryName	= hook.RepositoryName
	relationDetailed.BranchName	= hook.BranchName
	relationDetailed.WebHookTag	= hook.Tag
	relationDetailed.Shell		= hook.Shell
	relationDetailed.WebHookStatus	= hook.Status
	relationDetailed.Key		= hook.Key
	relationDetailed.Secure		= hook.Secure

	relationDetailed.ServerName	= server.Name
	relationDetailed.ServerType	= server.Type
	relationDetailed.IpAddress	= server.IpAddress
	relationDetailed.Port		= server.Port
	relationDetailed.Account	= server.Account
	relationDetailed.PrivateKey	= server.PrivateKey
	relationDetailed.ServerTag	= server.Tag
	relationDetailed.ServerStatus	= server.Status

	return relationDetailed,nil
}

// 指定条件查询完整的关系对象
func FindRelationDetailedByWhere(where string,params ...interface{}) ([]RelationDetailed,error) {
	o := orm.NewOrm()

	sql := "SELECT relation_id,member_id,relation.web_hook_id,server.server_id,name AS server_name,type, ip_address, port, account, private_key,server.tag as server_tag,server.status AS server_status, repo_name AS repository_name,branch_name,hook.tag AS web_hook_tag,shell, hook.status AS web_hook_status,hook.key,secure FROM relations AS relation " +
		"LEFT JOIN servers AS server ON relation.server_id = server.server_id " +
		"LEFT JOIN webhooks as hook ON relation.web_hook_id = hook.web_hook_id WHERE 1=1 ";

	if where != "" {
		sql += where
	}

	rawSetter := o.Raw(sql,params)

	var results []RelationDetailed

	_,err := rawSetter.QueryRows(&results)

	return results,err
}

type ServerRelation struct {
	ServerId int
	RelationId int
	WebHookId int
	MemberId int
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

	sql := "SELECT servers.*,relations.member_id,relations.create_time,web_hook_id,relation_id FROM relations LEFT JOIN servers ON relations.server_id = servers.server_id WHERE web_hook_id = ? AND create_at = ? ORDER BY relation_id DESC "

	_,err := o.Raw(sql,webHookId,memberId).QueryRows(&res)


	return res,err
}