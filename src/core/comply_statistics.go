package core

import (
	"fmt"
	"sort"

	"github.com/guojia99/my-cubing/src/core/model"
)

func ParserSorSort(players []model.Player, bestSingle, bestAvg map[model.Project][]model.Score) (single, avg map[model.SorStatisticsKey][]SorScore) {
	single, avg = make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap()))

	var singlePlayerDict = make(map[model.Project]map[uint]model.Score)
	var avgPlayerDict = make(map[model.Project]map[uint]model.Score)

	// 1. 做map缓存
	for _, pj := range model.AllProjectRoute() {
		singlePlayerDict[pj] = make(map[uint]model.Score)
		avgPlayerDict[pj] = make(map[uint]model.Score)

		if _, ok := bestSingle[pj]; ok {
			for _, val := range bestSingle[pj] {
				singlePlayerDict[pj][val.PlayerID] = val
			}
		}
		if _, ok := bestAvg[pj]; ok {
			for _, val := range bestAvg[pj] {
				avgPlayerDict[pj][val.PlayerID] = val
			}
		}
	}

	// 2. 全项目排序
	for sorKey, projects := range model.SorKeyMap() {
		var s = make([]SorScore, 0)
		var a = make([]SorScore, 0)
		var playerCache = make(map[uint]*SorScore)

		for _, player := range players {
			playerCache[player.ID] = &SorScore{Player: player}

			for _, pj := range projects {
				if val, ok := singlePlayerDict[pj][player.ID]; ok {
					playerCache[val.PlayerID].SingleCount += int64(val.Rank)
					playerCache[val.PlayerID].SingleProjects += 1
				} else if len(bestSingle[pj]) > 0 {
					playerCache[player.ID].SingleCount += int64(bestSingle[pj][len(bestSingle[pj])-1].Rank + 1)
				}

				if val, ok := avgPlayerDict[pj][player.ID]; ok {
					playerCache[val.PlayerID].AvgCount += int64(val.Rank)
					playerCache[val.PlayerID].AvgProjects += 1
				} else if len(bestAvg[pj]) > 0 {
					playerCache[player.ID].AvgCount += int64(bestAvg[pj][len(bestAvg[pj])-1].Rank + 1)
				}
			}
		}

		for _, val := range playerCache {
			s = append(s, SorScore{Player: val.Player, SingleCount: val.SingleCount, SingleProjects: val.SingleProjects})
			a = append(a, SorScore{Player: val.Player, AvgCount: val.AvgCount, AvgProjects: val.AvgProjects})
		}

		sort.Slice(s, func(i, j int) bool {
			if s[i].SingleCount == s[j].SingleCount {
				return s[i].SingleProjects < s[j].SingleProjects
			}
			return s[i].SingleCount < s[j].SingleCount
		})
		sort.Slice(a, func(i, j int) bool {
			if a[i].AvgCount == a[j].AvgCount {
				return a[i].AvgProjects < a[j].AvgProjects
			}
			return a[i].AvgCount < a[j].AvgCount
		})

		single[sorKey] = s
		avg[sorKey] = a
	}
	return
}

// getSorScore 获取所有玩家的Sor排名
func (c *client) getSorScore() (single, avg map[model.SorStatisticsKey][]SorScore) {

	var players []model.Player
	if err := c.db.Find(&players).Error; err != nil {
		return
	}
	bestSingle, bestAvg := c.getAllPlayerBestScore()
	single, avg = ParserSorSort(players, bestSingle, bestAvg)

	return
}

// getSorScoreByContest 获取比赛的Sor排名
func (c *client) getSorScoreByContest(contestID uint) (single, avg map[model.SorStatisticsKey][]SorScore) {
	single, avg = make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap()))

	// 查这场比赛所有选手
	var playerIDs []uint64
	if c.db.Model(&model.Score{}).Distinct("player_id").Where("contest_id = ?", contestID).Pluck("player_id", &playerIDs); len(playerIDs) == 0 {
		return
	}
	var players []model.Player
	c.db.Where("id in ?", playerIDs).Find(&players)

	bestSingleCache, bestAvgCache := c.getContestAllBestScores(contestID)
	fmt.Println(bestAvgCache[model.Cube333BF])
	single, avg = ParserSorSort(players, bestSingleCache, bestAvgCache)
	return
}
