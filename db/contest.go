/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:47.
 *  * Author: guojia(https://github.com/guojia99)
 */

package db

import "time"

type (
	Route struct {
		Project      Project  `json:"Project"`      // 项目
		RoutePeoples []int    `json:"RoutePeoples"` // 每轮人数
		Players      [][]uint `json:"Players"`      // 每轮比赛的参赛者记录
	}

	ContestContent struct {
		Routes []Route `json:"Routes"` // 每个项目
	}
)

// Contest 比赛表，记录某场比赛
type Contest struct {
	Model

	Name        string    `json:"Name" gorm:"unique;not null;column:name"`        // 比赛名
	Content     string    `json:"contest" gorm:"column:content;null"`             // 比赛论次描述, 为 ContestRoutes 的结构体json化
	Description string    `json:"Description" gorm:"not null;column:description"` // 描述
	IsEnd       bool      `json:"IsEnd" gorm:"null;column:is_end"`                // 是否已结束
	StartTime   time.Time `json:"StartTime" gorm:"column:start_time"`             // 开始时间
	EndTime     time.Time `json:"EndTime" gorm:"column:end_time"`                 // 结束时间
}
