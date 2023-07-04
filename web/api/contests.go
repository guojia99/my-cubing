/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:28.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"

	"my-cubing/db"
)

type (
	CreateContestRequest struct {
		Contest
	}

	Contest struct {
		ID          uint              `json:"ID"`
		Name        string            `json:"Name"`
		Description string            `json:"Description"`
		Routes      db.ContestContent `json:"Content"`
		StartTime   int64             `json:"StartTime"`
		EndTime     int64             `json:"EndTime"`
	}

	ReadContestsResponse struct {
		Contests []Contest `json:"Contests"`
	}
)

func CreateContest(ctx *gin.Context) {
	var req CreateContestRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if req.StartTime == 0 || req.EndTime == 0 {
		req.StartTime = time.Now().Unix()
		req.EndTime = time.Now().Unix() + int64((time.Hour * 24 * 7).Seconds())
	}

	routes, _ := json.Marshal(req.Routes)
	c := &db.Contest{
		Model:       db.Model{},
		Name:        req.Name,
		Content:     string(routes),
		Description: req.Description,
		StartTime:   time.Unix(req.StartTime, 0),
		EndTime:     time.Unix(req.EndTime, 0),
	}
	err := db.DB.Create(c).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": c.ID})
}

func ReadContests(ctx *gin.Context) {
	var contests []db.Contest
	if err := db.DB.Model(&db.Contest{}).Order("created_at DESC").Find(&contests).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var resp = new(ReadContestsResponse)
	for _, contest := range contests {
		var route db.ContestContent
		_ = json.UnmarshalFromString(contest.Content, &route)
		resp.Contests = append(resp.Contests, Contest{
			ID:          contest.ID,
			Name:        contest.Name,
			Routes:      route,
			Description: contest.Description,
			StartTime:   contest.StartTime.Unix(),
			EndTime:     contest.EndTime.Unix(),
		})
	}
	ctx.JSON(http.StatusOK, resp)
}

func UpdateContest(ctx *gin.Context) {

}

func DeleteContest(ctx *gin.Context) {

}
