package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"hackathon/initDB"
	"hackathon/model"
	"hackathon/util"
	"log"
	"net/http"
)


type Verifyphone struct {
	gorm.Model
	Phone string `json:"phone" binding:"required,len=11"`
}
//手机号验证
func VerifyPhone(c *gin.Context) {
	//获取需要验证的手机号
	var verifyPhone Verifyphone
	err := c.BindJSON(&verifyPhone)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{"errcode":400, "Msg": "Post Data Err"})
		return
	}

	//判断手机号是否已经被注册
	if util.IsTelephoneExist(initDB.Db, verifyPhone.Phone) {
		c.JSON(http.StatusOK, gin.H{"status": "wphone"})
		return
	}

	//发送验证短信
	captcha := util.GenValidateCode(6)
	 util.SendSms(verifyPhone.Phone, captcha)


	//返回数据
	c.JSON(http.StatusOK, gin.H{
		"status": "sms_success",
		"Msg": captcha,
	})
}


type BindingUser struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}
//注册
func Register(c *gin.Context) {
	//获取用户信息
	var userMsg BindingUser
	err := c.BindJSON(&userMsg)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"errcode":400,
			"Msg": "Post Data Err",
		})
		return
	}

	//将信息写入数据库
	newUser := model.User{
		Model:    gorm.Model{},
		UserNickname: userMsg.Nickname,
		UserPhone:    userMsg.Phone,
		UserPassword: userMsg.Password,
	}
	initDB.Db.Create(&newUser)

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"status":0,
		"Msg": "注册成功",
	})
}
