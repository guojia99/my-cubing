/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

import (
	json "github.com/json-iterator/go"
)

// Player 选手表
type Player struct {
	Model

	Name       string `json:"Name" gorm:"unique;not null;column:name"` // 选手名
	WcaID      string `json:"WcaID,omitempty" gorm:"column:wca_id"`    // 选手WcaID，用于查询选手WCA的成绩
	ActualName string `json:"ActualName,omitempty" gorm:"actual_name"` // 真实姓名

	Titles    string   `json:"-" gorm:"titles"` // 头衔
	TitlesVal []string `json:"TitlesVal,omitempty" gorm:"-"`

	//DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}

func (c *Player) GetTitles() []string {
	var out []string
	_ = json.UnmarshalFromString(c.Titles, &out)
	c.TitlesVal = out
	return out
}

func (c *Player) SetTitles(in []string) *Player {
	c.TitlesVal = append(c.TitlesVal, in...)
	data, _ := json.MarshalToString(c.TitlesVal)
	c.Titles = data
	return c
}
