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
	Rank  int         `json:"Rank,omitempty"` // 排名
	Score model.Score `json:"Score,omitempty"`
}

type RoutesScores struct {
	Final  bool          `json:"final,omitempty"`
	Round  []model.Round `json:"Round,omitempty"`
	Scores []model.Score `json:"Scores,omitempty"`
}

type ScoresByContest struct {
	Contest model.Contest `json:"Contest,omitempty"`
	Scores  []model.Score `json:"Scores,omitempty"`
	Rounds  []model.Round `json:"Rounds,omitempty"`
}

type Podiums struct {
	Player         model.Player    `json:"Player,omitempty"`
	Gold           int64           `json:"Gold,omitempty"`
	Silver         int64           `json:"Silver,omitempty"`
	Bronze         int64           `json:"Bronze,omitempty"`
	PodiumsResults []PodiumsResult `json:"PodiumsResults,omitempty"`
}

type PodiumsResult struct {
	Contest model.Contest `json:"Contest,omitempty"`
	Score   model.Score   `json:"Score,omitempty"`
}

type SorScore struct {
	Player         model.Player `json:"Player,omitempty"`
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
	Record  model.Record  `json:"Record,omitempty"`
	Player  model.Player  `json:"Player,omitempty"`
	Score   model.Score   `json:"Score,omitempty"`
	Contest model.Contest `json:"Contest,omitempty"`
}

type PlayerDetail struct {
	model.Player

	ContestNumber       int `json:"ContestNumber,omitempty"`
	ValidRecoveryNumber int `json:"ValidRecoveryNumber,omitempty"`
	RecoveryNumber      int `json:"RecoveryNumber,omitempty"`
}
