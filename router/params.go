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
	UserName    string   `json:"user_name"`
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
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

type CreatePermission struct {
	Name     string `json:"name"`
	RoleName string `json:"role_name"`
}

type CreateUserPermission struct {
	UserId     string `json:"user_id"`
	Permission string `json:"permission"`
}
