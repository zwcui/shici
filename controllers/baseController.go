package controllers

import (
	"github.com/astaxie/beego"
	"regexp"
	"strings"
	_ "github.com/astaxie/beego/cache/redis"
	"shici/util"
	"time"
	"shici/base"
	"fmt"
	"errors"
)

//数据返回结构体
type BaseController struct {
	beego.Controller
	Err        error
	ErrCode    int
	ReturnData interface{}
}

//api接口统一controller
type apiController struct {
	BaseController
	NeedBaseAuthList []RequestPathAndMethod
}

//需要验证的请求路径
type RequestPathAndMethod struct {
	PathRegexp string
	Method     string
}

const (
	REDIS_BATHAUTH = "BaseAuth_"
)

//接口调用返回code
const (
	ServerApiSuccuess 		= 1000		//调用成功
	ServerApiUndefinedFail 	= 999		//未知错误
	ServerApiIllegalParam 	= 900		//接口参数不合法

)

//返回数据的head
const (
	headerCodeKey 	= "code"
	headerDesKey 	= "description"
)

const (
	DEFAULT_PAGESIZE = 15	//默认分页15条
)


//默认请求之前加路径head身份验证，默认所有方法都需要验证，各个api可以重写该方法
func (this *apiController) Prepare(){
	util.Logger.Info("apiController Prepare")
	this.NeedBaseAuthList = []RequestPathAndMethod{{".+", "post"}, {".+", "patch"}, {".+", "delete"}, {".+", "put"}}
	this.bathAuth()
}

//对路径进行校验
func (this *apiController) bathAuth(){
	pathNeedAuth := false
	for _, value := range this.NeedBaseAuthList {
		if ok, _ := regexp.MatchString(value.PathRegexp, this.Ctx.Request.URL.Path); ok && strings.ToUpper(this.Ctx.Request.Method) == strings.ToUpper(value.Method) {
			pathNeedAuth = true
			break
		}
	}

	//要求head中放Authorization，内容格式为 "Basic 18800000000:123456"  密码为加密后的密文，加密方式为base64
	if pathNeedAuth {
		phoneNumber, encryptedPassword, ok := this.Ctx.Request.BasicAuth()
		if !ok {
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"empty auth"+`"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			this.ServeJSON()
			this.StopRun()
		}
		redisTemp := base.RedisCache.Get(REDIS_BATHAUTH+phoneNumber)
		if redisTemp == nil {
			user, err := UserWithPhoneNumber(phoneNumber)
			if err != nil || user == nil {
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+err.Error()+`"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			//校验密码
			passwordByte := util.Base64Encode([]byte(encryptedPassword))
			password := string(passwordByte)
			hashedPwd, _ := util.EncryptPasswordWithSalt(password, user.Salt)
			if hashedPwd != user.Password {
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"password error"+`"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}

			//存入redis
			if user != nil {
				base.RedisCache.Put(REDIS_BATHAUTH+phoneNumber, encryptedPassword, 60*60*2*time.Second)
			}
		} else {
			if encryptedPassword != redisTemp {
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"password redis error"+`"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
		}

	}
}

//func (this *BaseController) Init() {
//	util.Logger.Info("BaseController Init")
//}

func (this *BaseController) Prepare() {
	util.Logger.Info("BaseController Prepare")
}

func (this *BaseController) Get() {
	util.Logger.Info("BaseController Get")
}

func (this *BaseController) Post() {
	util.Logger.Info("BaseController Post")
}

func (this *BaseController) Delete() {
	util.Logger.Info("BaseController Delete")
}

func (this *BaseController) Put() {
	util.Logger.Info("BaseController Put")
}

func (this *BaseController) Head() {
	util.Logger.Info("BaseController Head")
}

func (this *BaseController) Patch() {
	util.Logger.Info("BaseController Patch")
}

func (this *BaseController) Options() {
	util.Logger.Info("BaseController Options")
}

//func (this *BaseController) Render() {
//	util.Logger.Info("BaseController Render")
//}

//func (this *BaseController) XSRFToken() {
//	util.Logger.Info("BaseController XSRFToken")
//}

//func (this *BaseController) CheckXSRFCookie() {
//	util.Logger.Info("BaseController CheckXSRFCookie")
//}

//func (this *BaseController) HandlerFunc() {
//	util.Logger.Info("BaseController HandlerFunc")
//}

func (this *BaseController) URLMapping() {
	util.Logger.Info("BaseController URLMapping")
}

//取参数错误返回
func (this *BaseController) Failed() {
	util.Logger.Info("BaseController Failed")
	if this.ErrCode == 0 {
		this.ErrCode = ServerApiUndefinedFail
	}
	this.Data["json"] = map[string]interface{}{
		"header": map[string]string{
			headerCodeKey: fmt.Sprintf("%d", this.ErrCode),
			headerDesKey:  this.Err.Error(),
		},
	}
	this.ServeJSON()
	this.StopRun()
}

// 函数结束时,组装成json结果返回
func (this *BaseController) Finish() {
	util.Logger.Info("BaseController Finish")
	if this.Err != nil {
		this.Failed()
	}
	r := struct {
		Header interface{} `json:"header"`
		Data   interface{} `json:"data"`
	}{}

	r.Header = map[string]string{
		headerCodeKey: fmt.Sprintf("%d", ServerApiSuccuess),
		headerDesKey:  "success",
	}

	r.Data = this.ReturnData
	this.Data["json"] = r
	this.ServeJSON()
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustString(key string) string {
	v := this.GetString(key)
	if v == "" {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustInt64(key string) int64 {
	v, err := this.GetInt64(key)
	if err != nil {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustFloat64(key string) float64 {
	v, err := this.GetFloat(key)
	if err != nil {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustInt(key string) int {
	v, err := this.GetInt(key)
	if err != nil {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

func (this *BaseController) GetPageSize(key string) int {
	v, _ := this.GetInt(key)
	if v == 0 {
		return DEFAULT_PAGESIZE
	}
	return v
}