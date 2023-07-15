/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import "github.com/guojia99/my-cubing/src/core/model"

type RoutesScores struct {
	Round  model.Round   `json:"Round"`
	Scores []model.Score `json:"Scores"`
}

type ScoresByContest struct {
	Contest model.Contest `json:"Contest"`
	Scores  []model.Score `json:"Scores"`
}

type Podiums struct {
	Player model.Player `json:"Player"`
	Gold   int64        `json:"Gold"`
	Silver int64        `json:"Silver"`
	Bronze int64        `json:"Bronze"`
}

type SorScore struct {
	Player      model.Player `json:"Player"`
	SingleCount int64        `json:"SingleCount"`
	AvgCount    int64        `json:"Count"`
}
