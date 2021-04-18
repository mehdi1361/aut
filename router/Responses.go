package router

import (
	"encoding/json"
	"login_service/common"
	"login_service/models"
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
	dbResult := db.Where(&models.User{UserName: param.UserName, Password: common.GetMD5Hash(param.Password)}).First(&user)
	if dbResult.Error != nil {
		return 400, nil, err
	}

	s, err := convertParamToDict(
		&LoginParamResponse{
			UserName: param.UserName,
			Token:    common.TokenGenerator(user.UserName, user.Password),
		},
	)
	if err != nil {
		return 400, nil, err
	}
	return 200, s, nil
}
