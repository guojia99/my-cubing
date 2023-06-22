/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:30.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import "github.com/gin-gonic/gin"

func curd(route *gin.RouterGroup, c, r, rMui, u, d gin.HandlerFunc) {
	if c != nil {
		route.POST("/", c)
	}

	if r != nil {
		route.GET("/:id", r)
	}

	if rMui != nil {
		route.GET("/", rMui)
	}

	if u != nil {
		route.PUT("/:id", u)
	}

	if d != nil {
		route.DELETE("/", d)
	}
}

func AddApiRoute(route *gin.RouterGroup) {
	api := route.Group("/api")
	{
		curd(api.Group("/players"), CreatePlayer, nil, ReadPlayers, UpdatePlayers, nil)             // 选手增删改查
		curd(api.Group("/contest"), CreateContest, nil, ReadContests, UpdateContest, DeleteContest) // 比赛增删改查

		// 成绩
		score := api.Group("/score")
		{
			score.POST("/user/:user_id/contest/:contest_id")   // 上传成绩
			score.DELETE("/user/:user_id/contest/:contest_id") // 移除成绩
		}

		// 统计报表
		report := score.GET("/report")
		{
			report.GET("/user/:user_id")        // 个人成绩统计
			report.GET("/contests/:contest_id") // 某比赛成绩统计
			report.GET("/best")                 // 所有项目的最佳成绩总排名积分
			report.GET("/best/:project")        // 某项目的最佳成绩统计
		}
	}
}
