/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

const (
	RecordByAvg    = 1
	RecordBySingle = 2
)

type Record struct {
	Model

	RType        uint    `json:"RType,omitempty" gorm:"column:rtype"`       // 记录类型
	ScoreId      uint    `json:"score_id,omitempty" gorm:"column:score_id"` // 成绩记录位置
	ScoreValue   Score   `json:"ScoreValue,omitempty" gorm:"-"`
	PlayerID     uint    `json:"PlayerID,omitempty" gorm:"index;not null;column:player_id"`   // 选手的ID
	PlayerName   string  `json:"PlayerName,omitempty" gorm:"column:player_name"`              // 玩家名
	ContestID    uint    `json:"ContestID,omitempty" gorm:"index;not null;column:contest_id"` // 比赛的ID
	ContestValue Contest `json:"ContestValue,omitempty" gorm:"-"`
}
