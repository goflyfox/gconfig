package system

import (
	"gcs/utils"
	"gcs/utils/base"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type MenuAction struct {
	base.BaseRouter
}

var (
	actionNameMenu = "MenuAction"
)

// path: /index
func (action *MenuAction) Index(r *ghttp.Request) {
	tplFile := "pages/system/menu_index.html"
	err := r.Response.WriteTpl(tplFile, g.Map{
		"now": gtime.Datetime(),
	})

	if err != nil {
		glog.Error(err)
	}
}

// path: /get/{id}
func (action *MenuAction) Get(r *ghttp.Request) {
	id := r.GetInt("id")
	model := SysMenu{Id: id}.Get()
	if model.Id <= 0 {
		base.Fail(r, actionNameMenu+" get fail")
	}

	base.Succ(r, model)
}

// path: /delete/{id}
func (action *MenuAction) Delete(r *ghttp.Request) {
	id := r.GetInt("id")

	model := SysMenu{Id: id}
	model.UpdateId = base.GetUser(r).Id
	model.UpdateTime = utils.GetNow()

	num := model.Delete()
	if num <= 0 {
		base.Fail(r, actionNameMenu+" delete fail")
	}

	base.Succ(r, "")
}

// path: /save
func (action *MenuAction) Save(r *ghttp.Request) {
	model := SysMenu{}
	err := gconv.Struct(r.GetPostMap(), &model)
	if err != nil {
		glog.Error(actionNameMenu+" save struct error", err)
		base.Error(r, "save error")
	}

	userId := base.GetUser(r).Id

	model.UpdateId = userId
	model.UpdateTime = utils.GetNow()

	var num int64
	if model.Id <= 0 {
		model.CreateId = userId
		model.CreateTime = utils.GetNow()
		num = model.Insert()
	} else {
		num = model.Update()
	}

	if num <= 0 {
		base.Fail(r, actionNameMenu+" save fail")
	}

	base.Succ(r, "")
}

// path: /list
func (action *MenuAction) List(r *ghttp.Request) {
	form := base.NewForm(r.GetPostMap())
	model := SysMenu{}

	list := model.List(&form)
	base.Succ(r, list)
}

// path: /tree
func (action *MenuAction) Tree(r *ghttp.Request) {
	form := base.NewForm(r.GetPostMap())
	model := SysMenu{}

	list := model.List(&form)
	base.Succ(r, list)
}

// path: /page
func (action *MenuAction) Page(r *ghttp.Request) {
	form := base.NewForm(r.GetPostMap())
	model := SysMenu{}

	page := model.Page(&form)
	base.Succ(r, g.Map{"list": page, "form": form})
}

// path: /jqgrid
func (action *MenuAction) Jqgrid(r *ghttp.Request) {
	form := base.NewForm(r.GetPostMap())
	model := SysMenu{}

	page := model.Page(&form)
	r.Response.WriteJson(g.Map{
		"page":    form.Page,
		"rows":    page,
		"total":   form.TotalPage,
		"records": form.TotalSize,
	})
}
