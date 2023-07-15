/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

import (
	"fmt"
	"testing"
)

func TestScore_SetResult(t *testing.T) {
	var s = &Score{
		Project: Cube333,
	}
	err := s.SetResult([]float64{14.74, 11.67, 11.3, 12.92, 14.5})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
