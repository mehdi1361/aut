package router

import (
	"aut/common"
	"aut/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"strconv"
	"strings"
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
		c.JSON(500, gin.H{"Message": "error in connect to database"})
		return
	}
	user := models.User{
		UserName:    param.UserName,
		Password:    common.GetMD5Hash(param.Password),
		UserId:      common.UuidGenerator(),
		MobileNo:    param.MobileNo,
		UserType:    param.UserType,
		IsSuperUser: param.IsSuperuser,
	}
	result := db.Create(&user)
	defer db.Close()

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("user created, %v", result.Error),
		})
	} else {
		c.JSON(201, gin.H{
			"message": fmt.Sprintf("user created, %v", user.UserName),
			"id":      user.UserId,
		})
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

func ListRole(c *gin.Context) {
	var records []models.Role
	db, err := models.Connect()

	if err != nil {
		c.JSON(200, gin.H{"message": fmt.Sprintf("error in connect to database %s", err)})
		return
	}
	db.Find(&records)
	c.JSON(200, records)
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

func ListPermission(c *gin.Context) {
	var records []models.Permission
	db, err := models.Connect()

	if err != nil {
		c.JSON(200, gin.H{"message": fmt.Sprintf("error in connect to database %s", err)})
		return
	}
	db.Find(&records)
	c.JSON(200, records)
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

	userResult := db.Where("id = ?", param.UserId).First(&user)
	if userResult.Error != nil {
		c.JSON(400, gin.H{"message": "error find user"})
		return
	}

	lstPermission := strings.Split(param.Permission, ",")
	for _, v := range lstPermission {
		var permission models.Permission
		data, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(400, gin.H{"message": "error permission id is not valid"})
			return
		}
		permissionResult := db.Where("id = ?", data).First(&permission)
		if permissionResult.Error != nil {
			c.JSON(400, gin.H{"message": "error find permission"})
			return
		}
		db.Model(&user).Association("Permissions").Append(&permission)
	}

	c.JSON(200, gin.H{"message": "permissions append to user"})
	defer db.Close()
}

func EditUser(c *gin.Context) {
	var param CreateUserParam
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(500, gin.H{"Message": "error in connect to database"})
		return
	}
	user := models.User{
		UserName:    param.UserName,
		Password:    common.GetMD5Hash(param.Password),
		UserId:      common.UuidGenerator(),
		MobileNo:    param.MobileNo,
		UserType:    param.UserType,
		IsSuperUser: param.IsSuperuser,
	}
	result := db.Create(&user)
	defer db.Close()

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("user created, %v", result.Error),
		})
	} else {
		c.JSON(201, gin.H{
			"message": fmt.Sprintf("user created, %v", user.UserName),
			"id":      user.UserId,
		})
	}
}

func ListCustomerRole(c *gin.Context) {
	var records []models.CustomerRole
	db, err := models.Connect()

	if err != nil {
		c.JSON(200, gin.H{"message": fmt.Sprintf("error in connect to database %s", err)})
		return
	}
	db.Find(&records)
	c.JSON(200, records)
}

func CustomerRoleCreated(c *gin.Context) {
	var param CreateCustomerRole
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

	permission := models.CustomerRole{
		Name: param.Name,
		Type: param.Type,
	}
	db.Create(&permission)

	c.JSON(201, gin.H{"message": "permission created"})
	defer db.Close()
}

func CustomerUserRole(c *gin.Context) {
	var param CreateCustomerRoleUser
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
	var customerRole models.CustomerRole

	userResult := db.Where(&models.User{UserId: param.UserId}).First(&user)
	if userResult.Error != nil {
		c.JSON(400, gin.H{"message": "error find user"})
		return
	}

	permissionResult := db.Where(&models.CustomerRole{Name: param.CustomerRole}).First(&customerRole)
	if permissionResult.Error != nil {
		c.JSON(400, gin.H{"message": "error find permission"})
		return
	}

	db.Model(&user).Association("CustomerRole").Append(&customerRole)
	c.JSON(200, gin.H{"message": "permission append to user"})
	defer db.Close()
}

func UpdateActiveState(c *gin.Context) {
	var param ChangeActiveState
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(400, gin.H{"message": "error in connect to database"})
		return
	}
	db.Model(&models.User{}).Where("id = ?", param.UserId).Update("active", param.State)
	c.JSON(200, gin.H{"message": fmt.Sprintf("change active to %t", param.State)})
	defer db.Close()
}
func UpdateUserData(c *gin.Context) {
	var param UserDataUpdate
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		log.Printf("%+v", err)
		c.JSON(400, err)
		return
	}
	db, err := models.Connect()
	if err != nil {
		c.JSON(400, gin.H{"message": "error in connect to database"})
		return
	}

	if len(param.MobileNo) > 11 {
		c.JSON(400, gin.H{"message": fmt.Sprintf("len %s is greater than 11", param.MobileNo)})
		return
	}
	db.Model(&models.User{}).Where("id = ?", param.UserId).Update("is_super_user", param.IsSuperUser)
	db.Model(&models.User{}).Where("id = ?", param.UserId).Update("mobile_no", param.MobileNo)
	c.JSON(200, gin.H{"message": fmt.Sprintf("user with id:%d updated", param.UserId)})
	defer db.Close()
}
