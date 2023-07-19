/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/cache"

	"github.com/guojia99/my-cubing/src/core/model"
)

type (
	Core interface {
		Read
		// ReloadCache 重置缓存
		ReloadCache()
		// AddScore 添加成绩
		AddScore(playerName string, contestID uint, project model.Project, routeNum int, result []float64) error
		// RemoveScore 删除成绩
		RemoveScore(playerName string, contestID uint, project model.Project, routeNum int) error
		// StatisticalRecordsAndEndContest 结束比赛并统计记录
		StatisticalRecordsAndEndContest(contestId uint) error
	}

	Read interface {
		// GetBestScores 获取所有项目最佳成绩
		GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score)
		// GetAllPlayerBestScore 获取所有人最佳成绩
		GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score)
		// GetAllPlayerBestScoreByProject 获取某个项目最佳成绩
		GetAllPlayerBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score)
		// GetSorScore 获取排名总和
		GetSorScore() (single, avg []SorScore)
		// GetSorScoreByContest 获取某场比赛的排名总和
		GetSorScoreByContest(contestID uint) (single, avg []SorScore)
		// GetScoreByContest 获取某场比赛成绩排名
		GetScoreByContest(contestID uint) map[model.Project][]RoutesScores
		// GetPlayerScore 获取选手所有成绩
		GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest)
		// GetPodiumsByPlayer 获取玩家领奖台数据
		GetPodiumsByPlayer(playerID uint) Podiums
		// GetPodiumsByContest 获取比赛的领奖台数据
		GetPodiumsByContest(contestID uint) []Podiums
		// GetAllPodium 获取全部人的领奖台排行
		GetAllPodium() []Podiums
		// GetRecordByContest 获取一场比赛中的记录
		GetRecordByContest(contestID uint) []model.Record
		// GetRecordByPlayer 获取一个人的记录
		GetRecordByPlayer(playerID uint) []model.Record
	}
)

func NewScoreCore(db *gorm.DB, debug bool) Core {
	return &client{
		debug: debug,
		db:    db,
		cache: cache.NewLRUExpireCache(255),
	}
}

type client struct {
	debug bool

	db    *gorm.DB
	cache *cache.LRUExpireCache
}

func (c *client) ReloadCache() {
	for _, key := range c.cache.Keys() {
		c.cache.Remove(key)
	}
}

func (c *client) AddScore(playerName string, contestID uint, project model.Project, routeNum int, result []float64) error {
	if err := c.addScore(playerName, contestID, project, routeNum, result); err != nil {
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
	c.cache.Add(key, [2]map[model.Project]model.Score{bestSingle, bestAvg}, time.Minute*15)
	return
}

func (c *client) GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	key := "GetAllPlayerBestScore"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]map[model.Project][]model.Score)
		return result[0], result[1]
	}

	bestSingle, bestAvg = c.getAllPlayerBestScore()
	c.cache.Add(key, [2]map[model.Project][]model.Score{bestSingle, bestAvg}, time.Minute*15)
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
	c.cache.Add(key, [2][]SorScore{single, avg}, time.Minute*15)
	return
}

func (c *client) GetSorScoreByContest(contestID uint) (single, avg []SorScore) {
	key := fmt.Sprintf("GetSorScoreByContest%d", contestID)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2][]SorScore)
		return result[0], result[1]
	}

	single, avg = c.getSorScoreByContest(contestID)
	c.cache.Add(key, [2][]SorScore{single, avg}, time.Minute*15)
	return
}

func (c *client) GetScoreByContest(contestID uint) map[model.Project][]RoutesScores {
	key := fmt.Sprintf("GetScoreByContest%d", contestID)

	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(map[model.Project][]RoutesScores)
	}

	out := c.getScoreByContest(contestID)
	c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	key := fmt.Sprintf("GetPlayerScore%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([]interface{})
		return result[0].([]model.Score), result[1].([]model.Score), result[2].([]ScoresByContest)
	}

	bestSingle, bestAvg, scores = c.getPlayerScore(playerID)
	c.cache.Add(key, []interface{}{bestSingle, bestAvg, scores}, time.Minute*15)
	return
}

func (c *client) GetPodiumsByPlayer(playerID uint) Podiums {
	key := fmt.Sprintf("GetPodiumsByPlayer%d", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(Podiums)
	}

	out := c.getPodiumsByPlayer(playerID)
	c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetPodiumsByContest(contestID uint) []Podiums {
	key := fmt.Sprintf("GetPodiumsByContest%d", contestID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out := c.getPodiumsByContest(contestID)
	c.cache.Add(key, out, time.Minute*15)
	return out
}

func (c *client) GetAllPodium() []Podiums {
	key := "GetAllPodium"
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}
	out := c.getAllPodium()
	c.cache.Add(key, out, time.Minute*5)
	return out
}

func (c *client) GetRecordByContest(contestID uint) []model.Record {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetRecordByPlayer(playerID uint) []model.Record {
	//TODO implement me
	panic("implement me")
}
