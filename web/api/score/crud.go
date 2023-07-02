/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/1 下午3:40.
 *  * Author: guojia(https://github.com/guojia99)
 */

package score

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"my-cubing/db"
)

func GetUserContestScore(ctx *gin.Context) {
	var req GetUserContestScoreRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var p db.Player
	if err := db.DB.First(&p, "name = ?", req.PlayerName).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var scores []db.Score
	db.DB.Model(&db.Score{}).Where("player_id = ?", p.ID).Where("contest_id = ?", req.ContestID).Find(&scores)
	ctx.JSON(http.StatusOK, gin.H{"data": scores})
	fmt.Println(scores)
}

func CreateScore(ctx *gin.Context) {
	var req CreateScoreRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	pj := db.StrToProject(req.ProjectName)
	if pj == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "project not found"})
		return
	}

	// 找到玩家或者生成一个新玩家
	var p = db.Player{Name: req.PlayerName}
	if err := db.DB.FirstOrCreate(&p, "name = ?", req.PlayerName).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var score db.Score
	// 找到上一个本比赛的成绩
	if err := db.DB.Model(&db.Score{}).
		Where("player_id = ?", p.ID).
		Where("contest_id = ?", req.ContestID).
		Where("route_number = ?", 1).
		Where("project = ?", pj).First(&score).Error; err == nil && score.ID != 0 {

		if err = score.SetResult(req.Results); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		fmt.Println(score)
		if err = db.DB.Save(&score).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
		return
	}

	score = db.Score{
		CreatedAt:   time.Now(),
		PlayerID:    p.ID,
		ContestID:   req.ContestID,
		RouteNumber: req.RouteNumber,
		Project:     db.StrToProject(req.ProjectName),
	}
	if err := score.SetResult(req.Results); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := db.DB.Save(&score).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func DeleteScore(ctx *gin.Context) {}
