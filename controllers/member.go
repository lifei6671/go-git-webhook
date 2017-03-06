package controllers

type MemberController struct {
	BaseController
}

func (c *MemberController) Index() {
	c.Layout = ""
	c.TplName = "member/index.html"
}


func (c *MemberController) My(){
	c.Layout = ""
	c.TplName = "member/my.html"
}