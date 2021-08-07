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
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginParamResponse struct {
	UserName    string   `json:"username"`
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
	UserType    string   `json:"user_type"`
	IsSuperuser bool     `json:"is_superuser"`
}

type CreateUserParam struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	MobileNo    string `json:"mobile_no"`
	UserType    string `json:"user_type"`
	IsSuperuser bool   `json:"is_superuser"`
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
	UserId     int    `json:"user_id"`
	Permission string `json:"permission"`
}

type CreateCustomerRole struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateCustomerRoleUser struct {
	UserId       string `json:"user_id"`
	CustomerRole string `json:"customer_role"`
}

type ChangeActiveState struct {
	UserId int  `json:"user_id"`
	State  bool `json:"state"`
}

type UserDataUpdate struct {
	UserId      int    `json:"user_id"`
	MobileNo    string `json:"mobile_no"`
	IsSuperUser bool   `json:"is_superuser"`
}

type GetUserPermissionParam struct {
	UserId int `json:"user_id"`
}

type GetUserPermissionResponse struct {
	Permissions []string `json:"permissions"`
}
