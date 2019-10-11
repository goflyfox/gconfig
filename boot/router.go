package boot

import (
	"gcs/module/api"
	"gcs/module/common"
	"gcs/module/component/middle"
	"gcs/module/config"
	"gcs/module/constants"
	"gcs/module/system"
	"gcs/utils/base"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"strings"
)

/*
绑定业务路由
*/
func bindRouter() {
	urlPath := g.Config().GetString("url-path")
	s := g.Server()
	// 首页
	s.BindHandler(urlPath+"/", common.Login)
	s.BindHandler(urlPath+"/main.html", common.Index)
	s.BindHandler(urlPath+"/login", common.Login)

	s.BindHandler(urlPath+"/admin/welcome.html", common.Welcome)
	// 中间件
	s.Group(urlPath+"/", func(g *ghttp.RouterGroup) {
		g.Middleware(middle.MiddlewareLog, middle.MiddlewareCommon)
	})

	s.Group(urlPath+"/system", func(g *ghttp.RouterGroup) {
		// 系统路由
		userAction := new(system.UserAction)
		g.ALL("user", userAction)
		g.GET("/user/get/{id}", userAction.Get)
		g.ALL("user/delete/{id}", userAction.Delete)

		departAction := new(system.DepartmentAction)
		g.ALL("department", departAction)
		g.GET("/department/get/{id}", departAction.Get)
		g.ALL("/department/delete/{id}", departAction.Delete)

		logAction := new(system.LogAction)
		g.ALL("log", logAction)
		g.GET("/log/get/{id}", logAction.Get)
		g.ALL("/log/delete/{id}", logAction.Delete)

		menuAction := new(system.MenuAction)
		g.ALL("menu", menuAction)
		g.GET("/menu/get/{id}", menuAction.Get)
		g.ALL("/menu/delete/{id}", menuAction.Delete)

		roleAction := new(system.RoleAction)
		g.ALL("role", roleAction)
		g.GET("/role/get/{id}", roleAction.Get)
		g.ALL("/role/delete/{id}", roleAction.Delete)

		configAction := new(system.ConfigAction)
		g.ALL("config", configAction)
		g.GET("/config/get/{id}", configAction.Get)
		g.ALL("/config/delete/{id}", configAction.Delete)
	})

	s.Group(urlPath+"/admin", func(g *ghttp.RouterGroup) {
		// 项目
		projectAction := new(config.ProjectAction)
		g.ALL(urlPath+"/project", projectAction)
		g.GET(urlPath+"/project/get/{id}", projectAction, "Get")
		g.ALL(urlPath+"/project/delete/{id}", projectAction, "Delete")
		// 发布
		configPublicAction := new(config.ConfigPublicAction)
		g.ALL(urlPath+"/configpublic", configPublicAction)
		g.GET(urlPath+"/configpublic/get/{id}", configPublicAction, "Get")
		g.ALL(urlPath+"/configpublic/delete/{id}", configPublicAction, "Delete")
		g.ALL(urlPath+"/configpublic/rollback/{id}", configPublicAction, "Rollback")
	})

	// 启动gtoken
	base.Token = &gtoken.GfToken{
		//Timeout:         10 * 1000,
		CacheMode:        g.Config().GetInt8("cache-mode"),
		LoginPath:        "/login/submit",
		LoginBeforeFunc:  common.LoginSubmit,
		LogoutPath:       "/user/logout",
		LogoutBeforeFunc: common.LogoutBefore,
		AuthPaths:        g.SliceStr{"/user", "/system", "/admin"},
		AuthBeforeFunc: func(r *ghttp.Request) bool {
			// 静态页面不拦截
			if r.IsFileRequest() {
				return false
			}

			if strings.HasSuffix(r.URL.Path, "index") ||
				strings.HasSuffix(r.URL.Path, ".html") {
				return false
			}

			return true
		},
	}
	base.Token.Start()

	// 对外接口
	s.Group(urlPath+"/config/api", func(g *ghttp.RouterGroup) {
		g.Middleware(middle.MiddlewareApiAuth)

		// 版本和数据接口
		configApiAction := new(api.ConfigApiAction)
		g.ALL("/", configApiAction)
	})

}

/*
统一路由注册
*/
func initRouter() {

	s := g.Server()

	// 绑定路由
	bindRouter()

	if constants.DEBUG {
		g.DB().SetDebug(constants.DEBUG)
	}

	// 上线建议关闭
	s.BindHandler("/debug", common.Debug)

	// 301错误页面
	s.BindStatusHandler(301, common.Error301)
	// 404错误页面
	s.BindStatusHandler(404, common.Error404)
	// 500错误页面
	s.BindStatusHandler(500, common.Error500)

	// 某些浏览器直接请求favicon.ico文件，特别是产生404时
	s.SetRewrite("/favicon.ico", "/resources/images/favicon.ico")

	// 管理接口
	s.EnableAdmin("/administrator")

	// 为平滑重启管理页面设置HTTP Basic账号密码
	//s.BindHookHandler("/admin/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
	//	user := g.Config().GetString("admin.user")
	//	pass := g.Config().GetString("admin.pass")
	//	if !r.BasicAuth(user, pass) {
	//		r.ExitAll()
	//	}
	//})

	// 强制跳转到HTTPS访问
	//g.Server().BindHookHandler("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
	//    if !r.IsFileServe() && r.TLS == nil {
	//        r.Response.RedirectTo(fmt.Sprintf("https://%s%s", r.Host, r.URL.String()))
	//        r.ExitAll()
	//    }
	//})
}
