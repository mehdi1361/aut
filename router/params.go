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

type LoginParamResponse struct {
	UserName string `json:"user_name"`
	Token    string `json:"token"`
}

type CreateUserParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	MobileNo string `json:"mobile_no"`
}

type CreateRole struct {
	Name   string `json:"name"`
	FaName string `json:"fa_name"`
}
