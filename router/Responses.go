package router

import (
	"encoding/json"
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
