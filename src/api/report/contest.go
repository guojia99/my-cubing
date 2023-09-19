/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/19 下午2:17.
 *  * Author: guojia(https://github.com/guojia99)
 */

package report

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/core"
	"github.com/guojia99/my-cubing/src/core/model"
	"github.com/guojia99/my-cubing/src/svc"
)

type ContestRequest struct {
	ContestID uint `uri:"contest_id"`
}

type (
	ContestSorReportResponse struct {
		Single map[model.SorStatisticsKey][]core.SorScore `json:"Single"`
		Avg    map[model.SorStatisticsKey][]core.SorScore `json:"Avg"`
	}
)

func ContestSorReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ContestRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		single, avg := svc.Core.GetSorScoreByContest(req.ContestID)
		ctx.JSON(http.StatusOK, ContestSorReportResponse{
			Single: single,
			Avg:    avg,
		})
	}
}

type ContestScoreReportResponse struct {
	Scores map[model.Project][]core.RoutesScores `json:"Scores"`
}

func ContestScoreReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ContestRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, ContestScoreReportResponse{
			Scores: svc.Core.GetScoreByContest(req.ContestID),
		})
	}
}

func ContestPodiumReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ContestRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, svc.Core.GetPodiumsByContest(req.ContestID))
	}
}

func ContestRecord(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ContestRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, svc.Core.GetRecordByContest(req.ContestID))
	}
}
