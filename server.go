package main

import (
	"fmt"
	"log"
	"login_service/models"
)

func main() {
	db, err := models.Connect()
	if err != nil {
		log.Fatal("error in connect to db", err)
	}
	app := models.App{
		ClientId: "9101C3DC96E2FBC91FD4D85B11E344F441E8673891326FAABF4606AB1745828C",
		SecretKey: "CBAB3FA247EF0FE79DC14984F86398767BBFCACC502B727088A7F02040F400AB",
		WhiteListIp: "192.168.1.1,192.1681.2",
	}
	db.Create(&app)
	defer db.Close()
	//models.GetDB().
	fmt.Println("server run")
}