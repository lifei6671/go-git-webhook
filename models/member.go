package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"go-git-webhook/modules/passwords"
	"fmt"
)

type Member struct {
	MemberId int		`orm:"pk;auto;unique;column(member_id)"`
	Account string 		`orm:"size(255);column(account)"`
	Password string 	`orm:"size(1000);column(password)"`
	Email string 		`orm:"size(255);column(email)"`
	Phone string 		`orm:"size(255);column(phone)"`
	Avatar string 		`orm:"size(1000);column(avatar)"`
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt int		`orm:"type(int);column(create_at)"`
	LastLoginTime time.Time	`orm:"type(datetime);column(last_login_time);null"`
}

func (m *Member) TableName() string {
	return "members"
}

func (m *Member) TableEngine() string {
	return "INNODB"
}


func NewMember() *Member {
	return new(Member)
}

//根据用户ID查找用户
func (m *Member) Find(id int) (error) {
	o := orm.NewOrm()
	m.MemberId = id

	err := o.Read(m)

	if err == orm.ErrNoRows {
		return ErrMemberNoExist
	}

	return nil
}

//用户登录
func (m *Member) Login(account string,password string) (*Member,error) {
	o := orm.NewOrm()

	member := &Member{}

	err := o.QueryTable(m.TableName()).Filter("account",account).One(member);

	if err != nil {
		return  member,ErrMemberNoExist
	}

	ok,err := passwords.PasswordVerify(member.Password,password) ;

	fmt.Println(err)
	fmt.Println(ok)

	if ok && err == nil {
		return member,nil
	}

	return member,ErrorMemberPasswordError
}

//分页获取用户列表
func (m *Member) GetMemberList(pageIndex int,pageSize int) ([]Member,int64,error) {

	offset  := (pageIndex -1) * pageSize

	o := orm.NewOrm()

	var members []Member

	_,err := o.QueryTable(m.TableName()).Limit(pageSize).Offset(offset).OrderBy("-member_id").All(&members)

	if err != nil {
		return nil,0,err
	}

	count,err := o.QueryTable(m.TableName()).Count()

	return members,count,nil
}

//添加一个用户
func (member *Member) Add () (error) {
	o := orm.NewOrm()

	hash ,err := passwords.PasswordHash(member.Password);

	if  err != nil {
		return err
	}

	member.Password = hash

	_,err = o.Insert(member)

	if err != nil {
		return err
	}
	return  nil
}

//更新用户信息
func (m *Member) Update(cols... string) (error) {
	o := orm.NewOrm()

	if _,err := o.Update(m,cols...);err != nil {
		return err
	}
	return nil
}