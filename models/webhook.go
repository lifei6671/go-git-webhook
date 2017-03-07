package models

import (
	"time"
	"errors"
	"github.com/astaxie/beego/orm"
	"crypto/md5"
)

type WebHook struct {
	WebHookId int			`orm:"pk;auto;unique;column(web_hook_id)"`
	RepositoryName string		`orm:"size(255);column(repo_name);not null"`
	BranchName string		`orm:"size(255);column(branch_name);not null"`
	ServerId int			`orm:"type(int);column(server_id);not null"`
	Tag string			`orm:"size(1000);column(tag)"`
	Shell string			`orm:"size(1000);column(shell)"`
	Status int			`orm:"type(int);column(status);default(0)"`
	Key string			`orm:"size(1000);column(key)"`
	LastExecTime time.Time		`orm:"type(datetime);column(last_exec_time)"`
	CreateTime time.Time		`orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt int			`orm:"type(int);column(create_at)"`
}


func (m *WebHook) TableName() string {
	return "webhooks"
}

func (m *WebHook) TableEngine() string {
	return "INNODB"
}
func NewWebHook() *WebHook {
	return &WebHook{}
}

func (m *WebHook) Find(id int) error {
	o := orm.NewOrm()

	m.WebHookId = id

	if err := o.Read(m) ;err != nil {
		return err
	}
	return nil;
}

func (m *WebHook) DeleteMulti (id ...int) error {
	if len(id) > 0 {
		o := orm.NewOrm()
		ids := make([]int,len(id))
		params := ""

		for i := 0;i<len(id);i++ {
			ids[i] = id[i]
			params = params + ",?"
		}
		_,err := o.Raw("DELETE webhooks WHERE web_hook_id IN ("+ params[1:] +")",ids).Exec()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Invalid parameter")
}

func (m *WebHook) DeleteForServerId(serverId int) error {
	o := orm.NewOrm()

	_,err := o.Raw("DELETE webhooks WHERE server_id = ?",serverId).Exec()

	if err != nil {
		return err
	}
	return nil
}

func (m *WebHook) Save() {
	o := orm.NewOrm()
	var err error;

	if m.ServerId > 0 {
		_,err = o.Update(m)
	}else{
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(time.Now().String() + m.RepositoryName + m.BranchName ))
		cipherStr := md5Ctx.Sum(nil)

		m.Key = string(cipherStr)

		_,err = o.Insert(m)
	}

	return err
}