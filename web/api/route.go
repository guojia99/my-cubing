/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:30.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import (
	"github.com/gin-gonic/gin"

	"my-cubing/web/api/score"
)

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
		scoreRoute := api.Group("/score")
		{
			scoreRoute.GET("/player/:player_name/contest/:contest_id", score.GetUserContestScore)                 // 获取某个选手某场成绩
			scoreRoute.POST("/", score.CreateScore)                                                               // 上传成绩
			scoreRoute.DELETE("/player/:player_name/contest/:contest_id/project/:project_key", score.DeleteScore) // 移除成绩

			scoreRoute.GET("/report/all_project_score", score.GetProjectScores)       // 获取所有项目成绩列表
			scoreRoute.GET("/report/all_project_best", score.GetAllProjectBestScore)  // 获取所有项目最佳的成绩, 有且仅有一个最佳和一个单次
			scoreRoute.GET("/report/all_sor", score.GetSorScores)                     // sor 排名
			scoreRoute.GET("/report/contest/:contest_id", score.GetContestScores)     // 某场比赛的成绩汇总
			scoreRoute.POST("/report/contest/:contest_id/end", score.EndContestScore) // 结束某场比赛并开始统计
		}
	}
}
