/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/20 下午1:57.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

import (
	"github.com/gin-gonic/gin"
)

func (c *Client) AddRoute() {
	route := c.e.Group("/")

	api := route.Group("/api")
	{
		curd(api.Group("/players"), nil, nil, nil, nil, nil)  // 选手增删改查
		curd(api.Group("/contests"), nil, nil, nil, nil, nil) // 比赛增删改查

		// 成绩
		score := api.Group("/score")
		{
			score.POST("/user/:user_id/contest/:contest_id")   // 上传成绩
			score.DELETE("/user/:user_id/contest/:contest_id") // 移除成绩

			report := score.GET("/report")
			report.GET("/user/:user_id")        // 个人成绩统计
			report.GET("/contests/:contest_id") // 某比赛成绩统计
			report.GET("/best")
		}
	}

	web := route.Group("/web")
	{
		web.GET("/")
	}
}

func curd(route *gin.RouterGroup, c, r, rMui, u, d gin.HandlerFunc) {
	route.GET("/:id", r)
	route.GET("/", rMui)
	route.POST("/", c)
	route.PUT("/:id", u)
	route.DELETE("/", d)
}
