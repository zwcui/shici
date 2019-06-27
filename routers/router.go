// @APIVersion 1.0.0
// @Title shici 接口基础工程
// @Description 接入Redis,MongoDB,xorm,seelog等，加入部署脚本
// @Contact zwcui2017@163.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"shici/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/message",
			beego.NSInclude(
				&controllers.MessageController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
