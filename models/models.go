package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

type User struct {
	gorm.Model
	UserName       string `json:"user_name" gorm:"size:60;unique"`
	Password       string `json:"password" gorm:"size:250"`
	MobileNo       string `json:"mobile_no" gorm:"size:11;unique"`
	Active         bool   `json:"active" gorm:"default:true"`
	ChangePassword bool   `json:"change_password" gorm:"default:true"`
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
	Name        string       `json:"name" gorm:"size:50;unique"`
	FaName      string       `json:"fa_name"`
	Permissions []Permission `json:"permissions" gorm:"foreignKey:RoleName"`
}

func (r *Role) TableName() string {
	return "auth_service_role"
}

type Permission struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:50"`
	RoleName string `json:"role_name"`
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
	_ = db.AutoMigrate(&User{}, &App{}, &Role{}, &Permission{})
}

func GetDB() *gorm.DB {
	return db
}
