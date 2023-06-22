/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:45.
 * Author:  guojia(https://github.com/guojia99)
 */

package db

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primaryKey;column:id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;deleted_at"`
}

// Contest 比赛表，记录某场比赛
type Contest struct {
	Model

	Name        string    `json:"Name" gorm:"uniqueIndex;not null;column:name"`   // 比赛名
	Description string    `json:"Description" gorm:"not null;column:description"` // 描述
	StartTime   time.Time `json:"StartTime" gorm:"column:start_time"`             // 开始时间
	EndTime     time.Time `json:"EndTime" gorm:"column:end_time"`                 // 结束时间
}

// Player 选手表
type Player struct {
	Model

	Name  string `json:"Name" gorm:"uniqueIndex;not null;column:name"` // 选手名
	WcaID string `json:"WcaID" gorm:"column:wca_id"`                   // 选手WcaID，用于查询选手WCA的成绩
}

// Score 成绩表
type Score struct {
	Model

	ContestID   uint    `json:"ContestID" gorm:"index;not null;column:contest_id"` // 比赛的ID
	PlayerID    uint    `json:"PlayerID" gorm:"index;not null;column:player_id"`   // 选手的ID
	RouteNumber uint    `json:"RouteNumber" gorm:"not null;column:route_number"`   // 该项目的轮次
	Project     string  `json:"Project" gorm:"not null;column:project"`            // 分项目 333/222/444等
	Value       float64 `json:"Score" gorm:"gorm;column:value"`                    // 成绩
}
