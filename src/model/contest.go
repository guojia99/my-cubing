/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:47.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

import (
	"time"

	json "github.com/json-iterator/go"
)

// Contest 比赛表，记录某场比赛
type Contest struct {
	Model

	Name        string    `json:"Name" gorm:"unique;not null;column:name"`        // 比赛名
	Content     string    `json:"contest" gorm:"column:content;null"`             // 比赛论次描述, 为 ContestRoutes 的结构体json化
	Description string    `json:"Description" gorm:"not null;column:description"` // 描述
	IsEnd       bool      `json:"IsEnd" gorm:"null;column:is_end"`                // 是否已结束
	RoundIds    string    `json:"RoundIds" gorm:"column:round_ids"`               // 轮次ID
	StartTime   time.Time `json:"StartTime" gorm:"column:start_time"`             // 开始时间
	EndTime     time.Time `json:"EndTime" gorm:"column:end_time"`                 // 结束时间
}

func (c *Contest) GetRoundIds() []uint {
	var out []uint
	_ = json.UnmarshalFromString(c.RoundIds, &out)
	return out
}

func (c *Contest) SetRoundIds(in []uint) *Contest {
	data, _ := json.MarshalToString(in)
	c.RoundIds = data
	return c
}

// Round 轮次及打乱
type Round struct {
	Model
	Project Project `json:"Project" gorm:"column:project"` // 项目
	Rank    int     `json:"Rank" gorm:"column:rank"`       // 轮次
	Name    string  `json:"Name" grom:"column:name"`       // 名
	Final   bool    `json:"Final" gorm:"column:final"`     // 是否是最后一轮
	Upsets  string  `json:"Upsets" gorm:"column:upsets"`   // 打乱 UpsetDetail
}

func (r *Round) GetUpsets() []string {
	var out []string
	_ = json.UnmarshalFromString(r.Upsets, &out)
	return out
}

func (r *Round) SetUpsets(in []string) *Round {
	data, _ := json.MarshalToString(in)
	r.Upsets = data
	return r
}
