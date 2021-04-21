package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/patrickmn/go-cache"
	"log"
	"login_service/router"
)

func main() {

	r := gin.Default()
	r.GET("/ping", router.PingHandler)
	r.GET("/app_access", router.AppAccessHandler)
	r.POST("/login", router.UserLoginHandler)
	r.POST("/user", router.CreateUser)
	r.GET("/check", router.CheckPermission)
	r.POST("/role", router.RoleCreate)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
