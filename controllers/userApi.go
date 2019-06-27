package controllers

import (
	"shici/base"
	"shici/models"
	"shici/util"
)

//继承apiController
//用户模块
type UserController struct {
	apiController
}

//当前api请求之前调用，用于配置哪些接口需要进行head身份验证
func (this *UserController) Prepare(){
	//this.NeedBaseAuthList = []RequestPathAndMethod{{".+", "post"}, {".+", "patch"}, {".+", "delete"}}
	this.NeedBaseAuthList = []RequestPathAndMethod{}
	this.bathAuth()
	util.Logger.Info("UserController beforeRequest ")
}

// @Title 新增用户
// @Description create users
// @Param	nickName		formData		string  		true		"昵称"
// @Success 200 {string} success
// @Failure 403 create users failed
// @router / [post]
func (this *UserController) Post() {
	nickName := this.MustString("nickName")

	var user models.User
	user.NickName = nickName
	base.DBEngine.Table("user").InsertOne(&user)

	this.ReturnData = "success"
}
//
//// @Title CreateUser
//// @Description create users
//// @Param	body		body 	models.User	true		"body for user content"
//// @Success 200 {int} models.User.Id
//// @Failure 403 body is empty
//// @router / [post]
//func (u *UserController) Post() {
//	var user models.User
//	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
//	uid := models.AddUser(user)
//	u.Data["json"] = map[string]string{"uid": uid}
//	u.ServeJSON()
//}
//
//// @Title GetAll
//// @Description get all Users
//// @Success 200 {object} models.User
//// @router / [get]
//func (u *UserController) GetAll() {
//	users := models.GetAllUsers()
//	u.Data["json"] = users
//	u.ServeJSON()
//}
//
//// @Title Get
//// @Description get user by uid
//// @Param	uid		path 	string	true		"The key for staticblock"
//// @Success 200 {object} models.User
//// @Failure 403 :uid is empty
//// @router /:uid [get]
//func (u *UserController) Get() {
//	uid := u.GetString(":uid")
//	if uid != "" {
//		user, err := models.GetUser(uid)
//		if err != nil {
//			u.Data["json"] = err.Error()
//		} else {
//			u.Data["json"] = user
//		}
//	}
//	u.ServeJSON()
//}
//
//// @Title Update
//// @Description update the user
//// @Param	uid		path 	string	true		"The uid you want to update"
//// @Param	body		body 	models.User	true		"body for user content"
//// @Success 200 {object} models.User
//// @Failure 403 :uid is not int
//// @router /:uid [put]
//func (u *UserController) Put() {
//	uid := u.GetString(":uid")
//	if uid != "" {
//		var user models.User
//		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
//		uu, err := models.UpdateUser(uid, &user)
//		if err != nil {
//			u.Data["json"] = err.Error()
//		} else {
//			u.Data["json"] = uu
//		}
//	}
//	u.ServeJSON()
//}
//
//// @Title Delete
//// @Description delete the user
//// @Param	uid		path 	string	true		"The uid you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 uid is empty
//// @router /:uid [delete]
//func (u *UserController) Delete() {
//	uid := u.GetString(":uid")
//	models.DeleteUser(uid)
//	u.Data["json"] = "delete success!"
//	u.ServeJSON()
//}
//
//// @Title Login
//// @Description Logs user into the system
//// @Param	username		query 	string	true		"The username for login"
//// @Param	password		query 	string	true		"The password for login"
//// @Success 200 {string} login success
//// @Failure 403 user not exist
//// @router /login [get]
//func (u *UserController) Login() {
//	username := u.GetString("username")
//	password := u.GetString("password")
//	if models.Login(username, password) {
//		u.Data["json"] = "login success"
//	} else {
//		u.Data["json"] = "user not exist"
//	}
//	u.ServeJSON()
//}
//
//// @Title logout
//// @Description Logs out current logged in user session
//// @Success 200 {string} logout success
//// @router /logout [get]
//func (u *UserController) Logout() {
//	u.Data["json"] = "logout success"
//	u.ServeJSON()
//}



//通过手机号获取用户信息
func UserWithPhoneNumber(phoneNumber string) (user *models.User, err error) {
	hasUser, err := base.DBEngine.Table("user").Where("phone_number=?", phoneNumber).Get(&user)
	if err != nil {
		return nil, err
	}
	if !hasUser {
		return nil, nil
	}
	return user, nil
}