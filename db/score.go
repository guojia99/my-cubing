/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:47.
 *  * Author: guojia(https://github.com/guojia99)
 */

package db

import "time"

// Score 成绩表
type Score struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	// 这里的设计是
	// 一次复原成绩记录一条
	PlayerID    uint `json:"PlayerID" gorm:"index;not null;column:player_id"`   // 选手的ID
	ContestID   uint `json:"ContestID" gorm:"index;not null;column:contest_id"` // 比赛的ID
	RouteNumber uint `json:"RouteNumber" gorm:"not null;column:route_number"`   // 该项目的轮次

	Project Project `json:"Project" gorm:"not null;column:project"` // 分项目 333/222/444等
	Result1 float64 `json:"R1" gorm:"column:r1;NULL"`               // 成绩1
	Result2 float64 `json:"R2" gorm:"column:r2;NULL"`               // 成绩2
	Result3 float64 `json:"R3" gorm:"column:r3;NULL"`               // 成绩3
	Result4 float64 `json:"R4" gorm:"column:r4;NULL"`               // 成绩4
	Result5 float64 `json:"R5" gorm:"column:r5;NULL"`               // 成绩5
	Best    float64 `json:"Best" gorm:"column:best;NULL"`           // 五把最好成绩
	Avg     float64 `json:"Avg" gorm:"column:avg;NULL"`             // 五把平均成绩
}

func (s Score) GetResult() []float64 {
	return []float64{s.Result1, s.Result2, s.Result3, s.Result4, s.Result5}
}

func (s Score) GetBest() float64 {
	return 0
}

func (s Score) GetAvg() float64 {
	return 0
}
