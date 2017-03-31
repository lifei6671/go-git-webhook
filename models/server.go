package models

import (
	"time"
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

// 服务器对象
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

// 获取对应数据库表名
func (m *Server) TableName() string {
	return "servers"
}
// 获取数据使用的引擎
func (m *Server) TableEngine() string {
	return "INNODB"
}
// 新建服务器对象
func NewServer() *Server {
	return &Server{}
}

// 根据ID查找对象
func (m *Server) Find() (error) {

	if m.ServerId <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	if err := o.Read(m) ;err != nil {
		return err
	}
	return nil;
}

// 创建或更新
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

// 删除
func (m *Server) Delete() error {
	o := orm.NewOrm()
	_,err := o.Delete(m)

	return err
}

// 搜索指定用户的服务器
func (m *Server) Search(keyword string, memberId int,excludeServerId ...int) ([]Server,error) {
	o := orm.NewOrm()

	keyword = "%" + keyword + "%"

	sql := "SELECT * FROM servers WHERE create_at = ? AND (name LIKE ? OR servers.tag LIKE ?)"

	if len(excludeServerId) > 0 {
		sql += " AND server_id not in ("
		for _,num := range excludeServerId {
			sql += strconv.Itoa(num) + ","
		}
		sql += "0)"
	}

	var servers []Server

	_,err := o.Raw(sql,memberId,keyword,keyword).QueryRows(&servers)

	if err != nil {
		logs.Error("",err.Error())
		return servers,err
	}

	return servers,nil
}

// 根据server_id和用户id查询服务器信息列表
func (m *Server) QueryServerByServerId(serverIds []int,memberId ...int) ([]*Server,error) {
	o := orm.NewOrm()

	query := o.QueryTable(m.TableName()).Filter("server_id__in",serverIds)

	if len(memberId) > 0 {
		query = query.Filter("create_at",memberId[0])
	}

	var servers []*Server

	_,err := query.All(&servers)

	return servers,err
}