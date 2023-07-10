/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:45.
 * Author:  guojia(https://github.com/guojia99)
 */

package model

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}

var Models = []interface{}{
	&Contest{},
	&Round{},
	&Player{},
	&Score{},
}
