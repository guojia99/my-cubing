/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午7:14.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"sort"

	"github.com/guojia99/my-cubing/src/model"
)

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

func (c *client) getScoreByContest(contestID uint) map[model.Project][]RoutesScores {
	var scores []model.Score
	c.db.Where("contest_id = ?", contestID).Find(&scores)
}

func (c *client) getSorScoreByContest(contestID uint) (single, avg []SorScore) {
	//TODO implement me
	panic("implement me")
}

func (c *client) getPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	//TODO implement me
	panic("implement me")
}

func (c *client) getPodiumsByPlayer(playerID uint) Podiums {
	//TODO implement me
	panic("implement me")
}

func (c *client) getPodiumsByContest(contestID uint) []Podiums {
	//TODO implement me
	panic("implement me")
}
