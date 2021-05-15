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

func errHandler(c *gin.Context, errorMessage string) {
	if r := recover(); r != nil {
		c.JSON(400, gin.H{"message": errorMessage})
	}
}

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
	defer db.Close()

	if result.Error != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("user created, %v", result.Error)})
	} else {
		c.JSON(201, gin.H{"message": fmt.Sprintf("user created, %v", user.UserName)})
	}
}

func CheckPermission(c *gin.Context) {
	defer errHandler(c, "error to parse header")
	token := c.Request.Header["Token"][0]
	d, _ := AppCache.Get(token)
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(string(d)), &result)
	if result != nil {
		c.JSON(200, result)
		return
	}
	c.JSON(200, gin.H{"message": "user not found"})

}

func RoleCreate(c *gin.Context) {
	var param CreateRole
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(200, gin.H{"message": "error in connect to database"})
		return
	}
	role := &models.Role{
		Name:   param.Name,
		FaName: param.FaName,
	}
	result := db.Create(&role)

	if result.Error != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("error:, %v", result.Error)})
	} else {
		c.JSON(201, gin.H{"message": fmt.Sprintf("role created, %v", role.Name)})
	}
	defer db.Close()

}

func PermissionCreated(c *gin.Context) {
	var param CreatePermission
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(200, gin.H{"message": "error in connect to database"})
		return
	}
	var role models.Role
	dbResult := db.Where(&models.Role{Name: param.RoleName}).First(&role)
	if dbResult.Error != nil {
		c.JSON(400, gin.H{"message": "error find role"})
		return
	}

	permission := models.Permission{
		Name:   param.Name,
		RoleId: role.ID,
	}
	db.Create(&permission)

	c.JSON(201, gin.H{"message": "permission created"})
	defer db.Close()
}

func UserPermission(c *gin.Context) {
	var param CreateUserPermission
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(200, gin.H{"message": "error in connect to database"})
		return
	}
	var user models.User
	var permission models.Permission

	userResult := db.Where(&models.User{UserId: param.UserId}).First(&user)
	if userResult.Error != nil {
		c.JSON(400, gin.H{"message": "error find user"})
		return
	}

	permissionResult := db.Where(&models.Permission{Name: param.Permission}).First(&permission)
	if permissionResult.Error != nil {
		c.JSON(400, gin.H{"message": "error find permission"})
		return
	}

	db.Model(&user).Association("Permissions").Append(&permission)
	c.JSON(400, gin.H{"message": "permission append to user"})
	defer db.Close()
}
