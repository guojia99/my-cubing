/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:30.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import "github.com/gin-gonic/gin"

func curd(route *gin.RouterGroup, c, rMui, u, d gin.HandlerFunc) {
	if c != nil {
		route.POST("/", c)
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
		curd(api.Group("/players"), CreatePlayer, ReadPlayers, UpdatePlayers, nil)              // 选手增删改查
		curd(api.Group("/contests"), CreateContest, ReadContests, UpdateContest, DeleteContest) // 比赛增删改查

		// 成绩
		score := api.Group("/score")
		{
			score.GET("/player/:player_name/contest/:contest_id", GetUserContestScore) // 获取某个选手某场成绩
			score.POST("/", CreateScore)                                               // 上传成绩
			score.DELETE("/", DeleteScore)                                             // 移除成绩
		}

		// 统计报表
		report := score.Group("/report")
		{
			report.GET("/player/:user_name", ReadReportByUser)       // 个人成绩统计
			report.GET("/contests/:contest_id", ReadReportByContest) // 某比赛成绩统计
			report.GET("/best", ReadReportBest)                      // 所有项目的最佳成绩总排名积分
			report.GET("/best/:project", ReadReportByProjectBest)    // 某项目的最佳成绩统计
		}
	}
}
