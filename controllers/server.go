package controllers

type ServerController struct {
	BaseController
}

func (c *ServerController) Index() {
	c.Layout = ""
	c.TplName = "server/index.html"


}
