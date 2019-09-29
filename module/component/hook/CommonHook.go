package hook

import (
	"fmt"
	"gcs/module/constants"
	"gcs/utils/base"
	"gcs/utils/resp"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

func CommonBefore(r *ghttp.Request) {
	r.SetParam("BASE_PATH", "")
}

func LogBeforeServe(r *ghttp.Request) {
	if !constants.DEBUG {
		return
	}

	now := gtime.Millisecond()
	r.SetParam("_now", now)

	if r.IsFileRequest() {
		return
	}

	var params map[string]interface{}
	if r.Method == "GET" {
		params = r.GetQueryMap()
	} else if r.Method == "POST" {
		params = r.GetPostMap()
	} else {
		base.Error(r, "Request Method is ERROR! ")
		return
	}

	no := gconv.String(params["no"])
	if no == "" {
		no = gconv.String(now)
	}

	glog.Info(fmt.Sprintf("[REQUEST_%s_%d][url:%s][params:%s]",
		no, r.Id, r.URL.Path, params))

}

func LogBeforeOutput(r *ghttp.Request) {
	if !constants.DEBUG {
		return
	}

	data := string(r.Response.Buffer())
	if r.URL.Path == "" || r.URL.Path == "/" || gstr.Contains(
		r.URL.Path, "index") || gstr.Contains(
		r.URL.Path, "html") || gstr.Contains(
		r.URL.Path, "login") {
		data = ""
	}
	var params map[string]interface{}
	if r.Method == "GET" {
		params = r.GetQueryMap()
	} else if r.Method == "POST" {
		params = r.GetPostMap()
	} else {
		r.Response.Writeln("Request Method is ERROR! ")
		return
	}

	now := gtime.Millisecond()
	rTime := gconv.Int64(r.GetParam("_now"))
	no := gconv.String(params["no"])
	if no == "" {
		if rTime == 0 {
			no = gconv.String(gtime.Millisecond())
		} else {
			no = gconv.String(rTime)
		}
	}

	if r.IsFileRequest() {
		glog.Info(fmt.Sprintf("[FILE_%s_%d][diff:%d][url:%s][params:%s]",
			no, r.Id, now-rTime, r.URL.Path, params))
	} else if now-rTime > 1000 {
		glog.Warning(fmt.Sprintf("[RESPONSE_%s_%d][diff:%d][url:%s][params:%s][data:%s]",
			no, r.Id, now-rTime, r.URL.Path, params, data))
	} else {
		glog.Info(fmt.Sprintf("[RESPONSE_%s_%d][diff:%d][url:%s][params:%s][data:%s]",
			no, r.Id, now-rTime, r.URL.Path, params, data))
	}

}

func AuthAfterFunc(r *ghttp.Request, respData resp.Resp) {
	if !respData.Success() {
		var params map[string]interface{}
		if r.Method == "GET" {
			params = r.GetQueryMap()
		} else if r.Method == "POST" {
			params = r.GetPostMap()
		} else {
			r.Response.Writeln("Request Method is ERROR! ")
			return
		}

		no := gconv.String(gtime.Millisecond())

		glog.Info(fmt.Sprintf("[AUTH_%s_%d][url:%s][params:%s][data:%s]",
			no, r.Id, r.URL.Path, params, respData.Json()))
		respData.Msg = "请求错误或登录超时"
		r.Response.WriteJson(respData)
		r.ExitAll()
	}
}

func AuthBeforeFunc(r *ghttp.Request) bool {
	// 静态页面不拦截
	if r.IsFileRequest() {
		return false
	}
	// 静态页面不拦截
	if r.IsAjaxRequest() {
		return true
	}

	return false
}
