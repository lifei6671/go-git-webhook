package models

import "time"

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
