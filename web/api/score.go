/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:41.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"my-cubing/db"
)

type CreateScoreRequest struct {
	PlayerName  string    `json:"PlayerName"`
	ContestID   uint      `json:"ContestID"`
	RouteNumber uint      `json:"RouteNumber"`
	ProjectName string    `json:"ProjectName"`
	Results     []float64 `json:"Results"`
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

	var p = db.Player{Name: req.PlayerName}
	if err := db.DB.FirstOrCreate(&p, "name = ?", req.PlayerName).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var score db.Score
	if err := db.DB.Model(&db.Score{}).
		Where("player_id = ?", p.ID).
		Where("contest_id = ?", req.ContestID).
		Where("route_number = ?", 1).
		Where("project = ?", pj).First(&score).Error; err == nil && score.ID != 0 {
		// todo update

		score.Result1 = req.Results[0]
		score.Result2 = req.Results[0]
		score.Result3 = req.Results[0]
		score.Result4 = req.Results[0]
		score.Result5 = req.Results[0]

	}

}

func DeleteScore(ctx *gin.Context) {}
