package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
)

type Server struct {
	ServerId int			`orm:"pk;auto;unique;column(server_id)" json:"server_id"`
	Name string			`orm:"size(255);column(name)" json:"name"`
	Type string			`orm:"size(255);column(type);default(ssh)" json:"type"`
	IpAddress string		`orm:"size(255);column(ip_address)" json:"ip_address"`
	Port int			`orm:"type(int);column(port)" json:"port"`
	Account string			`orm:"size(255);column(account)" json:"account"`
	PrivateKey string		`orm:"type(text);column(private_key)" json:"private_key"`
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

//创建或更新
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

//删除
func (m *Server) Delete() error {
	o := orm.NewOrm()
	_,err := o.Delete(m)

	return err
}


func (m *Server) Search(keyword string, memberId int) ([]Server,error) {
	o := orm.NewOrm()

	keyword = "%" + keyword + "%"

	var servers []Server

	_,err := o.Raw("SELECT * FROM servers WHERE create_at = ? AND (name LIKE ? OR servers.tag LIKE ?)",memberId,keyword,keyword).QueryRows(&servers)

	if err != nil {
		fmt.Println(err)
		return servers,err
	}

	return servers,nil
}