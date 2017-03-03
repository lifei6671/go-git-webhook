package models

import "time"

type History struct {
	Id int			`orm:"pk;auto;unique;column(server_id)"`
	Status int		`orm:"type(int);column(status);default(0)"`
	ShellLog string		`orm:"type(text);column(shell_log);not null"`
	Data string		`orm:"type(text);column(data);not null"`
	PushUser string		`orm:"type(text);column(push_user);not null"`
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add"`
	UpdateTime time.Time	`orm:"type(datetime);column(update_time)"`
	WebHookId int		`orm:"type(int);column(web_hook_id)"`
}


func (m *History) TableName() string {
	return "historys"
}

func (m *History) TableEngine() string {
	return "INNODB"
}
