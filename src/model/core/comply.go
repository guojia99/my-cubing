/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午7:14.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/cache"

	"github.com/guojia99/my-cubing/src/model"
)

type client struct {
	db    *gorm.DB
	cache *cache.LRUExpireCache
}

func (c client) ReloadCache() {
	//TODO implement me
	panic("implement me")
}

func (c client) GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	//TODO implement me
	panic("implement me")
}

func (c client) GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	//TODO implement me
	panic("implement me")
}

func (c client) GetSorScore() (bestSingle, bestAvg []model.Score) {
	//TODO implement me
	panic("implement me")
}

func (c client) GetSorScoreByContest(contestID uint) (bestSingle, bestAvg []model.Score) {
	//TODO implement me
	panic("implement me")
}

func (c client) GetScoreByContest(contestID uint) map[model.Project][]RoutesScores {
	//TODO implement me
	panic("implement me")
}

func (c client) GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	//TODO implement me
	panic("implement me")
}

func (c client) GetPodiumsByPlayer(playerID uint) Podiums {
	//TODO implement me
	panic("implement me")
}

func (c client) GetPodiumsByContest(contestID uint) []Podiums {
	//TODO implement me
	panic("implement me")
}

func (c client) StatisticalRecordsAndEndContest(contestId uint) {
	//TODO implement me
	panic("implement me")
}
