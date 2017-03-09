package models

import (
	"time"
	"errors"
	"github.com/astaxie/beego/orm"
	"go-git-webhook/modules/hash"
)

type WebHook struct {
	WebHookId int			`orm:"pk;auto;unique;column(web_hook_id)" json:"web_hook_id"`
	RepositoryName string		`orm:"size(255);column(repo_name)" json:"repository_name"`
	BranchName string		`orm:"size(255);column(branch_name)" json:"branch_name"`
	ServerId int			`orm:"type(int);column(server_id)" json:"-"`
	Tag string			`orm:"size(1000);column(tag)" json:"tag"`
	Shell string			`orm:"size(1000);column(shell)" json:"shell"`
	Status int			`orm:"type(int);column(status);default(0)" json:"status"`
	Key string			`orm:"size(255);column(key);unique" json:"key"`
	Secure string			`orm:"size(255);column(secure);unique" json:"secure"`
	LastExecTime time.Time		`orm:"type(datetime);column(last_exec_time);null" json:"last_exec_time"`
	CreateTime time.Time		`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
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

func (m *WebHook) Delete() error {
	o := orm.NewOrm()
	_,err := o.Delete(m)

	return err
}
func (m *WebHook) DeleteForServerId(serverId int) error {
	o := orm.NewOrm()

	_,err := o.Raw("DELETE webhooks WHERE server_id = ?",serverId).Exec()

	if err != nil {
		return err
	}
	return nil
}

func (m *WebHook) FindByKey(key string) error {
	o := orm.NewOrm()

	if err := o.QueryTable(m.TableName()).Filter("key",key).One(m);err != nil {
		return err
	}
	return nil
}

func (m *WebHook) Save() error {
	o := orm.NewOrm()
	var err error;

	if m.WebHookId > 0 {
		_,err = o.Update(m)
	}else{
		key := (time.Now().String() + m.RepositoryName + m.BranchName )

		m.Key = hash.Md5(key)

		m.Secure = hash.Md5(key + key + time.Now().String())

		_,err = o.Insert(m)
	}

	return err
}