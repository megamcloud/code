package helpers

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/actions"
)

type UserMustAuth struct {
	Username string
}

func (this *UserMustAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	var action = actionPtr.Object()
	var session = action.Session()
	var username = session.GetString("username")
	if len(username) == 0 {
		this.login(action)
		return false
	}

	this.Username = username

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	// 初始化变量
	action.Data["teaTitle"] = "TeaWeb管理平台"
	action.Data["teaUsername"] = username
	action.Data["teaMenu"] = ""
	action.Data["teaModules"] = []map[string]interface{}{
		{
			"code":     "proxy",
			"menuName": "代理设置",
			"icon":     "paper plane outline",
		},
		{
			"code":     "log",
			"menuName": "日志",
			"icon":     "history",
		},
		{
			"code":     "stat",
			"menuName": "统计",
			"icon":     "chart area",
		},
		/**{
			"code":     "services",
			"menuName": "服务",
			"icon":     "gem outline",
		},**/
		{
			"code":     "plugins",
			"menuName": "插件",
			"icon":     "puzzle piece",
		},
		/**{
			"code":     "team",
			"menuName": "团队",
			"icon":     "users",
		},**/
		/**{
			"code":     "lab",
			"menuName": "实验室",
			"icon":     "medapps",
		},**/
	}
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}
	action.Data["teaVersion"] = teaconst.TeaVersion

	return true
}

func (this *UserMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/login")
}
