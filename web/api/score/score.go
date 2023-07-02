/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/1 下午3:38.
 *  * Author: guojia(https://github.com/guojia99)
 */

package score

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/cache"

	"my-cubing/db"
)

// todo 加缓存

var caches = cache.NewLRUExpireCache(100)

// GetAllProjectBestScore 获取所有成绩最佳
func GetAllProjectBestScore(ctx *gin.Context) {
	key := "GetAllProjectBestScore"
	if data, ok := caches.Get(key); ok {
		ctx.JSON(http.StatusOK, data)
		return
	}

	var out = GetAllProjectBestScoreResponse{
		Data: make([]BestScore, 0),
	}
	for _, project := range db.ProjectList() {
		var best db.Score
		var avg db.Score
		bestScore := BestScore{Project: project.Cn(), AvgPlayer: "-", BestPlayer: "-"}
		if err := db.DB.Where("best != ?", 0).Where("project = ?", project).Order("best").First(&best).Error; err == nil {
			var bestPlayer db.Player
			if err = db.DB.Where("id = ?", best.PlayerID).First(&bestPlayer).Error; err == nil {
				bestScore.BestPlayer = bestPlayer.Name
			}
			bestScore.BestScore = best.Best
		}
		if err := db.DB.Where("avg != ?", 0).Where("project = ?", project).Order("avg").First(&avg).Error; err == nil {
			var avgPlayer db.Player
			if err = db.DB.Where("id = ?", avg.PlayerID).First(&avgPlayer).Error; err == nil {
				bestScore.AvgPlayer = avgPlayer.Name
			}
			bestScore.AvgScore = best.Avg
		}
		out.Data = append(out.Data, bestScore)
	}

	caches.Add(key, out, time.Minute*5)
	ctx.JSON(http.StatusOK, out)
}

// GetProjectScores 获取所有成绩最佳排名
func GetProjectScores(ctx *gin.Context) {
	key := "GetProjectScores"
	if data, ok := caches.Get(key); ok {
		ctx.JSON(http.StatusOK, data)
		return
	}

	// 1. 查所有的角色
	var players []db.Player
	db.DB.Find(&players)
	var out = GetProjectScoresResponse{
		ProjectList: []string{},
		Avg:         make(map[string][]ProjectScores),
		Best:        make(map[string][]ProjectScores),
	}
	for _, project := range db.ProjectList() {
		out.ProjectList = append(out.ProjectList, project.Cn())
		out.Best[project.Cn()] = make([]ProjectScores, 0)
		out.Avg[project.Cn()] = make([]ProjectScores, 0)
	}

	// 2. 这个角色查所有的项目最佳成绩
	for _, player := range players {
		for _, project := range db.ProjectList() {
			var (
				best db.Score
				avg  db.Score
			)

			if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("best != ?", 0).Order("best").First(&best).Error; err == nil {
				out.Best[project.Cn()] = append(out.Best[project.Cn()], ProjectScores{
					Player: player.Name,
					Score:  best.Best,
				})
			}

			if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("avg != ?", 0).Order("avg").First(&avg).Error; err == nil {
				out.Avg[project.Cn()] = append(out.Avg[project.Cn()], ProjectScores{
					Player: player.Name,
					Score:  avg.Avg,
				})
			}
		}
	}

	// 3. 给所有的项目排序
	for _, project := range db.ProjectList() {
		sort.Slice(out.Best[project.Cn()], func(i, j int) bool {
			return out.Best[project.Cn()][i].Score < out.Best[project.Cn()][j].Score
		})
		sort.Slice(out.Avg[project.Cn()], func(i, j int) bool {
			return out.Avg[project.Cn()][i].Score < out.Avg[project.Cn()][j].Score
		})
	}

	caches.Add(key, out, time.Minute*5)
	ctx.JSON(http.StatusOK, out)
}

// GetSorScores 查询所有角色的排名最佳排名总和
func GetSorScores(ctx *gin.Context) {
	key := "GetSorScores"
	if data, ok := caches.Get(key); ok {
		ctx.JSON(http.StatusOK, data)
		return
	}

	// 1. 查所有角色
	var players []db.Player
	db.DB.Find(&players)

	// 2. 查各个项目的所有角色最佳成绩
	var (
		bestCache = make(map[db.Project][]SorScoreDetail)
		avgCache  = make(map[db.Project][]SorScoreDetail)
	)

	for _, project := range db.ProjectList() {
		bestCache[project] = make([]SorScoreDetail, 0)
		avgCache[project] = make([]SorScoreDetail, 0)

		for _, player := range players {
			var (
				best db.Score
				avg  db.Score
			)

			if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("best != ?", 0).Order("best").First(&best).Error; err == nil {
				bestCache[project] = append(bestCache[project], SorScoreDetail{Player: player.Name, Value: best.Best})
			}

			if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("avg != ?", 0).Order("avg").First(&avg).Error; err == nil {
				avgCache[project] = append(avgCache[project], SorScoreDetail{Player: player.Name, Value: best.Best})
			}
		}
	}

	// 3. 排序
	for _, project := range db.ProjectList() {
		sort.Slice(bestCache[project], func(i, j int) bool {
			return bestCache[project][i].Value < bestCache[project][j].Value
		})
		sort.Slice(avgCache[project], func(i, j int) bool {
			return avgCache[project][i].Value < avgCache[project][j].Value
		})
	}

	// 4. 给各个玩家汇总
	var playerCache = make(map[string][]SorScoreDetail)
	for _, player := range players {
		playerCache[player.Name] = []SorScoreDetail{{Player: player.Name, Count: 0}, {Player: player.Name, Count: 0}}
	}

	for _, project := range db.ProjectList() {
		for _, player := range players {
			var bestAdd bool
			var avgAdd bool

			for idx, val := range bestCache[project] {
				if val.Player == player.Name {
					playerCache[val.Player][0].Count += idx + 1
					bestAdd = true
				}
			}

			for idx, val := range avgCache[project] {
				if val.Player == player.Name {
					playerCache[val.Player][1].Count += idx + 1
					avgAdd = true
				}
			}

			if !bestAdd {
				playerCache[player.Name][0].Count += len(players)
			}
			if !avgAdd {
				playerCache[player.Name][1].Count += len(players)
			}
		}
	}

	// 5. 统计结果
	var out = GetSorScoresResponse{
		Best: make([]SorScoreDetail, 0),
		Avg:  make([]SorScoreDetail, 0),
	}
	for _, val := range playerCache {
		out.Best = append(out.Best, val[0])
		out.Avg = append(out.Avg, val[1])
	}

	sort.Slice(out.Best, func(i, j int) bool { return out.Best[i].Count < out.Best[j].Count })
	sort.Slice(out.Avg, func(i, j int) bool { return out.Avg[i].Count < out.Avg[j].Count })

	caches.Add(key, out, time.Minute*5)
	ctx.JSON(http.StatusOK, out)
}

// GetContestScores 获取单场比赛成绩汇总
func GetContestScores(ctx *gin.Context) {
	var req GetContestScoresRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// todo 加缓存

	// 1. 获取比赛的信息
	var contest db.Contest
	if err := db.DB.Where("id = ?", req.ContestID).First(&contest).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	// 2. 获取这场比赛所有的数据
	var scores []db.Score
	if err := db.DB.Where("contest_id = ?", req.ContestID).Find(&scores).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		fmt.Println(err)
		return
	}

	// 3. 给各个项目分类排序
	var playerMap = make(map[string]GetContestScoresPlayer)
	var scoreCache = make(map[db.Project][]GetContestScoresDetail)
	for _, score := range scores {
		if _, ok := scoreCache[score.Project]; !ok {
			scoreCache[score.Project] = make([]GetContestScoresDetail, 0)
		}

		var player db.Player
		if err := db.DB.Where("id = ?", score.PlayerID).First(&player).Error; err != nil {
			continue
		}
		playerMap[player.Name] = GetContestScoresPlayer{Name: player.Name, Id: player.ID}

		scoreCache[score.Project] = append(scoreCache[score.Project], GetContestScoresDetail{
			Player:    player.Name,
			Result1:   score.Result1,
			Result2:   score.Result2,
			Result3:   score.Result3,
			Result4:   score.Result4,
			Result5:   score.Result5,
			Best:      score.Best,
			Avg:       score.Avg,
			IsBest:    score.IsBest,
			IsBestAvg: score.IsBestAvg,
		})
	}

	// 4. 整理并输出
	var out = GetContestScoresResponse{
		ProjectList: make([]string, 0),
		Players:     make([]GetContestScoresPlayer, 0),
		Data:        make(map[string][]GetContestScoresDetail),
	}

	var hasProject []db.Project
	for key, _ := range scoreCache {
		hasProject = append(hasProject, key)
	}
	sort.Slice(hasProject, func(i, j int) bool { return hasProject[i] < hasProject[j] })
	for _, val := range hasProject {
		out.ProjectList = append(out.ProjectList, val.Cn())
	}

	for _, project := range hasProject {
		ss := scoreCache[project]
		sort.Slice(ss, func(i, j int) bool {
			iHasAvg := ss[i].Avg != 0
			jHasAvg := ss[i].Avg != 0

			// 一方有平均， 另一方无平均, 有平均排前
			if (iHasAvg && !jHasAvg) || (jHasAvg && !iHasAvg) {
				return iHasAvg
			}
			// 双方都有平均, 小的排前, 相同则最佳成绩的排前
			if iHasAvg && jHasAvg {
				if ss[i].Avg == ss[i].Avg {
					return ss[i].Best < ss[j].Best
				}
				return ss[i].Avg < ss[j].Avg
			}
			// 双方都无平均, 按最佳成绩排
			return ss[i].Best < ss[j].Best
		})
		out.Data[project.Cn()] = append(out.Data[project.Cn()], ss...)
	}
	out.ContestName = contest.Name

	for _, val := range playerMap {
		out.Players = append(out.Players, val)
	}

	ctx.JSON(http.StatusOK, out)
}
