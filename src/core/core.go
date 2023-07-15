/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
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
	}
)

func NewScoreCore(db *gorm.DB) Core {
	return &client{
		db:    db,
		cache: cache.NewLRUExpireCache(255),
	}
}

type client struct {
	db    *gorm.DB
	cache *cache.LRUExpireCache
}

func (c *client) ReloadCache() {
	for _, key := range c.cache.Keys() {
		c.cache.Remove(key)
	}
}

func (c *client) AddScore(playerName string, contestID uint, project model.Project, routeNum int, result []float64) error {
	return c.addScore(playerName, contestID, project, routeNum, result)
}

func (c *client) RemoveScore(playerName string, contestID uint, project model.Project, routeNum int) error {
	return c.removeScoreByContestID(playerName, contestID, project, routeNum)
}

func (c *client) StatisticalRecordsAndEndContest(contestId uint) error {
	defer c.ReloadCache()
	return c.statisticalRecordsAndEndContest(contestId)
}

func (c *client) GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetSorScore() (single, avg []SorScore) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetSorScoreByContest(contestID uint) (single, avg []SorScore) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetScoreByContest(contestID uint) map[model.Project][]RoutesScores {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetPodiumsByPlayer(playerID uint) Podiums {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetPodiumsByContest(contestID uint) []Podiums {
	//TODO implement me
	panic("implement me")
}
