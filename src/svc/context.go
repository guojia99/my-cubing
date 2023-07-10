/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午6:24.
 *  * Author: guojia(https://github.com/guojia99)
 */

package svc

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing/src/model"
	"github.com/guojia99/my-cubing/src/model/core"
)

type Context struct {
	DB   *gorm.DB
	Cfg  *Config
	Core core.Core
}

func NewContext(config string) (*Context, error) {
	ctx := &Context{
		Cfg: &Config{},
	}
	if err := ctx.Cfg.Load(config); err != nil {
		return nil, err
	}
	var err error
	switch ctx.Cfg.DB.Driver {
	case "sqlite":
		ctx.DB, err = gorm.Open(sqlite.Open(ctx.Cfg.DB.DSN), &gorm.Config{})
	case "mysql":
		ctx.DB, err = gorm.Open(mysql.New(mysql.Config{DSN: ctx.Cfg.DB.DSN}), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}
	if err = ctx.DB.AutoMigrate(model.Models...); err != nil {
		return nil, err
	}

	ctx.Core = core.NewScoreCore(ctx.DB)
	return ctx, nil
}
