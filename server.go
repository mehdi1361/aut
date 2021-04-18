package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"login_service/router"
)

func main() {
	r := gin.Default()
	r.GET("/ping", router.PingHandler)
	r.GET("/app_access", router.AppAccessHandler)
	r.GET("/login", router.UserLoginHandler)
	r.POST("/user", router.CreateUser)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
