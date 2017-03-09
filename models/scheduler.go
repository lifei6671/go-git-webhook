package models

import "time"

//任务调度器储存表
type Scheduler struct {
	SchedulerId int 		`orm:"pk;auto;unique;column(scheduler_id)" json:"scheduler_id"`
	WebHookId int			`orm:"type(int);column(web_hook_id)" json:"web_hook_id"`
	ServerId int			`orm:"type(int);column(server_id)" json:"server_id"`
	Status string			`orm:"column(status);default(wait)" json:"status"` //状态：wait 等待执行/executing 执行中/suspend 中断执行/ failure 执行失败/ success 执行成功
	CreateTime time.Time		`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"` //添加时间
	StartExecTime time.Time		`orm:"type(datetime);column(start_exec_time);null" json:"start_exec_time"` //开始执行时间
	EndExecTime time.Time		`orm:"type(datetime);column(end_exec_time);null" json:"end_exec_time"` //执行结束时间
	Data string			`orm:"type(text);column(data);not null" json:"data"`	//推送的数据
	PushUser string			`orm:"column(push_user);null;size(255)" json:"push_user"` //推送用户
	LogContent string		`orm:"column(log_content);type(text);" json:"log_content"`
}


func (m *Scheduler) TableName() string {
	return "scheduler"
}

func (m *Scheduler) TableEngine() string {
	return "INNODB"
}

