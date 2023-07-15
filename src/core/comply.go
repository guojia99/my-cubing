/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"fmt"
	"sort"
	"time"

	"github.com/guojia99/my-cubing/src/core/model"
)

// addScore 添加一条成绩
func (c *client) addScore(playerName string, contestID uint, project model.Project, routeNum int, result []float64) (err error) {
	// 1. 确定比赛是否存在
	var contest model.Contest
	if err = c.db.Where("id = ?", contestID).First(&contest).Error; err != nil || contest.IsEnd {
		return fmt.Errorf("the contest id end or error %+v", err)
	}

	// 2. 获取轮次信息
	var round model.Round
	if err = c.db.Where("contest_id = ?", contestID).Where("project = ?", project).Where("number = ?", routeNum).First(&round).Error; err != nil {
		return err
	}

	// 3. 玩家信息
	var player = model.Player{Name: playerName}
	if err = c.db.FirstOrCreate(&player, "name = ?", playerName).Error; err != nil {
		return err
	}

	// 4. 尝试找到本场比赛成绩
	var score model.Score
	err = c.db.Model(&model.Score{}).
		Where("player_id = ?", player.ID).Where("contest_id = ?", contestID).
		Where("route_id = ?", round.ID).First(&score).Error

	if err != nil || score.ID == 0 {
		score = model.Score{
			PlayerID:   player.ID,
			PlayerName: playerName,
			ContestID:  contestID,
			RouteID:    round.ID,
			Project:    project,
		}
	}

	if err = score.SetResult(result); err != nil {
		return err
	}

	if err = c.db.Save(&score).Error; err != nil {
		return err
	}

	// 5. 最佳成绩查询, 确定是否该玩家刷新了最佳成绩
	if score.Best == 0 {
		return nil
	}
	var (
		bestSingle model.Score
		bestAvg    model.Score
	)
	err = c.db.Where("player_id = ?", player.ID).Where("project = ?", project).
		Where("best != ?", 0).Where("id != ?", score.ID).Order("best").First(&bestSingle).Error
	if ((err != nil || bestSingle.Best == 0) && score.Best != 0) || (score.IsBestScore(bestSingle)) {
		// 之前没有成绩, 且当前有成绩
		// 之前有成绩, 当前成绩好
		score.IsBestSingle = true
		c.db.Save(&score)
	}

	err = c.db.Where("player_id = ?", player.ID).Where("project = ?", project).
		Where("avg != ?", 0).Where("id != ?", score.ID).Order("avg").First(&bestAvg).Error
	if ((err != nil || bestAvg.Avg == 0) && score.Avg != 0) || score.IsBestAvgScore(bestAvg) {
		score.IsBestAvg = true
		c.db.Save(&score)
	}
	return nil
}

// removeScoreByContestID 删除一条成绩
func (c *client) removeScoreByContestID(playerName string, contestID uint, project model.Project, routeNum int) (err error) {
	// 1. 确定比赛是否存在
	var contest model.Contest
	if err = c.db.Where("id = ?", contestID).First(&contest).Error; err != nil || contest.IsEnd {
		return fmt.Errorf("the contest id end or error %+v", err)
	}

	// 2. 获取轮次信息
	var round model.Round
	if err = c.db.Where("contest_id = ?", contestID).Where("project = ?", project).Where("number = ?", routeNum).First(&round).Error; err != nil {
		return err
	}

	// 3. 玩家信息
	var player = model.Player{Name: playerName}
	if err = c.db.FirstOrCreate(&player, "name = ?", playerName).Error; err != nil {
		return err
	}

	// 4. 尝试找到本场比赛成绩
	var score model.Score
	err = c.db.Model(&model.Score{}).
		Where("player_id = ?", player.ID).Where("contest_id = ?", contestID).
		Where("route_id = ?", round.ID).First(&score).Error

	if err != nil {
		return err
	}

	return c.db.Delete(&score).Error
}

// statisticalRecordsAndEndContest 结束一场比赛并获取记录
func (c *client) statisticalRecordsAndEndContest(contestID uint) (err error) {
	// 1. 确定比赛是否存在 且非结束的
	var contest model.Contest
	if err = c.db.Where("id = ?", contestID).First(&contest).Error; err != nil || contest.IsEnd {
		return fmt.Errorf("the contest id end or error %+v", err)
	}

	// 2. 获取本场比赛最佳
	thisContestBestSingle, thisContestBestAvg := c.getContestBestSingle(contestID, false), c.getContestBestAvg(contestID, false)
	oldContestBest, oldContestAvg := c.getContestBestSingle(contestID, true), c.getContestBestAvg(contestID, true)

	var records []model.Record
	for key, score := range thisContestBestSingle {
		if _, ok := oldContestBest[key]; ok && score.IsBestScore(oldContestBest[key]) {
			records = append(records, model.Record{
				RType:      model.RecordBySingle,
				ScoreId:    score.ID,
				PlayerID:   score.PlayerID,
				PlayerName: score.PlayerName,
				ContestID:  score.ContestID,
			})
		}
	}

	for key, score := range thisContestBestAvg {
		if _, ok := oldContestAvg[key]; ok && score.IsBestAvgScore(oldContestAvg[key]) {
			records = append(records, model.Record{
				RType:      model.RecordByAvg,
				ScoreId:    score.ID,
				PlayerID:   score.PlayerID,
				PlayerName: score.PlayerName,
				ContestID:  score.ContestID,
			})
		}
	}
	if err = c.db.Save(&records).Error; err != nil {
		return err
	}

	// 3. 结束比赛
	contest.IsEnd = true
	contest.EndTime = time.Now()
	return c.db.Save(&contest).Error
}

// getBestScores 获取所有成绩中最佳成绩
func (c *client) getBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	for _, project := range model.WCAProjectRoute() {
		var best, avg model.Score
		if project == model.Cube333MBF {
			if err := c.db.Where("project = ?", project).Where("r1 != ?", 0).
				Order("best").Order("r2 DESC").Order("r3").First(&best).Error; err == nil {
				bestSingle[project] = best
			}
			continue
		}

		if err := c.db.Where("best != ?", 0).Where("project = ?", project).Order("best").First(&best).Error; err == nil {
			bestSingle[project] = best
		}
		if err := c.db.Where("avg != ?", 0).Where("project = ?", project).Order("avg").First(&avg).Error; err == nil {
			bestAvg[project] = avg
		}
	}
	return
}

// getAllPlayerBestScore 获取所有玩家各自的全项目最佳成绩
func (c *client) getAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	bestSingle, bestAvg = make(map[model.Project][]model.Score), make(map[model.Project][]model.Score)

	var players []model.Player
	c.db.Find(&players)

	for _, project := range model.WCAProjectRoute() {
		bestSingle[project] = make([]model.Score, 0)
		bestAvg[project] = make([]model.Score, 0)
	}

	for _, project := range model.WCAProjectRoute() {
		for _, player := range players {
			var best, avg model.Score
			if project == model.Cube333MBF {
				if err := c.db.Where("player_id = ?", player.ID).Where("project = ?", project).Where("r1 != ?", 0).
					Order("best").Order("r2 DESC").Order("r3").First(&best).Error; err == nil {
					bestSingle[project] = append(bestSingle[project], best)
				}
				continue
			}
			if err := c.db.Where("player_id = ?", player.ID).Where("project = ?", project).Where("best != ?", 0).Order("best").First(&best).Error; err == nil {
				bestSingle[project] = append(bestSingle[project], best)
			}
			if err := c.db.Where("player_id = ?", player.ID).Where("project = ?", project).Where("avg != ?", 0).Order("avg").First(&avg).Error; err == nil {
				bestAvg[project] = append(bestAvg[project], avg)
			}
		}

		// sort
		sort.Slice(bestSingle[project], func(i, j int) bool { return bestSingle[project][i].IsBestScore(bestSingle[project][j]) })
		sort.Slice(bestAvg[project], func(i, j int) bool { return bestAvg[project][i].IsBestScore(bestAvg[project][j]) })
	}

	return
}

// getSorScore 获取所有玩家的Sor排名
func (c *client) getSorScore() (single, avg []SorScore) {
	var players []model.Player
	c.db.Find(&players)
	bestSingle, bestAvg := c.getAllPlayerBestScore()
	var playerCache = make(map[uint]*SorScore)

	for _, player := range players {
		playerCache[player.ID] = &SorScore{Player: player}
		for _, project := range model.WCAProjectRoute() {
			var bestUse, avgUse bool
			// best
			for idx, val := range bestSingle[project] {
				if val.PlayerID == player.ID {
					playerCache[val.PlayerID].SingleCount += int64(idx + 1)
					bestUse = true
					break
				}
			}

			// avg
			for idx, val := range bestAvg[project] {
				if val.PlayerID == player.ID {
					playerCache[val.PlayerID].AvgCount += int64(idx + 1)
					avgUse = true
					break
				}
			}

			if !bestUse {
				playerCache[player.ID].SingleCount += int64(len(players))
			}
			if !avgUse {
				playerCache[player.ID].AvgCount += int64(len(players))
			}
		}
	}

	for _, val := range playerCache {
		single = append(single, SorScore{Player: val.Player, SingleCount: val.SingleCount})
		avg = append(avg, SorScore{Player: val.Player, SingleCount: val.AvgCount})
	}

	sort.Slice(single, func(i, j int) bool { return single[i].SingleCount < single[j].SingleCount })
	sort.Slice(avg, func(i, j int) bool { return avg[i].AvgCount < avg[j].AvgCount })
	return
}

// getScoreByContest 获取一个比赛所有成绩
func (c *client) getScoreByContest(contestID uint) map[model.Project][]RoutesScores {
	var out = make(map[model.Project][]RoutesScores)

	for _, project := range model.WCAProjectRoute() {
		// 查询该比赛
		var scores []model.Score
		c.db.Where("contest_id = ?", contestID).Where("project = ?", project).Find(&scores)
		if len(scores) == 0 {
			continue
		}

		// 按轮次分类
		var routeCache = make(map[uint][]model.Score)
		for _, score := range scores {
			if _, ok := routeCache[score.RouteID]; !ok {
				routeCache[score.RouteID] = make([]model.Score, 0)
			}
			routeCache[score.RouteID] = append(routeCache[score.RouteID], score)
		}

		// 查出并给予相应的查询方法
		out[project] = make([]RoutesScores, 0)
		for key, val := range routeCache {
			var round model.Round
			_ = c.db.Where("id = ?", key).First(&round)
			out[project] = append(out[project], RoutesScores{Round: round, Scores: val})
		}
		sort.Slice(out[project], func(i, j int) bool { return out[project][i].Round.Number < out[project][j].Round.Number })
	}
	return out
}

// getSorScoreByContest 获取比赛的Sor排名
func (c *client) getSorScoreByContest(contestID uint) (single, avg []SorScore) {
	// 查这场比赛所有选手
	var playerIDs []uint64
	c.db.Model(&model.Score{}).Distinct("player_id").Where("contest_id = ?", contestID).Pluck("player_id", &playerIDs)
	if len(playerIDs) == 0 {
		return
	}
	var players []model.Player
	c.db.Where("id in ?", playerIDs).Find(&players)

	// 查询这个比赛所有角色的最佳成绩
	var (
		bestSingleCache = make(map[model.Project][]model.Score)
		bestAvgCache    = make(map[model.Project][]model.Score)
	)

	for _, project := range model.WCAProjectRoute() {
		bestSingleCache[project] = make([]model.Score, 0)
		bestAvgCache[project] = make([]model.Score, 0)
		for _, player := range players {
			var b, a model.Score
			if project == model.Cube333MBF {
				if err := c.db.Where("player_id = ?", player.ID).Where("project = ?", project).Where("r1 != ?", 0).
					Order("best").Order("r2 DESC").Order("r3").First(&b).Error; err == nil {
					bestSingleCache[project] = append(bestSingleCache[project], b)
				}
				continue
			}
			if err := c.db.Where("player_id = ?", player.ID).Where("project = ?", project).Where("best != ?", 0).Order("best").First(&b).Error; err == nil {
				bestSingleCache[project] = append(bestSingleCache[project], b)
			}
			if err := c.db.Where("player_id = ?", player.ID).Where("project = ?", project).Where("avg != ?", 0).Order("avg").First(&a).Error; err == nil {
				bestAvgCache[project] = append(bestAvgCache[project], a)
			}
		}
	}

	// 排序
	var playerCache = make(map[uint]*SorScore)
	for _, player := range players {
		playerCache[player.ID] = &SorScore{Player: player}
		for _, project := range model.WCAProjectRoute() {
			var bestUse, avgUse bool
			// best
			for idx, val := range bestSingleCache[project] {
				if val.PlayerID == player.ID {
					playerCache[val.PlayerID].SingleCount += int64(idx + 1)
					bestUse = true
					break
				}
			}
			// avg
			for idx, val := range bestAvgCache[project] {
				if val.PlayerID == player.ID {
					playerCache[val.PlayerID].AvgCount += int64(idx + 1)
					avgUse = true
					break
				}
			}

			if !bestUse {
				playerCache[player.ID].SingleCount += int64(len(players))
			}
			if !avgUse {
				playerCache[player.ID].AvgCount += int64(len(players))
			}
		}
	}
	for _, val := range playerCache {
		single = append(single, SorScore{Player: val.Player, SingleCount: val.SingleCount})
		avg = append(avg, SorScore{Player: val.Player, SingleCount: val.AvgCount})
	}

	sort.Slice(single, func(i, j int) bool { return single[i].SingleCount < single[j].SingleCount })
	sort.Slice(avg, func(i, j int) bool { return avg[i].AvgCount < avg[j].AvgCount })
	return
}

// getPlayerScore 获取玩家所有成绩
func (c *client) getPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scoresByContest []ScoresByContest) {
	var scores []model.Score
	c.db.Where("player_id = ?", playerID).Find(&scores)
	if len(scores) == 0 {
		return
	}

	var (
		cache     = make(map[uint][]model.Score)
		avgCache  = make(map[model.Project]model.Score)
		bestCache = make(map[model.Project]model.Score)
	)
	for _, score := range scores {
		if _, ok := cache[score.ContestID]; !ok {
			cache[score.ContestID] = make([]model.Score, 0)
		}
		cache[score.ContestID] = append(cache[score.ContestID], score)

		if got, ok := bestCache[score.Project]; !ok || score.IsBestScore(got) {
			bestCache[score.Project] = score
		}
		if got, ok := avgCache[score.Project]; !ok || score.IsBestAvgScore(got) {
			avgCache[score.Project] = score
		}
	}

	for key, val := range cache {
		var contest model.Contest
		c.db.Where("id = ?", key).First(&contest)
		scoresByContest = append(scoresByContest, ScoresByContest{
			Contest: contest,
			Scores:  val,
		})
	}

	for _, val := range avgCache {
		bestAvg = append(bestAvg, val)
	}
	for _, val := range bestCache {
		bestSingle = append(bestSingle, val)
	}
	return
}

// getPodiumsByPlayer 获取玩家领奖台成绩
func (c *client) getPodiumsByPlayer(playerID uint) Podiums {
	var player model.Player
	if err := c.db.Where("id = ?", playerID).First(&player).Error; err != nil {
		return Podiums{}
	}

	var out = Podiums{Player: player}

	// 查选手参加过的所有比赛且结束的
	var contestId []uint
	c.db.Model(&model.Score{}).Distinct("contest_id").Where("is_end = ?", 1).Where("player_id = ?", playerID).Pluck("player_id", &contestId)
	if len(contestId) == 0 {
		return out
	}

	// 查选手所有比赛的成绩
	for _, id := range contestId {
		topThree := c.getContestTop(id, 3)
		for _, score := range topThree {
			for idx, val := range score {
				if val.PlayerID == playerID {
					switch idx {
					case 0:
						out.Gold += 1
					case 1:
						out.Silver += 1
					case 2:
						out.Bronze += 1
					}
				}
			}
		}
	}
	return out
}

// getPodiumsByPlayer 获取比赛前top N成绩, 会依据不同项目按最佳成绩或最佳平均来区分输出
func (c *client) getContestTop(contestID uint, n int) map[model.Project][]model.Score {
	var contest model.Contest
	if err := c.db.Where("id = ? ", contestID).First(&contest).Error; err != nil || !contest.IsEnd {
		return nil
	}

	var out = make(map[model.Project][]model.Score)

	for _, project := range model.WCAProjectRoute() {
		var score []model.Score
		switch project {
		case model.Cube333MBF, model.Cube333BF, model.Cube444BF, model.Cube555BF:
			c.db.Where("contest_id = ?", contestID).Where("project = ?", project).Where("best != ?", 0).Order("best").Limit(n).Find(&score)
		default:
			c.db.Where("contest_id = ?", contestID).Where("project = ?", project).Where("avg != ?", 0).Order("avg").Limit(n).Find(&score)
		}
		if len(score) > 0 {
			out[project] = score
		}
	}
	return out
}

// getContestBestSingle 获取比赛每个项目的最佳成绩
func (c *client) getContestBestSingle(contestID uint, past bool) map[model.Project]model.Score {
	var out = make(map[model.Project]model.Score)

	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}

	for _, project := range model.WCAProjectRoute() {
		var score model.Score
		if err := c.db.Where(conn, contestID).Where("project = ?", project).Where("best != ?", 0).Order("best").First(&score).Error; err != nil {
			continue
		}
		out[project] = score
	}
	return out
}

// getContestBestAvg 获取比赛每个项目的最佳平均成绩
func (c *client) getContestBestAvg(contestID uint, past bool) map[model.Project]model.Score {
	var out = make(map[model.Project]model.Score)
	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}
	for _, project := range model.WCAProjectRoute() {
		var score model.Score
		if err := c.db.Where(conn, contestID).Where("project = ?", project).Where("avg != ?", 0).Order("avg").First(&score).Error; err != nil {
			continue
		}
		out[project] = score
	}
	return out
}

// getPodiumsByContest 获取某场比赛的领奖台
func (c *client) getPodiumsByContest(contestID uint) (out []Podiums) {
	// 未结束的比赛无领奖台
	var contest model.Contest
	if err := c.db.Where("id = ? ", contestID).First(&contest).Error; err != nil || !contest.IsEnd {
		return
	}

	// 查这场比赛所有选手
	var playerIDs []uint64
	c.db.Model(&model.Score{}).Distinct("player_id").Where("contest_id = ?", contestID).Pluck("player_id", &playerIDs)
	if len(playerIDs) == 0 {
		return
	}
	var players []model.Player
	c.db.Where("id in ?", playerIDs).Find(&players)

	var cache = make(map[uint]*Podiums)
	for _, tt := range c.getContestTop(contestID, 3) {
		for idx, val := range tt {
			if _, ok := cache[val.PlayerID]; !ok {
				cache[val.PlayerID] = &Podiums{}
			}

			switch idx {
			case 0:
				cache[val.PlayerID].Gold += 1
			case 1:
				cache[val.PlayerID].Silver += 1
			case 2:
				cache[val.PlayerID].Bronze += 1
			}
		}
	}

	for _, player := range players {
		podiums := Podiums{
			Player: player,
		}
		if val, ok := cache[player.ID]; ok {
			podiums.Gold = val.Gold
			podiums.Silver = val.Silver
			podiums.Bronze = val.Bronze
		}
		out = append(out, podiums)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Gold == out[j].Gold {
			if out[i].Silver == out[j].Silver {
				return out[i].Bronze > out[j].Bronze
			}
			return out[i].Silver > out[j].Silver
		}
		return out[i].Gold > out[j].Gold
	})
	return
}
