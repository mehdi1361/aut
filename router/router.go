package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"login_service/common"
	"login_service/models"
	"time"
)

func PingHandler(c *gin.Context) {
	//r, _ := AppCache.Get("ping")

	c.JSON(200, gin.H{
		"message": AppCache.KeyContain("il"),
	})

}

func AppAccessHandler(c *gin.Context) {
	var f ParamsApp
	if err := c.ShouldBindBodyWith(&f, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}

	result, _ := AppCache.Get(fmt.Sprintf("%s-%s", f.ClientId, f.SecretKey))
	if len(result) == 0 {
		data := AppCheckResponse(f)
		_ = AppCache.Set(fmt.Sprintf("%s-%s", f.ClientId, f.SecretKey), data, 10*time.Minute)
		c.JSON(200, data)

	} else {
		r := make(map[string]interface{})
		_ = json.Unmarshal([]byte(string(result)), &r)
		c.JSON(200, r)
	}

}

func UserLoginHandler(c *gin.Context) {
	var param LoginParam
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}

	statusCode, result, err := UserLogin(param)
	if err != nil {
		c.JSON(statusCode, err)
	} else {
		c.JSON(statusCode, result)
	}
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
		UserId:   common.UuidGenerator(),
		MobileNo: param.MobileNo,
	}
	result := db.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("user created, %v", result.Error)})
	} else {
		c.JSON(201, gin.H{"message": fmt.Sprintf("user created, %v", user.UserName)})
	}
}

func CheckPermission(c *gin.Context) {
	token := c.Request.Header["Token"][0]
	d, _ := AppCache.Get(token)
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(string(d)), &result)
	c.JSON(200, result)
}

type Role struct{}

func (r *Role) POST(c *gin.Context) {
	db, err := models.Connect()
	if err != nil {
		c.JSON(200, gin.H{"message": "error in connect to database"})
		return
	}
	role := &models.Role{
		Name:   "test",
		FaName: "test",
	}
	result := db.Create(&role)

	if result.Error != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("user created, %v", result.Error)})
	} else {
		c.JSON(201, gin.H{"message": fmt.Sprintf("user created, %v", role.Name)})
	}

}
