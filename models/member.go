package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"go-git-webhook/modules/passwords"
)

type Member struct {
	MemberId int		`orm:"pk;auto;unique;column(member_id)"`
	Account string 		`orm:"size(255);column(account);not null"`
	Password string 	`orm:"size(1000);column(password)"`
	Email string 		`orm:"size(255);column(email)"`
	Phone string 		`orm:"size(255);column(phone)"`
	Avatar string 		`orm:"size(1000);column(avatar)"`
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt int		`orm:"type(int);column(create_at)"`
	LastLoginTime time.Time	`orm:"type(datetime);column(last_login_time)"`
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
func (m *Member) Find(id int) (*Member,error) {
	o := orm.NewOrm()
	user := Member{MemberId: id}

	err := o.Read(&user)

	if err == orm.ErrNoRows {
		return nil,ErrMemberNoExist
	}
	return &user,nil
}


func (m *Member) Login(account string,password string) (*Member,error) {
	o := orm.NewOrm()

	member := &Member{}

	err := o.QueryTable(m.TableName()).Filter("account",account).One(member);

	if err != nil {
		return  member,ErrMemberNoExist
	}

	if ok,err := passwords.PasswordVerify(member.Password,password) ; ok && err == nil {
		return member,nil
	}

	return member,ErrorMemberPasswordError
}