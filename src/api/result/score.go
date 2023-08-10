/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/19 下午5:59.
 *  * Author: guojia(https://github.com/guojia99)
 */

package result

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/core/model"
	"github.com/guojia99/my-cubing/src/svc"
)

type (
	CreateScoreRequest struct {
		PlayerName string             `json:"PlayerName"`
		ContestID  uint               `json:"ContestID"`
		Project    model.Project      `json:"Project"`
		RouteNum   int                `json:"RouteNum"`
		Penalty    model.ScorePenalty `json:"Penalty"`
		Results    []float64          `json:"Results"`
	}
)

func CreateScore(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateScoreRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Core.AddScore(req.PlayerName, req.ContestID, req.Project, req.RouteNum, req.Results, req.Penalty); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

type (
	DeleteScoreRequest struct {
		PlayerName string        `json:"PlayerName"`
		ContestID  uint          `json:"ContestID"`
		Project    model.Project `json:"Project"`
		RouteNum   int           `json:"RouteNum"`
	}
)

func DeleteScore(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req DeleteScoreRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Core.RemoveScore(req.PlayerName, req.ContestID, req.Project, req.RouteNum); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

type EndContestRequest struct {
	ContestID uint `json:"ContestID"`
}

func EndContest(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req EndContestRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Core.StatisticalRecordsAndEndContest(req.ContestID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
