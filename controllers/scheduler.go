package controllers

import (
	"go-git-webhook/models"
	"go-git-webhook/modules/pager"
	"strconv"
	"fmt"
	"github.com/astaxie/beego/logs"
)

type SchedulerController struct {
	BaseController
}

func (c *SchedulerController) Index()  {
	c.Prepare()
	c.TplName = "scheduler/index.html"


	pageIndex, _ := c.GetInt("page", 1)
	relation_id,_ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if relation_id <= 0 {
		c.Abort("404")
	}

	relationDetailedResult ,err := models.FindRelationDetailedByWhere("AND relation_id = ?",relation_id)

	if err != nil {
		fmt.Printf("FindRelationDetailed Error : %s", err.Error())
		c.Abort("500")
	}
	var relationDetailed models.RelationDetailed

	if len(relationDetailedResult) > 0{
		relationDetailed = relationDetailedResult[0]
	}

	var schedulers []models.Scheduler

	pageOptions := pager.PageOptions{
		TableName:  models.NewScheduler().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		ParamName : "page",
		Conditions : " AND relation_id = " + strconv.Itoa(relationDetailed.RelationId) + " order by scheduler_id desc",
	}

	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err = rs.QueryRows(&schedulers)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("",err.Error())
	}

	c.Data["Model"] = relationDetailed

	c.Data["lists"] = schedulers
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
}

func (c *SchedulerController) Console() {

}

func (c *SchedulerController) Cancel() {
	schedulerId,err := strconv.Atoi(c.Ctx.Input.Param(":scheduler_id"))

	if err != nil {
		c.JsonResult(500,"Parameter error")
	}

	scheduler := models.NewScheduler()
	scheduler.SchedulerId = schedulerId

	if err := scheduler.Find();err != nil {
		c.JsonResult(500,"Error 50001: Query data error")
	}
	deailed,err := models.FindRelationDetailedByWhere("AND relation_id = ? AND member_id = ?", scheduler.RelationId,c.Member.MemberId)

	if err != nil || len(deailed) <= 0{
		c.JsonResult(404,"The data does not exist")
	}

	scheduler.Status = "suspend"

	if err := scheduler.Save(); err != nil {
		c.JsonResult(500,"Cancel failed")
	}
	c.JsonResult(0,"ok")
}

func (c *SchedulerController) Resume () {

}
