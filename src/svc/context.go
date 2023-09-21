/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午6:24.
 *  * Author: guojia(https://github.com/guojia99)
 */

package svc

import (
	"time"

	cache "github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/guojia99/my-cubing/src/core"
	"github.com/guojia99/my-cubing/src/core/model"
)

type Context struct {
	DB    *gorm.DB
	Cfg   *Config
	Core  core.Core
	Cache *cache.Cache
}

func NewContext(config string) (*Context, error) {
	ctx := &Context{
		Cfg:   &Config{},
		Cache: cache.New(time.Minute, time.Minute),
	}
	if err := ctx.Cfg.Load(config); err != nil {
		return nil, err
	}
	var err error
	switch ctx.Cfg.DB.Driver {
	case "sqlite":
		ctx.DB, err = gorm.Open(sqlite.Open(ctx.Cfg.DB.DSN), &gorm.Config{})
	case "mysql":
		ctx.DB, err = gorm.Open(mysql.New(mysql.Config{DSN: ctx.Cfg.DB.DSN}), &gorm.Config{
			Logger: logger.Discard,
		})
	}
	if err != nil {
		return nil, err
	}
	if err = ctx.DB.AutoMigrate(model.Models...); err != nil {
		return nil, err
	}

	//if err = ctx.DB.AutoMigrate(model.Models...); err != nil {
	//	return nil, err
	//}

	ctx.Core = core.NewScoreCore(ctx.DB, ctx.Cfg.Debug)
	return ctx, nil
}
