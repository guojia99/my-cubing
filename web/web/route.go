/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:30.
 * Author:  guojia(https://github.com/guojia99)
 */

package web

import "github.com/gin-gonic/gin"

func AddWebRoute(route *gin.RouterGroup) {
	web := route.Group("/web")
	{
		web.GET("/score")  // 成绩记录页
		web.GET("/report") // 报告页
	}
}
