package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	UserName       string       `json:"user_name" gorm:"size:60;unique"`
	Password       string       `json:"password" gorm:"size:250"`
	MobileNo       string       `json:"mobile_no" gorm:"size:11;unique"`
	UserId         string       `json:"user_id" gorm:"unique"`
	Active         bool         `json:"active" gorm:"default:true"`
	ChangePassword bool         `json:"change_password" gorm:"default:true"`
	Permissions    []Permission `json:"permissions" gorm:"many2many:auth_user_service_permission"`
}

func (u *User) TableName() string {
	return "auth_service_user"
}

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
	Name   string `json:"name" gorm:"size:50"`
	RoleId uint   `json:"role_id" gorm:"Column:role_id"`
}

func (p Permission) TableName() string {
	return "auth_service_permissions"
}

func Connect() (db *gorm.DB, err error) {
	conn, err := gorm.Open(
		"postgres",
		"host=localhost port=5432 user=mebco_user dbname=mebco_db password=mebco_1060 sslmode=disable",
	)
	return conn, err
}

func init() {
	conn, err := Connect()
	if err != nil {
		fmt.Print(err)
	}
	defer conn.Close()

	db := conn
	_ = db.AutoMigrate(&User{}, &Role{}, &Permission{})
	db.Model(&Permission{}).AddForeignKey("role_id", "auth_service_role(id)", "CASCADE", "CASCADE")
}
