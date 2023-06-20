/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/20 下午1:57.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

import (
	"time"

	"gorm.io/gorm"
)

// Contest 比赛表，记录某场比赛
type Contest struct {
	gorm.Model
	Name        string    `json:"Name"`        // 比赛名
	Description string    `json:"Description"` // 描述
	StartTime   time.Time `json:"StartTime"`   // 开始时间
	EndTime     time.Time `json:"EndTime"`     // 结束时间
}

// Score 成绩表
type Score struct {
	gorm.Model
	Score       float64 `json:"Score"`       // 成绩
	ContestID   int     `json:"ContestID"`   // 比赛的ID
	Project     string  `json:"Project"`     // 分项目 333/222/444等
	RouteNumber int     `json:"RouteNumber"` // 该项目的轮次
}

// Player 选手表
type Player struct {
	gorm.Model
	Name  string // 选手名
	WcaID string // 选手WcaID，用于查询选手WCA的成绩
}
