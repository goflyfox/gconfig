package bean

type SessionUser struct {
	Id       int    `form:"id" json:"id"`               // 主键
	Uuid     string `form:"uuid" json:"uuid"`           // UUID
	UserType int    `form:"userType" json:"user_type"`  // 用户类型
	Username string `form:"username" json:"username"`   // 登录名/11111
	RealName string `form:"real_name" json:"real_name"` // 真实姓名
}
