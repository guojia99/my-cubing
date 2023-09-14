/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"sort"

	"github.com/guojia99/my-cubing/src/core/model"
)

type RankScore struct {
	Rank  int         `json:"Rank"` // 排名
	Score model.Score `json:"Score"`
}

type RoutesScores struct {
	Final  bool          `json:"final"`
	Round  []model.Round `json:"Round"`
	Scores []model.Score `json:"Scores"`
}

type ScoresByContest struct {
	Contest model.Contest `json:"Contest"`
	Scores  []model.Score `json:"Scores"`
	Rounds  []model.Round `json:"Rounds"`
}

type Podiums struct {
	Player         model.Player    `json:"Player"`
	Gold           int64           `json:"Gold"`
	Silver         int64           `json:"Silver"`
	Bronze         int64           `json:"Bronze"`
	PodiumsResults []PodiumsResult `json:"PodiumsResults,omitempty"`
}

type PodiumsResult struct {
	Contest model.Contest `json:"Contest"`
	Score   model.Score   `json:"Score"`
}

type SorScore struct {
	Player         model.Player `json:"Player"`
	SingleCount    int64        `json:"SingleCount,omitempty"`
	SingleProjects int64        `json:"SingleProjects,omitempty"` // 参与项目数
	AvgCount       int64        `json:"AvgCount,omitempty"`
	AvgProjects    int64        `json:"AvgProjects,omitempty"` // 参与项目数
}

func SortPodiums(in []Podiums) {
	sort.Slice(in, func(i, j int) bool {
		if in[i].Gold == in[j].Gold {
			if in[i].Silver == in[j].Silver {
				return in[i].Bronze > in[j].Bronze
			}
			return in[i].Silver > in[j].Silver
		}
		return in[i].Gold > in[j].Gold
	})
}

type RecordMessage struct {
	Record  model.Record  `json:"Record"`
	Player  model.Player  `json:"Player"`
	Score   model.Score   `json:"Score"`
	Contest model.Contest `json:"Contest"`
}

type PlayerDetail struct {
	model.Player

	ContestNumber       int `json:"ContestNumber"`
	ValidRecoveryNumber int `json:"ValidRecoveryNumber"`
	RecoveryNumber      int `json:"RecoveryNumber"`
}
