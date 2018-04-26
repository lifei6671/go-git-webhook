package models

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

// Scheduler 任务调度器储存表
type Scheduler struct {
	SchedulerId   int       `orm:"pk;auto;unique;column(scheduler_id)" json:"scheduler_id"`
	WebHookId     int       `orm:"type(int);column(web_hook_id)" json:"web_hook_id"`
	ServerId      int       `orm:"type(int);column(server_id)" json:"server_id"`
	RelationId    int       `orm:"type(int);column(relation_id)" json:"relation_id"`
	Status        string    `orm:"column(status);default(wait)" json:"status"`                         //状态：wait 等待执行/executing 执行中/suspend 中断执行/ failure 执行失败/ success 执行成功
	CreateTime    time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"` //添加时间
	StartExecTime time.Time `orm:"type(datetime);column(start_exec_time);null" json:"start_exec_time"` //开始执行时间
	EndExecTime   time.Time `orm:"type(datetime);column(end_exec_time);null" json:"end_exec_time"`     //执行结束时间
	Data          string    `orm:"type(text);column(data)" json:"-"`                                   //推送的数据
	PushUser      string    `orm:"column(push_user);null;size(255)" json:"push_user"`                  //推送用户
	ShaValue      string    `orm:"column(sha_value);null;size(255)" json:"sha_value"`                  //当前请求的sha值
	LogContent    string    `orm:"column(log_content);type(text);" json:"log_content"`
	ExecuteType   int       `orm:"column(execute_type);type(int);default(0)" json:"execute_type"` //执行方式：0 自动触发 / 1 手动执行
}

// WebScheduler 前端使用的结构体
type WebScheduler struct {
	Scheduler
	Slogan    string `json:"slogan"`
	Consuming string `json:"consuming"`
}

// ToWebScheduler 转换为前端可用的优化信息对象
func (m *Scheduler) ToWebScheduler() WebScheduler {
	item := WebScheduler{}

	item.SchedulerId = m.SchedulerId
	item.WebHookId = m.WebHookId
	item.ServerId = m.ServerId
	item.RelationId = m.RelationId
	item.Status = m.Status
	item.CreateTime = m.CreateTime
	item.StartExecTime = m.StartExecTime
	item.EndExecTime = m.EndExecTime
	item.Data = m.Data
	item.PushUser = m.PushUser
	item.ShaValue = m.ShaValue
	item.LogContent = m.LogContent
	item.ExecuteType = m.ExecuteType
	item.Consuming = ""

	duration := time.Now().Sub(m.CreateTime)

	if year := int(duration.Seconds() / (365 * 24 * 3600)); year >= 1 {
		item.Slogan = strconv.Itoa(year) + "年前"
	} else if month := int(duration.Seconds() / (30 * 24 * 3600)); month >= 1 {
		item.Slogan = strconv.Itoa(month) + "月前"
	} else if day := int(duration.Seconds() / (24 * 3600)); day >= 1 {
		item.Slogan = strconv.Itoa(day) + "天前"
	} else if hour := int(duration.Hours()); hour >= 1 {
		item.Slogan = strconv.Itoa(hour) + "小时前"
	} else if minute := int(duration.Minutes()); minute >= 1 {
		item.Slogan = strconv.Itoa(minute) + "分钟前"
	} else {
		item.Slogan = "刚刚"
	}

	if m.StartExecTime.Year() > 2000 && m.EndExecTime.Year() > 2000 {
		sub := m.EndExecTime.Sub(m.StartExecTime)

		if sub.Hours() > 1 {
			item.Consuming += strconv.Itoa(int(sub.Hours())) + "时"
		}
		if sub.Minutes() > 1 {
			item.Consuming += strconv.Itoa(int(sub.Minutes())) + "分"
		}
		if sub.Seconds() > 1 {
			item.Consuming += strconv.Itoa(int(sub.Seconds())) + "秒"
		}
		millisecond := sub.Nanoseconds() / 1000000

		if item.Consuming == "" && millisecond > 1 {
			item.Consuming = strconv.Itoa(int(millisecond)) + "毫秒"
		}
	}
	return item
}

// TableName 获取对应数据库表名
func (m *Scheduler) TableName() string {
	return "scheduler"
}

// TableEngine 获取数据使用的引擎
func (m *Scheduler) TableEngine() string {
	return "INNODB"
}

// NewScheduler 新建对象
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// Find 根据ID查找对象
func (m *Scheduler) Find() error {
	if m.SchedulerId <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	return o.Read(m)
}

// InsertMulti 批量插入对象
func (m *Scheduler) InsertMulti(schedulers []Scheduler) (int64, error) {
	if len(schedulers) <= 0 {
		return 0, ErrInvalidParameter
	}
	o := orm.NewOrm()

	return o.InsertMulti(len(schedulers), schedulers)
}

// QuerySchedulerByState 根据状态查询
func (m *Scheduler) QuerySchedulerByState(state ...string) ([]Scheduler, error) {
	o := orm.NewOrm()

	var results []Scheduler

	_, err := o.QueryTable(m.TableName()).Filter("status__in", state).All(&results)
	return results, err
}

// Save 更新或插入
func (m *Scheduler) Save() error {
	o := orm.NewOrm()
	var err error
	if m.SchedulerId > 0 {
		_, err = o.Update(m)
	} else {
		_, err = o.Insert(m)
	}
	return err
}

// DeleteByWhere 根据条件删除
func (m *Scheduler) DeleteByWhere(where string, args ...interface{}) error {
	o := orm.NewOrm()

	sql := "DELETE FROM scheduler WHERE 1=1 " + where

	_, err := o.Raw(sql, args).Exec()

	return err
}

//获取当前最大SchedulerId
func (m *Scheduler) GetMaxSchedulerId() int {
	o := orm.NewOrm()
	sql := "select * FROM scheduler order by scheduler_id desc limit 1"

	var results []Scheduler
	_, err := o.Raw(sql).QueryRows(&results)
	if err != nil {
		return 0
	}
	for _, result := range results {
		return result.SchedulerId

	}
	return 0
}
