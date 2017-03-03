package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Server struct {
	ServerId int			`orm:"pk;auto;unique;column(server_id)"`
	Name string			`orm:"size(255);column(name);not null"`
	Type string			`orm:"size(255);column(type);not null;default(ssh)"`
	IpAddress string		`orm:"size(255);column(ip_address);not null"`
	Port int			`orm:"type(int);column(port);not null"`
	Account string			`orm:"size(255);column(account);not null"`
	PrivateKey string		`orm:"type(text);column(private_key);not null"`
	Tag string			`orm:"size(1000);column(tag)"`
	Status int 			`orm:"type(int);column(status);default(0)"`
	CreateTime time.Time		`orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt int			`orm:"type(int);column(create_at)"`
}


func (m *Server) TableName() string {
	return "servers"
}

func (m *Server) TableEngine() string {
	return "INNODB"
}

//根据ID查找对象
func (m *Server) Find(id int) (error) {
	o := orm.NewOrm()

	m.ServerId = id

	if err := o.Read(m) ;err != nil {
		return err
	}
	return nil;
}

//分页查询服务器
func (m *Server) GetServerList(pageIndex int,pageSize int) ([]Server,int64,error) {
	o := orm.NewOrm()
	offset  := (pageIndex -1) * pageSize
	var list []Server

	_,err := o.QueryTable(m.TableName()).Limit(pageSize).Offset(offset).OrderBy("-member_id").All(&list)

	if err != nil {
		if err == orm.ErrNoRows {
			return list,0,nil
		}
		return list,0,err
	}
	count,err := o.QueryTable(m.TableName()).Count()

	if err != nil {
		return list,0,err
	}
	return list,count,nil
}