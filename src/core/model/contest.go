/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
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

	Name        string    `json:"Name,omitempty" gorm:"unique;not null;column:name"`        // 比赛名
	Type        string    `json:"Type,omitempty" gorm:"column:c_type"`                      // 类型 正式 | 线上 | 线下
	Description string    `json:"Description,omitempty" gorm:"not null;column:description"` // 描述
	IsEnd       bool      `json:"IsEnd,omitempty" gorm:"null;column:is_end"`                // 是否已结束
	RoundIds    string    `json:"RoundIds,omitempty" gorm:"column:round_ids"`               // 轮次ID
	RoundIdsVal []uint    `json:"RoundIdsVal,omitempty" gorm:"-"`                           // 轮次ID实际内容
	StartTime   time.Time `json:"StartTime,omitempty" gorm:"column:start_time"`             // 开始时间
	EndTime     time.Time `json:"EndTime,omitempty" gorm:"column:end_time"`                 // 结束时间
}

func (c *Contest) GetRoundIds() []uint {
	var out []uint
	_ = json.UnmarshalFromString(c.RoundIds, &out)
	c.RoundIdsVal = out
	return out
}

func (c *Contest) SetRoundIds(in []uint) *Contest {
	data, _ := json.MarshalToString(in)
	c.RoundIds = data
	c.RoundIdsVal = in
	return c
}

// Round 轮次及打乱
type Round struct {
	Model
	Name      string   `json:"Name,omitempty" gorm:"column:name"`
	ContestID uint     `json:"ContestID,omitempty" gorm:"column:contest_id"` // 所属比赛
	Project   Project  `json:"Project,omitempty" gorm:"column:project"`      // 项目
	Number    int      `json:"Number,omitempty" gorm:"column:number"`        // 项目轮次
	Part      int      `json:"Part,omitempty" gorm:"column:part"`            // 该轮次第几份打乱
	Final     bool     `json:"Final,omitempty" gorm:"column:final"`          // 是否是最后一轮
	IsStart   bool     `json:"IsStart,omitempty" gorm:"column:is_start"`     // 是否已开始
	Upsets    string   `json:"-" gorm:"column:upsets"`                       // 打乱 UpsetDetail
	UpsetsVal []string `json:"UpsetsVal,omitempty" gorm:"-"`                 // 打乱 UpsetDetail 实际内容
}

func (r *Round) GetUpsets() []string {
	var out []string
	_ = json.UnmarshalFromString(r.Upsets, &out)
	r.UpsetsVal = out
	return out
}

func (r *Round) SetUpsets(in []string) *Round {
	data, _ := json.MarshalToString(in)
	r.Upsets = data
	r.UpsetsVal = in
	return r
}
