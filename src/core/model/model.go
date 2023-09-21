/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

import (
	"time"
)

type Model struct {
	ID        uint      `json:"ID" gorm:"primaryKey;column:id"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime;column:created_at"`
}

var Models = []interface{}{
	&Contest{},
	&Round{},
	&Player{},
	&Score{},
	&Record{},
}
