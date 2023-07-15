/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午6:36.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

func (c *Client) initRoute() {

	api := c.e.Group("/api")
	{
		user := api.Group("/user")
		user.POST("")
	}

	{
		contest := api.Group("/contest")
		contest.GET("/")
		contest.POST("/")
		contest.DELETE("/")

		player := api.Group("/player")
		player.GET("/")
		player.POST("/")
		player.DELETE("/")
	}

	{
		score := api.Group("/score")
		score.GET("/player/:player_name/contest/:contest_id")    // 上传成绩
		score.DELETE("/player/:player_name/contest/:contest_id") // 删除成绩
		score.POST("/contest/:contest_id/end")                   // 结束比赛并统计
	}

	{
		report := api.Group("/report")
		report.GET("/project_best_score")         // 所有角色所有成绩的最佳列表， 按project分类
		report.GET("/best_score")                 // 所有成绩最佳， 仅有一个单次和平均
		report.GET("/sor")                        // sor成绩统计
		report.GET("/sor/contest/:contest_id")    // 某比赛的sor
		report.GET("/contest/:contest_id/score")  // 某比赛的成绩统计
		report.GET("/contest/:contest_id/podium") // 某场比赛领奖台
		report.GET("/player/:player_name")        // 某个玩家的成绩汇总
		report.GET("/player/:player_name/podium") // 某个玩家的领奖台
		report.GET("/podium")                     // 获取所有玩家领奖台的排行
	}

}
