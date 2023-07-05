/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/1 下午3:40.
 *  * Author: guojia(https://github.com/guojia99)
 */

package score

import "my-cubing/db"

type GetUserContestScoreRequest struct {
	PlayerName string `uri:"player_name"`
	ContestID  uint   `uri:"contest_id"`
}

type CreateScoreRequest struct {
	PlayerName  string    `json:"PlayerName"`
	ContestID   uint      `json:"ContestID"`
	RouteNumber uint      `json:"RouteNumber"`
	ProjectName string    `json:"ProjectName"`
	Results     []float64 `json:"Results"`
}

type (
	BestScore struct {
		Project    string  `json:"Project"`
		BestPlayer string  `json:"BestPlayer"`
		BestScore  float64 `json:"BestScore"`
		AvgPlayer  string  `json:"AvgPlayer"`
		AvgScore   float64 `json:"AvgScore"`
	}

	GetAllProjectBestScoreResponse struct {
		Data []BestScore `json:"Data"`
	}
)

type (
	ProjectScores struct {
		Player string `json:"Player"`
		db.Score
	}

	GetProjectScoresResponse struct {
		ProjectList []string                   `json:"ProjectList"`
		Avg         map[string][]ProjectScores `json:"Avg"`
		Best        map[string][]ProjectScores `json:"Best"`
	}
)

type (
	SorScoreDetail struct {
		db.Score
		Player string `json:"Player"`
		Count  int    `json:"Count"`
	}

	GetSorScoresResponse struct {
		Best []SorScoreDetail `json:"Best"`
		Avg  []SorScoreDetail `json:"Avg"`
	}
)

type (
	GetContestScoresRequest struct {
		ContestID uint `uri:"contest_id"`
	}

	GetContestScoresDetail struct {
		Player string `json:"Player"`
		db.Score
	}

	GetContestScoresPlayer struct {
		Name string `json:"Name"`
		Id   uint   `json:"Id"`
	}

	GetContestScoresResponse struct {
		ContestName string                              `json:"ContestName"`
		Players     []GetContestScoresPlayer            `json:"Players"`
		ProjectList []string                            `json:"ProjectList"`
		Data        map[string][]GetContestScoresDetail `json:"Data"`
	}
)

type (
	EndContestScoreRequest struct {
		ContestID uint `uri:"contest_id"`
	}
)
