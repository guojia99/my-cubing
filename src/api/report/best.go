/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/17 下午5:01.
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

type BestReportResponse struct {
	BestSingle map[model.Project]model.Score `json:"BestSingle"`
	BestAvg    map[model.Project]model.Score `json:"BestAvg"`
}

func BestReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bestSingle, bestAvg := svc.Core.GetBestScores()
		ctx.JSON(http.StatusOK, BestReportResponse{
			BestSingle: bestSingle,
			BestAvg:    bestAvg,
		})
	}
}

type (
	BestAllScoreReportRequest struct {
		Project string `query:"project"`
	}

	BestAllScoreReportResponse struct {
		BestSingle map[model.Project][]model.Score `json:"BestSingle"`
		BestAvg    map[model.Project][]model.Score `json:"BestAvg"`
	}
	BestAllScoreReportByProjectResponse struct {
		BestSingle []model.Score `json:"BestSingle"`
		BestAvg    []model.Score `json:"BestAvg"`
	}
)

func BestAllScoreReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		project := ctx.Query("project")
		if len(project) == 0 {
			bestSingle, bestAvg := svc.Core.GetAllPlayerBestScore()
			ctx.JSON(http.StatusOK, BestAllScoreReportResponse{
				BestSingle: bestSingle,
				BestAvg:    bestAvg,
			})
			return
		}

		bestSingle, bestAvg := svc.Core.GetAllPlayerBestScoreByProject(model.Project(project))
		ctx.JSON(http.StatusOK, BestAllScoreReportByProjectResponse{
			BestSingle: bestSingle,
			BestAvg:    bestAvg,
		})
	}
}

type BestSorReportResponse struct {
	BestSingle []core.SorScore `json:"BestSingle"`
	BestAvg    []core.SorScore `json:"BestAvg"`
}

func BestSorReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bestSingle, bestAvg := svc.Core.GetSorScore()
		ctx.JSON(http.StatusOK, BestSorReportResponse{
			BestSingle: bestSingle,
			BestAvg:    bestAvg,
		})
	}
}

func BestPodiumReport(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, svc.Core.GetAllPodium())
	}
}
