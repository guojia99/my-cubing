/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:47.
 *  * Author: guojia(https://github.com/guojia99)
 */

package db

import (
	"fmt"
	"sort"
	"time"
)

// Score 成绩表
type Score struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
	PlayerID     uint      `json:"PlayerID" gorm:"index;not null;column:player_id"`   // 选手的ID
	ContestID    uint      `json:"ContestID" gorm:"index;not null;column:contest_id"` // 比赛的ID
	RouteNumber  uint      `json:"RouteNumber" gorm:"not null;column:route_number"`   // 该项目的轮次
	Project      Project   `json:"Project" gorm:"not null;column:project"`            // 分项目 333/222/444等
	Result1      float64   `json:"R1" gorm:"column:r1;NULL"`                          // 成绩1
	Result2      float64   `json:"R2" gorm:"column:r2;NULL"`                          // 成绩2
	Result3      float64   `json:"R3" gorm:"column:r3;NULL"`                          // 成绩3
	Result4      float64   `json:"R4" gorm:"column:r4;NULL"`                          // 成绩4
	Result5      float64   `json:"R5" gorm:"column:r5;NULL"`                          // 成绩5
	Best         float64   `json:"Best" gorm:"column:best;NULL"`                      // 五把最好成绩
	Avg          float64   `json:"Avg" gorm:"column:avg;NULL"`                        // 五把平均成绩
	IsBest       bool      `json:"IsBest" grom:"column:is_best;NULL"`                 // 这是比往期最佳的还好的成绩
	IsBestAvg    bool      `json:"IsBestAvg" grom:"column:is_best_avg;NULL"`          // 这是比往期最佳的成绩还好的平均成绩
	IsBestRecord bool      `json:"BestRecord" gorm:"column:is_best_record;NULL"`      // 打破了以往的最佳记录
	IsAvgRecord  bool      `json:"AvgRecord" gorm:"column:is_avg_record;NULL"`        // 打破了以往的平均记录
}

func (s *Score) SetResult(in []float64) error {
	if len(in) == 0 {
		return nil
	}

	switch s.Project {
	case JuBaoHaoHao, OtherCola:
		if len(in) <= 0 {
			return fmt.Errorf("需要输入一个成绩")
		}
		s.Result1 = in[0]
		s.Best = in[0]
		s.Avg = in[0]
		// 五次的项目
	case Cube222, Cube333, Cube444, Cube555, CubeSk, CubePy, CubeSq1, CubeMinx, CubeClock, Cube333OH:
		if len(in) != 5 {
			return fmt.Errorf("该项目需要输入5个成绩")
		}
		s.Result1, s.Result2, s.Result3, s.Result4, s.Result5 = in[0], in[1], in[2], in[3], in[4]
		dnf := s.GetDNF()

		sort.Slice(in, func(i, j int) bool { return in[i] > in[j] })
		s.Best = in[1]
		switch {
		case dnf == 5:
			return nil
		case dnf >= 2:
			s.Avg = 0
		default:
			// 去头尾取平均
			s.Avg = (in[1] + in[2] + in[3]) / 3
		}
		return nil
		// 三次的项目
	case Cube666, Cube777, Cube333FM, Cube333BF, Cube444BF, Cube555BF:
		if len(in) < 3 {
			return fmt.Errorf("该项目需要输入3个成绩")
		}
		if s.GetDNF() == 5 {
			return nil
		}
		s.Result1, s.Result2, s.Result3 = in[0], in[1], in[2]
		sort.Slice(in, func(i, j int) bool { return in[i] > in[j] })
		s.Avg = (s.Result1 + s.Result2 + s.Result3) / 3
		s.Best = in[0]
		// 二个成绩
	case Cube333MBF:
		if len(in) < 2 {
			return fmt.Errorf("该项目需要2个成绩")
		}
	}
	return nil
}

func (s *Score) GetResult() []float64 {
	return []float64{s.Result1, s.Result2, s.Result3, s.Result4, s.Result5}
}

func (s *Score) GetDNF() int {
	dnf := 0
	for _, val := range s.GetResult() {
		if val == 0 {
			dnf += 1
		}
	}
	return dnf
}
