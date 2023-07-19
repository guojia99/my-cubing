/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/17 下午6:55.
 *  * Author: guojia(https://github.com/guojia99)
 */

package result

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/core/model"
	"github.com/guojia99/my-cubing/src/svc"
)

type (
	Contest struct {
		Contest model.Contest `json:"Contest"`
		Rounds  []model.Round `json:"Rounds"`
	}

	GetContestsResponse struct {
		Size     int64     `json:"Size"`
		Count    int64     `json:"Count"`
		Contests []Contest `json:"Contests"`
	}
)

func GetContests(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Try to get the cache.
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(ctx.DefaultQuery("size", "50"))
		if size > 50 {
			size = 50
		}

		offset := (page - 1) * size
		limit := size

		key := fmt.Sprintf("GetContests_%d_%d", page, size)
		if val, ok := svc.Cache.Get(key); ok {
			ctx.JSON(http.StatusOK, val)
			return
		}

		// Find Contests.
		var contests []model.Contest
		if err := svc.DB.Offset(offset).Limit(limit).Find(&contests).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var count int64
		if err := svc.DB.Model(&model.Contest{}).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert to interface contents.
		var resp = GetContestsResponse{
			Count:    count,
			Contests: []Contest{},
		}
		for _, contest := range contests {
			var round []model.Round
			svc.DB.Find(&round, "id in ?", contest.GetRoundIds())
			for _, r := range round {
				r.GetUpsets()
			}
			resp.Contests = append(resp.Contests, Contest{
				Contest: contest,
				Rounds:  round,
			})
		}
		_ = svc.Cache.Add(key, resp, time.Second*30)
		ctx.JSON(http.StatusOK, resp)
	}
}

type (
	CreateContestRequestRound struct {
		Project model.Project `json:"Project"`
		Number  int           `json:"Number"`
		Name    string        `json:"Name"`
		Final   bool          `json:"Final"`
		Upsets  []string      `json:"Upsets"`
	}

	CreateContestRequest struct {
		Name        string                      `json:"Name"`
		Description string                      `json:"Description"`
		Rounds      []CreateContestRequestRound `json:"Rounds"`
		StartTime   int64                       `json:"StartTime"`
		EndTime     int64                       `json:"EndTime"`
	}
)

var defaultProjectRounds = func() []CreateContestRequestRound {
	var out []CreateContestRequestRound
	for _, p := range model.WCAProjectRoute() {
		out = append(out, CreateContestRequestRound{
			Project: p,
			Number:  1,
			Name:    fmt.Sprintf("%s单轮赛", p.Cn()),
			Final:   true,
			Upsets:  []string{},
		})
	}
	return out
}()

// CreateContest
// @Summary 创建比赛
// @Produce  json
// @Param Name body string true "Name"
// @Param Description body string true "Description"
// @Param Rounds body []CreateContestRequestRound true "Rounds"
// @Param StartTime body int true "StartTime"
// @Param EndTime body int true "EndTime"
// @Success 200 {object} CreateContestRequest
// @Router /api/contests [post]
func CreateContest(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateContestRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var contest model.Contest
		if err := svc.DB.Where("name = ?", req.Name).First(&contest); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"error": "contest is exist"})
			return
		}
		contest = model.Contest{
			Name:        req.Name,
			Description: req.Description,
			StartTime:   time.Unix(req.StartTime, 0),
			EndTime:     time.Unix(req.EndTime, 0),
		}

		if err := svc.DB.Save(&contest).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		if len(req.Rounds) == 0 || req.Rounds == nil {
			req.Rounds = defaultProjectRounds
		}

		var roundIds []uint
		for _, val := range req.Rounds {
			var round = model.Round{
				ContestID: contest.ID,
				Project:   val.Project,
				Number:    val.Number,
				Name:      val.Name,
				Final:     val.Final,
			}
			round.SetUpsets(val.Upsets)
			svc.DB.Create(&round)
			roundIds = append(roundIds, round.ID)
		}

		fmt.Println(roundIds)
		contest.SetRoundIds(roundIds)
		svc.DB.Save(&contest)

		ctx.JSON(http.StatusOK, gin.H{})
	}
}

type (
	DeleteContestRequest struct {
		Id uint `uri:"contest_id"`
	}
)

func DeleteContest(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req DeleteContestRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var contest model.Contest
		if err := svc.DB.Where("id = ?", req.Id).First(&contest).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var count int64
		if err := svc.DB.Model(&model.Score{}).Where("contest_id = ?", req.Id).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "the contest has score, can't not delete"})
			return
		}

		_ = svc.DB.Delete(&contest)
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
