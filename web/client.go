/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:23.
 * Author:  guojia(https://github.com/guojia99)
 */

package web

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"my-cubing/web/api"
	"my-cubing/web/web"
)

type Client struct {
	e *gin.Engine
}

func (c *Client) Run(port string) {
	c.e = gin.Default()
	//c.e.Use(gin.Logger(), gin.Recovery())

	route := c.e.Group("/")
	api.AddApiRoute(route)
	web.AddWebRoute(c.e)

	if err := c.e.Run(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		panic(err)
	}
}

func NewClient() *Client {
	return &Client{}
}
