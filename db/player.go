/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:47.
 *  * Author: guojia(https://github.com/guojia99)
 */

package db

// Player 选手表
type Player struct {
	Model

	Name  string `json:"Name" gorm:"unique;not null;column:name"` // 选手名
	WcaID string `json:"WcaID" gorm:"column:wca_id"`              // 选手WcaID，用于查询选手WCA的成绩
}
