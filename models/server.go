package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Server struct {
	ServerId int			`orm:"pk;auto;unique;column(server_id)" json:"server_id"`
	Name string			`orm:"size(255);column(name);not null" json:"name"`
	Type string			`orm:"size(255);column(type);not null;default(ssh)" json:"type"`
	IpAddress string		`orm:"size(255);column(ip_address);not null" json:"ip_address"`
	Port int			`orm:"type(int);column(port);not null" json:"port"`
	Account string			`orm:"size(255);column(account);not null" json:"account"`
	PrivateKey string		`orm:"type(text);column(private_key);not null" json:"private_key"`
	Tag string			`orm:"size(1000);column(tag)" json:"tag"`
	Status int 			`orm:"type(int);column(status);default(0)" json:"status"`
	CreateTime time.Time		`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt int			`orm:"type(int);column(create_at)" json:"-"`
}


func (m *Server) TableName() string {
	return "servers"
}

func (m *Server) TableEngine() string {
	return "INNODB"
}

func NewServer() *Server {
	return &Server{}
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

func (m *Server) Save() error {
	o := orm.NewOrm()
	var err error;

	if m.ServerId > 0 {
		_,err = o.Update(m)
	}else{
		_,err = o.Insert(m)
	}

	return err
}
func (m *Server) Delete() error {
	o := orm.NewOrm()
	_,err := o.Delete(m)

	return err
}