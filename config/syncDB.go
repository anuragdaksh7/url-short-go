package config

import "github.com/anuragdaksh7/url-short-go/model"

func SyncDB() {
	DB.AutoMigrate(
		&model.User{},
		&model.Url{},
	)
}
