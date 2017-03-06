package controllers

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Prepare()

	c.Layout = ""
	c.TplName = "home/index.html"
}