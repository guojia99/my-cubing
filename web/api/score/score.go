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
}

func CreateScore(ctx *gin.Context) {
	// 0. 确认入参是否有效
	var req CreateScoreRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	pj := db.StrToProject(req.ProjectName)
	if pj == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "该项目类型错误"})
		return
	}

	// 1. 查看比赛是否有效
	var contest db.Contest
	if err := db.DB.Where("id = ?", req.ContestID).First(&contest).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if contest.IsEnd {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "比赛已结束"})
		return
	}

	// 2. 查看玩家是否存在，不存在创建一个
	var player = db.Player{Name: req.PlayerName}
	if err := db.DB.FirstOrCreate(&player, "name = ?", req.PlayerName).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// 3. 找到尝试找到本场比赛成绩, 找到就直接覆盖, 或者插入成绩
	var score db.Score
	err := db.DB.Model(&db.Score{}).
		Where("player_id = ?", player.ID).Where("contest_id = ?", req.ContestID).
		Where("route_number = ?", 1).Where("project = ?", pj).First(&score).Error

	if err != nil || score.ID == 0 {
		score = db.Score{
			CreatedAt:   time.Now(),
			PlayerID:    player.ID,
			ContestID:   req.ContestID,
			RouteNumber: req.RouteNumber,
			Project:     db.StrToProject(req.ProjectName),
		}
	}
	if err = score.SetResult(req.Results); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = db.DB.Save(&score).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
	// 4. 全D直接返回
	if score.Best == 0 {
		return
	}

	// 5. 判断是否是最佳成绩
	var (
		best    db.Score
		bestAvg db.Score
	)
	// 查之前所有的成绩
	err = db.DB.Where("player_id = ?", player.ID).Where("project = ?", pj).Where("best != ?", 0).Where("id != ?", score.ID).Order("best").First(&best).Error
	if ((err != nil || best.Best == 0) && score.Best != 0) || (score.IsBestScore(best)) {
		// 之前没有成绩, 且当前有成绩
		// 之前有成绩, 当前成绩好
		score.IsBest = true
		db.DB.Save(&score)
	}

	err = db.DB.Where("player_id = ?", player.ID).Where("project = ?", pj).Where("avg != ?", 0).Where("id != ?", score.ID).Order("avg").First(&bestAvg).Error
	if ((err != nil || best.Avg == 0) && score.Avg != 0) || score.IsBestAvgScore(bestAvg) {
		score.IsBestAvg = true
		db.DB.Save(&score)
	}
}

func DeleteScore(ctx *gin.Context) {
	var req DeleteScoreRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// 0. 项目
	pj := db.StrToProject(req.Project)
	if pj == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "该项目类型错误"})
		return
	}

	// 1. 查看比赛是否有效
	var contest db.Contest
	if err := db.DB.Where("id = ?", req.ContestID).First(&contest).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if contest.IsEnd {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "比赛已结束"})
		return
	}

	// 2. 查看玩家是否存在
	var player = db.Player{Name: req.PlayerName}
	if err := db.DB.First(&player, "name = ?", req.PlayerName).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// 3. 查成绩
	var score db.Score
	if err := db.DB.Where("player_id = ?", player.ID).Where("contest_id = ?", contest.ID).Where("project = ?", pj).First(&score).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = db.DB.Delete(&score).Error
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

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
	for _, project := range db.WCAProjectRoute() {
		if project == db.Cube333MBF {
			continue
		}

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
			bestScore.AvgScore = avg.Avg
		}
		out.Data = append(out.Data, bestScore)
	}

	var allMbfScore []db.Score
	if err := db.DB.Where("project = ?", db.Cube333MBF).Where("r1 != ?", 0).Find(&allMbfScore).Error; err == nil && len(allMbfScore) > 0 {
		sort.Slice(allMbfScore, func(i, j int) bool { return allMbfScore[i].IsBestScore(allMbfScore[j]) })
		best := allMbfScore[0]
		var bestPlayer db.Player
		if err = db.DB.Where("id = ?", allMbfScore[0].PlayerID).First(&bestPlayer).Error; err == nil {
			out.Data = append(out.Data, BestScore{Project: db.Cube333MBF.Cn(), BestPlayer: bestPlayer.Name, MBFScore: fmt.Sprintf("%d/%d", int(best.Result1), int(best.Result2)), BestScore: best.Result3})
		}
	}

	caches.Add(key, out, time.Minute*5)
	ctx.JSON(http.StatusOK, out)
}

// GetProjectScores 获取所有成绩
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
	for _, project := range db.WCAProjectRoute() {
		out.ProjectList = append(out.ProjectList, project.Cn())
		out.Best[project.Cn()] = make([]ProjectScores, 0)
		out.Avg[project.Cn()] = make([]ProjectScores, 0)
	}

	// 2. 这个角色查所有的项目最佳成绩
	for _, player := range players {
		for _, project := range db.WCAProjectRoute() {
			// 多盲需要独立查询
			if project == db.Cube333MBF {
				continue
			}
			var (
				best db.Score
				avg  db.Score
			)
			if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("best != ?", 0).Order("best").First(&best).Error; err == nil {
				out.Best[project.Cn()] = append(out.Best[project.Cn()], ProjectScores{
					Player: player.Name,
					Score:  best,
				})
			}

			if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("avg != ?", 0).Order("avg").First(&avg).Error; err == nil {
				out.Avg[project.Cn()] = append(out.Avg[project.Cn()], ProjectScores{
					Player: player.Name,
					Score:  avg,
				})
			}
		}
		// 查多盲, 只要有一把成了, 就可以记录进去成绩
		var mbfScore []db.Score
		if err := db.DB.Where("project = ?", db.Cube333MBF).Where("player_id = ?", player.ID).Where("r1 != ?", 0).Find(&mbfScore).Error; err == nil && len(mbfScore) > 0 {
			sort.Slice(mbfScore, func(i, j int) bool { return mbfScore[i].IsBestScore(mbfScore[j]) })
			out.Best[db.Cube333MBF.Cn()] = append(out.Best[db.Cube333MBF.Cn()], ProjectScores{
				Player: player.Name,
				Score:  mbfScore[0],
			})
		}
	}

	// 3. 给所有的项目排序
	for _, project := range db.WCAProjectRoute() {
		sort.Slice(out.Best[project.Cn()], func(i, j int) bool {
			return out.Best[project.Cn()][i].IsBestScore(out.Best[project.Cn()][j].Score)
		})
		sort.Slice(out.Avg[project.Cn()], func(i, j int) bool {
			return out.Avg[project.Cn()][i].IsBestAvgScore(out.Avg[project.Cn()][j].Score)
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

	for _, project := range db.WCAProjectRoute() {
		bestCache[project] = make([]SorScoreDetail, 0)
		avgCache[project] = make([]SorScoreDetail, 0)

		for _, player := range players {
			var (
				best db.Score
				avg  db.Score
			)

			if project != db.Cube333MBF {
				if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("best != ?", 0).Order("best").First(&best).Error; err == nil {
					bestCache[project] = append(bestCache[project], SorScoreDetail{Player: player.Name, Score: best})
				}

				if err := db.DB.Where("project = ?", project).Where("player_id = ?", player.ID).Where("avg != ?", 0).Order("avg").First(&avg).Error; err == nil {
					avgCache[project] = append(avgCache[project], SorScoreDetail{Player: player.Name, Score: avg})
				}
				continue
			}

			// 独立查询多盲
			var mbfScore []db.Score
			if err := db.DB.Where("project = ?", db.Cube333MBF).Where("player_id = ?", player.ID).Where("r1 != ?", 0).Find(&mbfScore).Error; err == nil && len(mbfScore) > 0 {
				sort.Slice(mbfScore, func(i, j int) bool { return mbfScore[i].IsBestScore(mbfScore[j]) })
				bestCache[project] = append(bestCache[project], SorScoreDetail{Player: player.Name, Score: mbfScore[0]})
			}
		}
	}

	// 3. 排序
	for _, project := range db.WCAProjectRoute() {
		sort.Slice(bestCache[project], func(i, j int) bool {
			return bestCache[project][i].IsBestScore(bestCache[project][j].Score)
		})
		sort.Slice(avgCache[project], func(i, j int) bool {
			return avgCache[project][i].IsBestAvgScore(avgCache[project][j].Score)
		})
	}

	// 4. 给各个玩家汇总
	var playerCache = make(map[string][]SorScoreDetail)
	for _, player := range players {
		playerCache[player.Name] = []SorScoreDetail{{Player: player.Name, Count: 0}, {Player: player.Name, Count: 0}}
	}

	for _, project := range db.WCAProjectRoute() {
		for _, player := range players {
			var bestAdd bool
			for idx, val := range bestCache[project] {
				if val.Player == player.Name {
					playerCache[val.Player][0].Count += idx + 1
					bestAdd = true
				}
			}
			if !bestAdd {
				playerCache[player.Name][0].Count += len(bestCache[project]) + 1
			}

			if project == db.Cube333MBF {
				continue
			}
			var avgAdd bool
			for idx, val := range avgCache[project] {
				if val.Player == player.Name {
					playerCache[val.Player][1].Count += idx + 1
					avgAdd = true
				}
			}
			if !avgAdd {
				playerCache[player.Name][1].Count += len(avgCache[project]) + 1
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

	key := fmt.Sprintf("GetContestScores_%d", req.ContestID)
	if data, ok := caches.Get(key); ok {
		ctx.JSON(http.StatusOK, data)
		return
	}

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
			Player: player.Name,
			Score:  score,
		})
	}

	// 4. 整理并输出
	var out = GetContestScoresResponse{
		ProjectList: make([]string, 0),
		Players:     make([]GetContestScoresPlayer, 0),
		Data:        make(map[string][]GetContestScoresDetail),
	}

	var hasProject = make(map[db.Project]struct{})
	for k, _ := range scoreCache {
		hasProject[k] = struct{}{}
	}

	for _, val := range db.WCAProjectRoute() {
		if _, ok := hasProject[val]; ok {
			out.ProjectList = append(out.ProjectList, val.Cn())
		}
	}

	for project, _ := range hasProject {
		ss := scoreCache[project]
		sort.Slice(ss, func(i, j int) bool {
			switch ss[i].Project {
			case db.Cube333MBF, db.Cube333BF, db.Cube444BF, db.Cube555BF: // 多盲规则
				return ss[i].IsBestScore(ss[j].Score)
			default:
				if ss[i].Avg+ss[j].Avg == 0 { // 都没有平均成绩
					return ss[i].IsBestScore(ss[j].Score)
				}

				if ss[i].Avg == 0 && ss[j].Avg != 0 {
					return false
				}
				if ss[i].Avg != 0 && ss[j].Avg == 0 {
					return true
				}

				return ss[i].IsBestAvgScore(ss[j].Score)
			}
		})
		out.Data[project.Cn()] = append(out.Data[project.Cn()], ss...)
	}
	out.ContestName = contest.Name

	for _, val := range playerMap {
		out.Players = append(out.Players, val)
	}

	caches.Add(key, out, time.Minute*5)
	ctx.JSON(http.StatusOK, out)
}

// EndContestScore 结束一场比赛
func EndContestScore(ctx *gin.Context) {
	var req EndContestScoreRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	var contest db.Contest
	if err := db.DB.Where("id = ?", req.ContestID).First(&contest).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	// 1. 查所有的该场比赛成绩, 把最佳成绩统计出来, 查所有之前的比赛, 把以往最佳成绩统计出来
	var (
		thisContestBest = make(map[db.Project]db.Score)
		thisContestAvg  = make(map[db.Project]db.Score)

		oldContestBest = make(map[db.Project]db.Score)
		oldContestAvg  = make(map[db.Project]db.Score)
	)
	for _, project := range db.ProjectList() {
		var bestScore db.Score
		if err := db.DB.Where("contest_id = ?", req.ContestID).Where("project = ?", project).Where("best != ?", 0).Order("best").First(&bestScore).Error; err == nil {
			thisContestBest[project] = bestScore
		}
		var avgScore db.Score
		if err := db.DB.Where("contest_id = ?", req.ContestID).Where("project = ?", project).Where("avg != ?", 0).Order("avg").First(&avgScore).Error; err == nil {
			thisContestAvg[project] = avgScore
		}

		var oldBestScore db.Score
		if err := db.DB.Where("contest_id != ?", req.ContestID).Where("project = ?", project).Where("best != ?", 0).Order("best").First(&oldBestScore).Error; err == nil {
			oldContestBest[project] = oldBestScore
		}
		var oldAvgScore db.Score
		if err := db.DB.Where("contest_id != ?", req.ContestID).Where("project = ?", project).Where("avg != ?", 0).Order("avg").First(&oldAvgScore).Error; err == nil {
			oldContestAvg[project] = oldAvgScore
		}
	}

	// 2. 循环所有当前的最佳记录
	for key, score := range thisContestBest {
		// 旧的成绩存在且新成绩差
		if _, ok := oldContestBest[key]; ok && oldContestBest[key].Best < score.Best {
			continue
		}
		// 不存在或新成绩好
		score.IsBestRecord = true
		db.DB.Save(&score)
	}

	for key, score := range thisContestAvg {
		// 旧的成绩存在且新成绩差
		if _, ok := oldContestAvg[key]; ok && oldContestAvg[key].Avg < score.Avg {
			continue
		}
		// 不存在或新成绩好
		score.IsAvgRecord = true
		db.DB.Save(&score)
	}

	contest.IsEnd = true
	contest.EndTime = time.Now()
	db.DB.Save(&contest)
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
