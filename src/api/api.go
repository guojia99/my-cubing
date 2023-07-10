/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午5:59.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/svc"
)

func NewClient(s *svc.Context) *Client { return &Client{svc: s} }

type Client struct {
	e   *gin.Engine
	svc *svc.Context
}

func (c *Client) Run() error {
	gin.SetMode(c.svc.Cfg.GinMode)
	c.e = gin.New()
	c.e.Use(gin.Logger(), gin.Recovery())

	return c.e.Run(fmt.Sprintf("0.0.0.0:%d", c.svc.Cfg.Port))
}
