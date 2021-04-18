package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"login_service/common"
	"login_service/models"
)

func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func AppAccessHandler(c *gin.Context) {
	var f ParamsApp
	if err := c.ShouldBindBodyWith(&f, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	c.JSON(200, AppCheckResponse(f))
}

func UserLoginHandler(c *gin.Context) {
	var param LoginParam
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}

	db, err := models.Connect()
	if err != nil {
		c.JSON(200, gin.H{"test": "t"})
		return
	}
	var user models.User
	db.Where(&models.User{UserName: param.UserName, Password: common.GetMD5Hash(param.Password)}).First(&user)
	s, _ := convertParamToDict(user)
	c.JSON(200, s)
}

func CreateUser(c *gin.Context) {
	var param CreateUserParam
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(200, gin.H{"test": "t"})
		return
	}
	user := models.User{
		UserName: param.UserName,
		Password: common.GetMD5Hash(param.Password),
		MobileNo: param.MobileNo,
	}
	result := db.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("user created, %v", result.Error)})
	} else {
		c.JSON(201, gin.H{"message": fmt.Sprintf("user created, %v", user.UserName)})

	}
}
