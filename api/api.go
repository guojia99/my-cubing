/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/20 下午10:47.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import "github.com/gin-gonic/gin"

type Client struct {
	e *gin.Engine
}

func (c *Client) Run() {
	c.e = gin.Default()
	c.e.Use(gin.Logger(), gin.Recovery())
	c.AddRoute()
}
