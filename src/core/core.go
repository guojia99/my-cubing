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

type (
	Core interface {
		Read
		ReadPlayer
		ReadContest
		// ReloadCache 重置缓存
		ReloadCache()
		// AddScore 添加成绩
		AddScore(playerName string, contestID uint, project model.Project, routeNum int, result []float64, penalty model.ScorePenalty) error
		// RemoveScore 删除成绩
		RemoveScore(playerName string, contestID uint, project model.Project, routeNum int) error
		// StatisticalRecordsAndEndContest 结束比赛并统计记录
		StatisticalRecordsAndEndContest(contestId uint) error
	}

	ReadPlayer interface {
		// GetPlayerDetail 获取玩家参赛信息
		GetPlayerDetail(playerId uint) PlayerDetail
		// GetPlayerBestScore 获取玩家最佳成绩及排名
		GetPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore)
		// GetAllPlayerBestScore 获取所有人最佳成绩
		GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score)
		// GetAllPlayerBestScoreByProject 获取某个项目最佳成绩
		GetAllPlayerBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score)
		// GetPodiumsByPlayer 获取玩家领奖台数据
		GetPodiumsByPlayer(playerID uint) Podiums
		// GetRecordByPlayer 获取一个人的记录
		GetRecordByPlayer(playerID uint) []RecordMessage
		// GetPlayerScore 获取选手所有成绩
		GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest)
	}

	ReadContest interface {
		// GetSorScoreByContest 获取某场比赛的排名总和
		GetSorScoreByContest(contestID uint) (single, avg []SorScore)
		// GetScoreByContest 获取某场比赛成绩排名
		GetScoreByContest(contestID uint) map[model.Project][]RoutesScores
		// GetPodiumsByContest 获取比赛的领奖台数据
		GetPodiumsByContest(contestID uint) []Podiums
	}

	Read interface {
		// GetBestScores 获取所有项目最佳成绩
		GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score)
		// GetSorScore 获取排名总和
		GetSorScore() (single, avg []SorScore)
		// GetAllPodium 获取全部人的领奖台排行
		GetAllPodium() []Podiums
		// GetRecordByContest 获取一场比赛中的记录
		GetRecordByContest(contestID uint) []RecordMessage
	}
)

func NewScoreCore(db *gorm.DB, debug bool) Core {
	return &client{
		debug: debug,
		db:    db,
		cache: cache.New(time.Minute*15, time.Minute*15),
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

func (c *client) AddScore(playerName string, contestID uint, project model.Project, routeNum int, result []float64, penalty model.ScorePenalty) error {
	if err := c.addScore(playerName, contestID, project, routeNum, result, penalty); err != nil {
		return err
	}
	c.ReloadCache()
	return nil
}

func (c *client) RemoveScore(playerName string, contestID uint, project model.Project, routeNum int) error {
	if err := c.removeScoreByContestID(playerName, contestID, project, routeNum); err != nil {
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
	_ = c.cache.Add(key, [2]map[model.Project]model.Score{bestSingle, bestAvg}, time.Minute*15)
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
	_ = c.cache.Add(key, [2]map[model.Project]RankScore{bestSingle, bestAvg}, time.Minute*15)
	return bestSingle, bestAvg
}

func (c *client) GetPlayerDetail(playerId uint) PlayerDetail {
	key := fmt.Sprintf("GetPlayerDetail%d", playerId)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(PlayerDetail)
	}

	out := c.getPlayerDetail(playerId)
	_ = c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	key := "GetAllPlayerBestScore"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.Project][]model.Score)
		return result[0], result[1]
	}

	bestSingle, bestAvg = c.getAllPlayerBestScore()
	_ = c.cache.Add(key, [2]map[model.Project][]model.Score{bestSingle, bestAvg}, time.Minute*15)
	return
}

func (c *client) GetAllPlayerBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) {
	best, avg := c.GetAllPlayerBestScore()
	return best[project], avg[project]
}

func (c *client) GetSorScore() (single, avg []SorScore) {
	key := "GetSorScore"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2][]SorScore)
		return result[0], result[1]
	}

	single, avg = c.getSorScore()
	_ = c.cache.Add(key, [2][]SorScore{single, avg}, time.Minute*15)
	return
}

func (c *client) GetSorScoreByContest(contestID uint) (single, avg []SorScore) {
	key := fmt.Sprintf("GetSorScoreByContest%d", contestID)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2][]SorScore)
		return result[0], result[1]
	}

	single, avg = c.getSorScoreByContest(contestID)
	_ = c.cache.Add(key, [2][]SorScore{single, avg}, time.Minute*15)
	return
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
	_ = c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	key := fmt.Sprintf("GetPlayerScore%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([]interface{})
		return result[0].([]model.Score), result[1].([]model.Score), result[2].([]ScoresByContest)
	}

	bestSingle, bestAvg, scores = c.getPlayerScore(playerID)
	_ = c.cache.Add(key, []interface{}{bestSingle, bestAvg, scores}, time.Minute*15)
	return
}

func (c *client) GetPodiumsByPlayer(playerID uint) Podiums {
	key := fmt.Sprintf("GetPodiumsByPlayer%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(Podiums)
	}

	out := c.getPodiumsByPlayer(playerID)
	_ = c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetPodiumsByContest(contestID uint) []Podiums {
	key := fmt.Sprintf("GetPodiumsByContest%d", contestID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out := c.getPodiumsByContest(contestID)
	_ = c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetAllPodium() []Podiums {
	key := "GetAllPodium"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}
	out := c.getAllPodium()
	_ = c.cache.Add(key, out, time.Minute*5)
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
