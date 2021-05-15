package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/patrickmn/go-cache"
	"log"
	"login_service/router"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/ping/", router.PingHandler)
	r.GET("/app_access/", router.AppAccessHandler)
	r.POST("/login/", router.UserLoginHandler)
	r.POST("/user/", router.CreateUser)
	r.GET("/check/", router.CheckPermission)
	r.POST("/role/", router.RoleCreate)
	r.POST("/permission/", router.PermissionCreated)
	r.POST("/user/permission/", router.UserPermission)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	port := os.Getenv("PORT")
	err = r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
