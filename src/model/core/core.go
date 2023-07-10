/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午6:07.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/cache"

	"github.com/guojia99/my-cubing/src/model"
)

type Core interface {
	// ReloadCache 重置缓存
	ReloadCache()

	// GetBestScores 获取所有项目最佳成绩
	GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score)
	// GetAllPlayerBestScore 获取所有人最佳成绩
	GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score)
	// GetSorScore 获取排名总和
	GetSorScore() (bestSingle, bestAvg []model.Score)
	// GetSorScoreByContest 获取某场比赛的排名总和
	GetSorScoreByContest(contestID uint) (bestSingle, bestAvg []model.Score)
	// GetScoreByContest 获取某场比赛成绩排名
	GetScoreByContest(contestID uint) map[model.Project][]RoutesScores
	// GetPlayerScore 获取选手所有成绩
	GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest)
	// GetPodiumsByPlayer 获取玩家领奖台数据
	GetPodiumsByPlayer(playerID uint) Podiums
	// GetPodiumsByContest 获取比赛的领奖台数据
	GetPodiumsByContest(contestID uint) []Podiums

	// StatisticalRecordsAndEndContest 结束比赛并统计记录
	StatisticalRecordsAndEndContest(contestId uint)
}

func NewScoreCore(db *gorm.DB) Core {
	return &client{
		db:    db,
		cache: cache.NewLRUExpireCache(255),
	}
}
