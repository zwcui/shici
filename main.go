package main

import (
	_ "baseApi/routers"
	"github.com/astaxie/beego"
	_ "baseApi/models"
	_ "baseApi/base"
	_ "baseApi/task"
	_ "baseApi/controllers"
	"baseApi/util"
	"github.com/astaxie/beego/plugins/cors"
)


/*
	1. govendor add +e 将项目使用到但为加入vendor的包加入工程，参考vendor.json
	2. xorm 同步结构体与表结构，默认驼峰  其他参考： https://www.kancloud.cn/kancloud/xorm-manual-zh-cn/56013
	3. 使用swagger 第一次执行  bee run -gendoc=true -downdoc=true；目前只能在swagger目录下的index.html中配置默认head
	4. 如果启动工程时想执行其他包中的init()方法，则引入进来，前面加 "_ "
	5. 为了对老版本的兼容，docker部署多个容器，由请求head中的api-version区分，通过nginx进行不同端口的跳转；可能同一套数据库多个服务会产生新版本更新时老版本仍需更新，配置主从数据库？
	6. 域名解析时，注意接口前缀如api，网页前缀如www
	7. 考虑加入gRPC，远程调用，方便创建分布式应用
	8. controller.go中，初始化后请求前调用Prepare()，请求后调用Finish()等等
	9. 集成了seelog日志
	10.指定路径git init;git clone XXXX.git;cd XXX;git checkout -b dev;git branch;
       deploy_dev.sh 用于部署开发服务器docker，取git最新程序同步至服务器；
	   deploy_test.sh用于部署测试服务器docker，取git最新程序同步至服务器；
       deploy_prod.sh用于部署正式服务器docker，取git最新程序同步至服务器；
	11.服务器go build编译时出现signal: killed，多半是服务器内存不够，free -h查看，
	   然后
		手动释放 echo 1 > /proc/sys/vm/drop_caches
		或者
		增加虚拟内存
			dd if=/dev/vda1 of=swapfile bs=1024000 count=2000
			mkswap swapfile
			swapon swapfile
			chmod 600 swapfile
			swapoff -v swapfile	停用虚拟内存
			vim /etc/fstab	; 最后加一行 /dev/swapfile            swap                    swap     defaults       0 0		开机启动虚拟内存
	12.加入定时任务，建议定时任务的设计原则是可以反复跑，以便宕机补数据；可以通过数据库记录跑任务的时间与状态，然后只跑未执行的任务
	13.接入websocket，ws://ip:8079/ws
	14.TODO:xorm支持的主从数据库读写分离
	15.TODO:接入微信、支付宝 支付

 */
func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// 跨域
	beego.InsertFilter("/v1/*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Api-Version", "AppNo", "Source"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	defer util.Logger.Flush()

	beego.Run()
}
