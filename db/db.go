/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:35.
 * Author:  guojia(https://github.com/guojia99)
 */

package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	dbDriver = "mysql"
)

func init() {
	var err error

	switch dbDriver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	case "mysql":
		cfg := mysql.Config{
			DSN: "root:my123456@tcp(127.0.0.1:3306)/mycube?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		}
		DB, err = gorm.Open(mysql.New(cfg), &gorm.Config{})
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
