package router

type ParamsApp struct {
	ClientId  string `json:"client_id"`
	SecretKey string `json:"secret_key"`
}

type ResponseParamApp struct {
	Message interface{} `json:"message"`
	Access  bool        `json:"access"`
}

type LoginParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type CreateUserParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	MobileNo string `json:"mobile_no"`
}
