package models

import (
	"aut/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/qor/validations"
	"os"
)

type User struct {
	gorm.Model
	UserName       string          `json:"user_name" gorm:"size:60;unique;index:idx_name"`
	Password       string          `json:"password" gorm:"size:250"`
	MobileNo       string          `json:"mobile_no" gorm:"size:11"`
	UserId         string          `json:"user_id" gorm:"unique"`
	Active         bool            `json:"active" gorm:"default:true"`
	ChangePassword bool            `json:"change_password" gorm:"default:true"`
	UserType       string          `json:"user_type" gorm:"size:60;index:idx_name"`
	IsSuperUser    bool            `json:"is_superuser" gorm:"default:false"`
	Permissions    []*Permission   `json:"permissions" gorm:"many2many:auth_user_service_permission"`
	CustomerRole   []*CustomerRole `json:"customer_role" gorm:"many2many:auth_user_service_customer_role"`
	Groups         string          `json:"groups" gorm:"size:60"`
	Branches       string          `json:"branches" gorm:"size:60"`
}

func (u *User) TableName() string {
	return "auth_service_user"
}

//func (u *User) PermissionUpdate(param router.CreateUserPermission) error {
//	db, err := Connect()
//	defer db.Close()
//	if err != nil {
//		return err
//	}
//	var lstUserPermId []int
//	for _, v := range u.Permissions {
//		lstUserPermId = append(lstUserPermId, int(v.ID))
//	}
//	var newValidPerm []int
//	lstPermission := strings.Split(param.Permission, ",")
//
//	for _, v := range lstPermission {
//		var permission Permission
//		data, err := strconv.Atoi(v)
//		if err != nil {
//			return err
//		}
//		newValidPerm = append(newValidPerm, data)
//		permissionResult := db.Where("id = ?", data).First(&permission)
//		if permissionResult.Error != nil {
//			return err
//		}
//		db.Model(&u).Association("Permissions").Append(&permission)
//	}
//
//	for _, v := range common.Difference(lstUserPermId, newValidPerm) {
//		var permission Permission
//		permissionResult := db.Where("id = ?", v).First(&permission)
//		if permissionResult.Error != nil {
//			return err
//		}
//		db.Model(&u).Association("Permissions").Delete(&permission)
//	}
//	db.Model(&u).Association("Permissions").Delete()
//	fmt.Println(common.Difference(lstUserPermId, newValidPerm))
//	return nil
//}

type App struct {
	gorm.Model
	ClientId    string `json:"client_id" gorm:"size:256;unique"`
	SecretKey   string `json:"secret_key" gorm:"size:256"`
	WhiteListIp string `json:"-"`
}

func (a *App) TableName() string {
	return "auth_service_app"
}

type Role struct {
	gorm.Model
	Name        string       `json:"name" gorm:"size:50;unique;not null"`
	FaName      string       `json:"fa_name" orm:"size:50;unique;not null"`
	Permissions []Permission `json:"permissions"`
}

func (r *Role) TableName() string {
	return "auth_service_role"
}

type Permission struct {
	gorm.Model
	Name   string  `json:"name" gorm:"size:50;unique"`
	RoleId uint    `json:"role_id" gorm:"Column:role_id"`
	Users  []*User `json:"users" gorm:"many2many:auth_user_service_permission"`
}

func (p Permission) TableName() string {
	return "auth_service_permissions"
}

type CustomerRole struct {
	gorm.Model
	Name  string  `json:"name" gorm:"size:50;unique"`
	Type  string  `json:"type" gorm:"size:50;unique"`
	Users []*User `json:"users" gorm:"many2many:auth_user_service_customer_role"`
}

func (cr CustomerRole) TableName() string {
	return "auth_service_customer_roles"
}

func (cr CustomerRole) Validate(db *gorm.DB) {
	customSl := &utils.StringSlice{}
	customSl.DataReader([]interface{}{"branch", "comex"})
	filter := &utils.Filter{}

	if filter.Contain(customSl, cr.Type) {
		_ = db.AddError(validations.NewError(cr, "Type", "Type need to be in "))
	}
}

type Group struct {
	Id    string  `json:"id" gorm:"column:customer_group_id"`
	Name  string  `json:"name" gorm:"size:50;unique;column:customer_group_name"`
	Users []*User `json:"users" gorm:"many2many:auth_user_service_group"`
}

func (g Group) TableName() string {
	return "customer_group"
}

type Branch struct {
	Id    string  `json:"id" gorm:"column:branch_id"`
	Name  string  `json:"name" gorm:"column:branch_name"`
	Users []*User `json:"users" gorm:"many2many:auth_user_service_branch"`
}

func (b Branch) TableName() string {
	return "broker_branch"
}

func Connect() (db *gorm.DB, err error) {
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}
	server := os.Getenv("DB_SERVER")
	database := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	conn, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", server, port, dbUser, database, password),
	)
	validations.RegisterCallbacks(conn)
	return conn, err
}

func init() {
	conn, err := Connect()
	if err != nil {
		fmt.Print(err)
	}
	defer conn.Close()

	db := conn
	_ = db.AutoMigrate(&App{}, &User{}, &Role{}, &Permission{}, CustomerRole{}, &Group{}, &Branch{})

	db.Model(&Permission{}).AddForeignKey("role_id", "auth_service_role(id)", "CASCADE", "CASCADE")
	db.Model(&User{}).AutoMigrate()
}
