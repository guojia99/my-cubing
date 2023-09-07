/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午6:36.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

import (
	swagFile "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/guojia99/my-cubing/src/api/report"
	"github.com/guojia99/my-cubing/src/api/result"
)

func (c *Client) initRoute() {
	api := c.e.Group("/v2/api")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swagFile.Handler))

	{ // 后台

		// 授权
		api.POST("/auth/token", c.ValidToken) // 获取授权

		// 比赛
		api.POST("/contest", c.AuthMiddleware, result.CreateContest(c.svc))               // 添加比赛
		api.DELETE("/contest/:contest_id", c.AuthMiddleware, result.DeleteContest(c.svc)) // 删除比赛

		// 玩家
		api.POST("/player", c.AuthMiddleware, result.CreatePlayer(c.svc)) //  添加玩家或修改玩家
		api.DELETE("/player", c.AuthMiddleware, result.DeletePlayer(c.svc))

		// 成绩
		api.GET("/score/player/:player_id/contest/:contest_id", c.AuthMiddleware, result.GetScores(c.svc)) // 获取某场比赛玩家的所有成绩
		api.POST("/score", c.AuthMiddleware, result.CreateScore(c.svc))                                    // 上传成绩
		api.DELETE("/score/:score_id", c.AuthMiddleware, result.DeleteScore(c.svc))                        // 删除成绩
		api.PUT("/score/end_contest", c.AuthMiddleware, result.EndContest(c.svc))                          // 结束比赛并统计
	}

	{ //开发日志，关键点
		xLog := api.Group("x-log")
		xLog.GET("/")
		xLog.PUT("/")
	}

	{ // 基础查询
		api.GET("/record", result.GetRecords(c.svc))
		api.GET("/contest", result.GetContests(c.svc))
		api.GET("/contest/:contest_id", result.GetContest(c.svc))
		api.GET("/player", result.GetPlayers(c.svc))
		api.GET("/player/:player_id", result.GetPlayer(c.svc))
		api.GET("/projects", result.ProjectList(c.svc))
	}

	{ // 榜单
		rp := api.Group("/report")
		{
			// 排行榜
			rp.GET("/best/score", report.BestReport(c.svc))              // 获取最佳成绩榜单，每个项目仅有一个单次和平均
			rp.GET("/best/all_scores", report.BestAllScoreReport(c.svc)) // 获取项目每个玩家最佳成绩
			rp.GET("/best/sor", report.BestSorReport(c.svc))             // 获取所有角色的sor汇总榜单
			rp.GET("/best/podium", report.BestPodiumReport(c.svc))       // 获取所有玩家领奖台的排行

			// 具体到比赛
			rp.GET("/contest/:contest_id/sor", report.ContestSorReport(c.svc))       // 某比赛的sor
			rp.GET("/contest/:contest_id/score", report.ContestScoreReport(c.svc))   // 某比赛的成绩统计
			rp.GET("/contest/:contest_id/podium", report.ContestPodiumReport(c.svc)) // 某场比赛领奖台
			rp.GET("/contest/:contest_id/record", report.ContestRecord(c.svc))       // 某场比赛的记录

			// 具体到个人
			rp.GET("/player/:player_id/best", report.PlayerBest(c.svc))           // 某玩家的最佳成绩
			rp.GET("/player/:player_id/score", report.PlayerScoreReport(c.svc))   // 某个玩家的成绩汇总
			rp.GET("/player/:player_id/podium", report.PlayerPodiumReport(c.svc)) // 某个玩家的领奖台
			rp.GET("/player/:player_id/record", report.PlayerRecord(c.svc))       // 某个玩家的记录
		}
	}
}
