/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:35.
 * Author:  guojia(https://github.com/guojia99)
 */

package db

import (
	"os"

	json "github.com/json-iterator/go"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type Config struct {
	Db struct {
		Driver string `json:"driver"`
		Path   string `json:"path"`
	} `json:"db"`
}

func Init() {
	configBody, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err = json.Unmarshal(configBody, &cfg); err != nil {
		panic(err)
	}

	switch cfg.Db.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.Db.Path), &gorm.Config{})
	case "mysql":
		DB, err = gorm.Open(mysql.New(mysql.Config{DSN: cfg.Db.Path}), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(
		&Contest{},
		&Score{},
		&Player{},
	)
	if err != nil {
		panic(err)
	}
}
