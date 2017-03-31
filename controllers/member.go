package controllers

import (
	"path/filepath"
	"strings"
	"time"
	"os"
	"strconv"
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"image/png"
	"image/gif"

	"github.com/lifei6671/go-git-webhook/modules/pager"
	"github.com/lifei6671/go-git-webhook/modules/passwords"
	"github.com/lifei6671/go-git-webhook/models"

	"github.com/astaxie/beego/logs"
)

// 会员控制器
type MemberController struct {
	BaseController
}

// 首页
func (c *MemberController) Index() {
	c.Prepare()

	if c.Member.Role != 0 {
		c.Abort("403")
	}

	c.Layout = ""
	c.TplName = "member/list.html"
	c.Data["MemberSelected"] = true

	pageIndex, _ := c.GetInt("page", 1)

	var members []models.Member

	pageOptions := pager.PageOptions{
		TableName:  models.NewMember().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		PageSize : 15,
		ParamName : "page",
		Conditions : " order by member_id desc",
	}



	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err := rs.QueryRows(&members)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("",err.Error())
	}

	c.Data["lists"] = members
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
}

// 个人中心
func (c *MemberController) My(){
	c.Prepare()
	c.Layout = ""
	c.TplName = "member/edit.html"
	c.Data["MemberSelected"] = true

	member := c.Member

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		status, _ := c.GetInt("status", 0)

		if password != "" {
			pass, _ := passwords.PasswordHash(password)
			member.Password = pass
		}
		if member.Role != 0 {
			member.Status = status
		}

		member.Email = c.GetString("email")
		member.Phone = c.GetString("phone")
		member.Avatar = c.GetString("avatar")

		if member.Avatar == "" {
			member.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if member.MemberId > 0 {

			result = member.Update()
		} else {
			result = member.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}

		view, err := c.ExecuteViewPathTemplate("member/list_item.html", *member)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view" : view,
		}
		c.SetMember(*member)
		c.JsonResult(0, "ok", data)

	}

	c.Data["Model"] = member
	c.Data["IsSelf"] = true
}

// 编辑信息
func (c *MemberController) Edit() {
	c.TplName = "member/edit.html"

	c.Prepare()


	member_id ,_ := c.GetInt(":id")

	member := models.NewMember()

	if member_id > 0 {
		member.MemberId = member_id
		if err := member.Find(); err != nil {
			c.ServerError("Data query error:" + err.Error())
		}
	}

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		account := c.GetString("account")
		status ,_ := c.GetInt("status",0)

		if member.MemberId > 0 {
			if password != "" {
				pass, _ := passwords.PasswordHash(password)
				member.Password = pass
			}
			if member.Role != 0 {
				member.Status = status
			}
		} else {
			if account == ""{
				c.JsonResult(500,"Account is require.")
			}
			if password == "" {
				c.JsonResult(500,"Password is require.")
			}
			member.Role = 1
			member.Account = account
			member.Password = password
		}

		member.Email = c.GetString("email")
		member.Phone = c.GetString("phone")
		member.Avatar = c.GetString("avatar")

		if member.Avatar == "" {
			member.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if member.MemberId > 0 {

			result = member.Update()
		} else {
			result = member.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}

		view, err := c.ExecuteViewPathTemplate("member/list_item.html", *member)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view" : view,
		}
		c.JsonResult(0, "ok", data)

	}


	c.Data["Model"] = member
	c.Data["IsSelf"] = false
	if member.MemberId == c.Member.MemberId {
		c.Data["IsSelf"] = true
	}
}

// 删除会员
func (c *MemberController) Delete() {
	c.Prepare()

	if c.Member.Role != 0 {
		c.Abort("403")
	}

	member_id ,err := c.GetInt(":id")

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"Parameter error.")
	}

	member := models.NewMember()
	member.MemberId = member_id

	if err := member.Find();err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"Data query error.")
	}

	if member.Role == 0 {
		c.JsonResult(500,"不能删除管理员用户")
	}
	if err := member.Delete();err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"删除失败")
	}
	c.JsonResult(0,"ok")
}

// 上传图片
func (c *MemberController) Upload() {
	file,moreFile,err := c.GetFile("image-file")
	defer file.Close()

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"读取文件异常")
	}

	ext := filepath.Ext(moreFile.Filename);
	if !strings.EqualFold(ext,".png") && !strings.EqualFold(ext,".jpg") && !strings.EqualFold(ext,".gif") && !strings.EqualFold(ext,".jpeg")  {
		c.JsonResult(500,"不支持的图片格式")
	}

	x ,_ := c.GetInt("x")
	y ,_ := c.GetInt("y")
	width ,_ := c.GetInt("width")
	height ,_ := c.GetInt("height")

	fileName := "avatar_" +  strconv.FormatInt(int64(time.Now().Nanosecond()), 16) + ext

	filePath := "static/uploads/" + time.Now().Format("200601") + "/" + fileName

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file",filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}

	fileBytes,err := ioutil.ReadFile(filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}
	buf := bytes.NewBuffer(fileBytes)

	m,_,_ := image.Decode(buf)

	rgbImg := m.(*image.YCbCr)

	subImg := rgbImg.SubImage(image.Rect(x, y, x + width, y + height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil{
		c.JsonResult(500,"保存图片失败")
	}
	defer f.Close()

	if strings.EqualFold(ext,".jpg") || strings.EqualFold(ext,".jpeg"){
		err = jpeg.Encode(f,subImg,&jpeg.Options{ Quality : 100 })
	}else if strings.EqualFold(ext,".png") {
		err = png.Encode(f,subImg)
	}else if strings.EqualFold(ext,".gif") {
		err = gif.Encode(f,subImg,&gif.Options{ NumColors : 256})
	}
	if err != nil {
		c.JsonResult(500,"图片剪切失败")
	}

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"保存文件失败")
	}
	url := "/" + filePath

	c.JsonResult(0,"ok",url)
}