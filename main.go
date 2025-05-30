package main

import (
	"log"

	"github.com/anuragdaksh7/url-short-go/config"
	"github.com/anuragdaksh7/url-short-go/internal/user"
	"github.com/anuragdaksh7/url-short-go/router"
)

var _config config.Config

func init() {
	var err error
	_config, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	config.ConnectDB()
	config.SyncDB()
}

func main() {
	userSvc := user.NewService()

	userHandler := user.NewHandler(userSvc)

	router.InitRouter(userHandler)
	log.Fatal(router.Start("0.0.0.0:" + _config.PORT))
}
