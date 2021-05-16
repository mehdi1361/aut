package router

import (
	"encoding/json"
	"fmt"
	"login_service/common"
	"login_service/models"
	"time"
)

func convertParamToDict(v interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(b), &result); err != nil {
		return nil, err
	} else {
		return result, nil
	}

}

func AppCheckResponse(f ParamsApp) map[string]interface{} {
	db, err := models.Connect()
	if err != nil {
		response, _ := convertParamToDict(ResponseParamApp{err, false})
		return response
	}

	var app models.App
	db.Where(&models.App{ClientId: f.ClientId}).First(&app)

	defer db.Close()
	if f.SecretKey == "" {
		response, _ := convertParamToDict(&ResponseParamApp{"secret key not found", false})
		return response
	}
	if app.SecretKey == f.SecretKey {
		response, _ := convertParamToDict(&ResponseParamApp{"access granted", true})
		return response
	} else {
		response, _ := convertParamToDict(&ResponseParamApp{"app authorization failed", false})
		return response
	}
}

func UserLogin(param LoginParam) (int, map[string]interface{}, error) {
	db, err := models.Connect()
	if err != nil {
		return 400, nil, err
	}

	var user models.User
	var permissions []models.Permission

	t := common.GetMD5Hash(param.Password)

	//_ = db.Preload("Permissions").Where(&models.User{UserName: param.UserName, Password: t}).First(&user)
	db.First(&user, "user_name=? and password=?", param.UserName, t)
	if t != user.Password {
		result := make(map[string]interface{})
		result["message"] = "user not found"
		defer db.Close()

		return 404, result, nil
	}
	db.Model(&user).Related(&permissions, "Permissions")

	var lstNamePermission []string
	for _, v := range permissions {
		lstNamePermission = append(lstNamePermission, v.Name)
	}
	defer db.Close()

	token := fmt.Sprintf("%s-%s", common.TokenGenerator(user), user.UserId)
	s, err := convertParamToDict(
		&LoginParamResponse{
			UserName:    param.UserName,
			Token:       token,
			Permissions: lstNamePermission,
			UserType:    user.UserType,
		},
	)
	if err != nil {
		return 400, nil, err
	}
	AppCache.KeyContainDelete(user.UserId)
	_ = AppCache.Set(token, s, 20*time.Hour)
	return 200, s, nil
}
