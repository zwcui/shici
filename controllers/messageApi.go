package controllers

import (
	"shici/base"
	"shici/models"
	"shici/util"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

//继承apiController
//消息模块
type MessageController struct {
	apiController
}

//当前api请求之前调用，用于配置哪些接口需要进行head身份验证
func (this *MessageController) Prepare(){
	//this.NeedBaseAuthList = []RequestPathAndMethod{{".+", "post"}, {".+", "patch"}, {".+", "delete"}}
	this.NeedBaseAuthList = []RequestPathAndMethod{}
	this.bathAuth()
}

// @Title 新增MongoDB消息
// @Description 新增MongoDB消息
// @Param	content		formData		string  		true		"内容"
// @Success 200 {object} models.MongoDBMessage
// @Failure 403 create message failed
// @router / [post]
func (this *MessageController) Post() {
	content := this.MustString("content")

	//var user models.User
	//user.NickName = nickName
	//base.DBEngine.Table("user").InsertOne(&user)

	session, mongoDB := base.MongoDB()
	defer session.Close()
	c := mongoDB.C("message")
	err := c.Insert(&models.MongoDBMessage{1, 2, 3, 1, content, util.UnixOfBeijingTime()})
	if err != nil {
		util.Logger.Info("Insert err:"+err.Error())
	}

	result := models.MongoDBMessage{}
	err = c.Find(bson.M{"content":content}).One(&result)
	if err != nil {
		util.Logger.Info("One err:"+err.Error())
	}
	util.Logger.Info("One result:"+strconv.FormatInt(result.MId, 10))
	util.Logger.Info("One result:"+strconv.FormatInt(result.GroupId, 10))
	util.Logger.Info("One result:"+strconv.FormatInt(result.SenderUid, 10))
	util.Logger.Info("One result:"+strconv.Itoa(result.Type))
	util.Logger.Info("One result:"+result.Content)
	util.Logger.Info("One result:"+strconv.FormatInt(result.Created, 10))


	this.ReturnData = result
}

