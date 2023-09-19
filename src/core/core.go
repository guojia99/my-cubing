/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"fmt"
	"runtime"
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing/src/core/model"
)

func NewScoreCore(db *gorm.DB, debug bool) Core {
	return &client{
		debug: debug,
		db:    db,
		cache: cache.New(time.Minute*60, time.Minute*60),
	}
}

type client struct {
	debug bool

	db    *gorm.DB
	cache *cache.Cache
}

func (c *client) ReloadCache() {
	c.cache.Flush()
	runtime.GC()
}

func (c *client) AddScore(playerID uint, contestID uint, project model.Project, roundId uint, result []float64, penalty model.ScorePenalty) error {
	if err := c.addScore(playerID, contestID, project, roundId, result, penalty); err != nil {
		return err
	}
	c.ReloadCache()
	return nil
}

func (c *client) RemoveScore(scoreID uint) error {
	if err := c.removeScore(scoreID); err != nil {
		return err
	}
	c.ReloadCache()
	return nil
}

func (c *client) StatisticalRecordsAndEndContest(contestId uint) error {
	if err := c.statisticalRecordsAndEndContest(contestId); err != nil {
		return err
	}
	c.ReloadCache()
	return nil
}

func (c *client) GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	key := "GetBestScores"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.Project]model.Score)
		return result[0], result[1]
	}

	bestSingle, bestAvg = c.getBestScores()
	_ = c.cache.Add(key, [2]map[model.Project]model.Score{bestSingle, bestAvg}, time.Minute*60)
	return
}

func (c *client) GetPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore) {
	key := fmt.Sprintf("GetPlayerBestScore_%d", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.Project]RankScore)
		return result[0], result[1]
	}

	bestSingle, bestAvg = make(map[model.Project]RankScore), make(map[model.Project]RankScore)
	allBest, allAvg := c.GetAllPlayerBestScore()
	for _, val := range allBest {
		for i := 0; i < len(val); i++ {
			if val[i].PlayerID == playerId {
				bestSingle[val[i].Project] = RankScore{
					Rank:  i + 1,
					Score: val[i],
				}
				break
			}
		}
	}
	for _, val := range allAvg {
		for i := 0; i < len(val); i++ {
			if val[i].PlayerID == playerId {
				bestAvg[val[i].Project] = RankScore{
					Rank:  i + 1,
					Score: val[i],
				}
				break
			}
		}
	}
	_ = c.cache.Add(key, [2]map[model.Project]RankScore{bestSingle, bestAvg}, time.Minute*60)
	return bestSingle, bestAvg
}

func (c *client) GetPlayerDetail(playerId uint) PlayerDetail {
	key := fmt.Sprintf("GetPlayerDetail%d", playerId)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(PlayerDetail)
	}

	out := c.getPlayerDetail(playerId)
	_ = c.cache.Add(key, out, time.Minute*60)
	return out
}

func (c *client) GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	key := "GetAllPlayerBestScore"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.Project][]model.Score)
		return result[0], result[1]
	}

	bestSingle, bestAvg = c.getAllPlayerBestScore()
	_ = c.cache.Add(key, [2]map[model.Project][]model.Score{bestSingle, bestAvg}, time.Minute*60)
	return
}

func (c *client) GetAllPlayerBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) {
	best, avg := c.GetAllPlayerBestScore()
	return best[project], avg[project]
}

func (c *client) GetSorScore() (single, avg map[model.SorStatisticsKey][]SorScore) {
	key := "GetSorScore"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.SorStatisticsKey][]SorScore)
		return result[0], result[1]
	}

	single, avg = c.getSorScore()
	_ = c.cache.Add(key, [2]map[model.SorStatisticsKey][]SorScore{single, avg}, time.Minute*60)
	return
}

func (c *client) GetSorScoreByContest(contestID uint) (single, avg map[model.SorStatisticsKey][]SorScore) {
	key := fmt.Sprintf("GetSorScoreByContest%d", contestID)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.SorStatisticsKey][]SorScore)
		return result[0], result[1]
	}

	single, avg = c.getSorScoreByContest(contestID)
	_ = c.cache.Add(key, [2]map[model.SorStatisticsKey][]SorScore{single, avg}, time.Minute*60)
	return
}

func (c *client) GetPlayerSor(playerID uint) (single, avg map[model.SorStatisticsKey]SorScore) {
	key := fmt.Sprintf("GetPlayerSor_%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.SorStatisticsKey]SorScore)
		return result[0], result[1]
	}

	single, avg = make(map[model.SorStatisticsKey]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey]SorScore, len(model.SorKeyMap()))
	singleCache, avgCache := c.GetSorScore()

	for k, _ := range model.SorKeyMap() {
		if _, ok := singleCache[k]; ok {
			for idx, score := range singleCache[k] {
				if score.Player.ID == playerID {
					score.SingleRank = int64(idx + 1)
					single[k] = score
					break
				}
			}
		}
		if _, ok := avgCache[k]; ok {
			for idx, score := range avgCache[k] {
				if score.Player.ID == playerID {
					score.AvgRank = int64(idx + 1)
					avg[k] = score
					break
				}
			}
		}
	}
	_ = c.cache.Add(key, [2]map[model.SorStatisticsKey]SorScore{single, avg}, time.Minute*60)
	return single, avg
}

func (c *client) GetScoreByContest(contestID uint) map[model.Project][]RoutesScores {
	key := fmt.Sprintf("GetScoreByContest%d", contestID)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(map[model.Project][]RoutesScores)
	}

	out := c.getScoreByContest(contestID)
	if out == nil {
		return out
	}
	_ = c.cache.Add(key, out, time.Minute*60)
	return out
}

func (c *client) GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	key := fmt.Sprintf("GetPlayerScore%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([]interface{})
		return result[0].([]model.Score), result[1].([]model.Score), result[2].([]ScoresByContest)
	}

	bestSingle, bestAvg, scores = c.getPlayerScore(playerID)
	_ = c.cache.Add(key, []interface{}{bestSingle, bestAvg, scores}, time.Minute*60)
	return
}

func (c *client) GetPodiumsByPlayer(playerID uint) Podiums {
	key := fmt.Sprintf("GetPodiumsByPlayer%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(Podiums)
	}

	out := c.getPodiumsByPlayer(playerID)
	_ = c.cache.Add(key, out, time.Minute*60)
	return out
}

func (c *client) GetPodiumsByContest(contestID uint) []Podiums {
	key := fmt.Sprintf("GetPodiumsByContest%d", contestID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out := c.getPodiumsByContest(contestID)
	_ = c.cache.Add(key, out, time.Minute*60)
	return out
}

func (c *client) GetAllPodium() []Podiums {
	key := "GetAllPodium"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}
	out := c.getAllPodium()
	_ = c.cache.Add(key, out, time.Minute*30)
	return out
}

func (c *client) GetRecordByContest(contestID uint) []RecordMessage {
	key := fmt.Sprintf("GetRecordByContest_%d", contestID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]RecordMessage)
	}
	out := c.getRecordByContest(contestID)
	_ = c.cache.Add(key, out, time.Minute*5)
	return out
}

func (c *client) GetRecordByPlayer(playerID uint) []RecordMessage {
	key := fmt.Sprintf("GetRecordByPlayer%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]RecordMessage)
	}
	out := c.getRecordByPlayer(playerID)
	_ = c.cache.Add(key, out, time.Minute*5)
	return out
}
