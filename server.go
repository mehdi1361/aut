package main

import (
	"aut/router"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/patrickmn/go-cache"
	"log"
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
	r.PUT("/user/", router.EditUser)
	r.GET("/check/", router.CheckPermission)
	r.POST("/role/", router.RoleCreate)
	r.GET("/role/", router.ListRole)
	r.POST("/permission/", router.PermissionCreated)
	r.GET("/permission/", router.ListPermission)
	r.GET("/customer_role/", router.ListCustomerRole)
	r.POST("/customer_role/", router.CustomerRoleCreated)
	r.POST("/user/permission/", router.UserPermission)
	r.GET("/user/permission/", router.GetUserPermission)
	r.POST("/user/customer_role/", router.CustomerUserRole)
	r.PATCH("/user/change_state/", router.UpdateActiveState)
	r.PATCH("/user/edit/", router.UpdateUserData)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	port := os.Getenv("PORT")
	err = r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
