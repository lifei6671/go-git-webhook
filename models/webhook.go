package models

import "time"

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
