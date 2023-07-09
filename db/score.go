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
	Result1      float64   `json:"R1" gorm:"column:r1;NULL"`                          // 成绩1 多盲时这个成绩是实际还原数
	Result2      float64   `json:"R2" gorm:"column:r2;NULL"`                          // 成绩2 多盲时这个成绩是尝试复原数
	Result3      float64   `json:"R3" gorm:"column:r3;NULL"`                          // 成绩3 多盲时这个成绩是计时
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

		sort.Slice(in, func(i, j int) bool { return in[i] < in[j] })
		for i := 0; i < len(in); i++ {
			if in[i] != 0 {
				s.Best = in[i]
				break
			}
		}

		dnf := s.GetDNF()
		switch {
		case dnf == 1: // 有一把D的情况下, 去掉最好成绩后取平均
			s.Avg = (in[2] + in[3] + in[4]) / 3
		case dnf >= 2: // 两把以上D直接无平均
			s.Avg = 0
		default: // 正常去头尾
			s.Avg = (in[1] + in[2] + in[3]) / 3
			fmt.Println(s.Avg, in[1], in[2], in[3])
			fmt.Println(in)
		}
		return nil
	case Cube666, Cube777, Cube333FM, Cube333BF, Cube444BF, Cube555BF: // 三次的项目
		if len(in) < 3 {
			return fmt.Errorf("该项目需要输入3个成绩")
		}
		s.Result1, s.Result2, s.Result3 = in[0], in[1], in[2]
		if s.GetDNF() == 5 {
			return nil
		}
		cache := []float64{in[0], in[1], in[2]}
		sort.Slice(cache, func(i, j int) bool { return cache[i] < cache[j] })

		hasAvg := true
		for i := 0; i < len(cache); i++ {
			if cache[i] != 0 {
				s.Best = cache[i]
				break
			}
			hasAvg = false
		}

		if hasAvg {
			s.Avg = (s.Result1 + s.Result2 + s.Result3) / 3
		}
	case Cube333MBF: // 多盲特殊规则
		if len(in) < 3 {
			return fmt.Errorf("该项目需要3个成绩")
		}
		s.Result1, s.Result2, s.Result3 = in[0], in[1], in[2]
		if s.Result1 != 0 {
			s.Best = s.Result1
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

func (s *Score) IsBestScore(other Score) bool {
	switch s.Project {
	case Cube333MBF:
		// 多盲特殊评分规则
		// 1. 还原数多的, 排名优先
		// 2. 还原数相等, 尝试还原数少的排名优先
		// 3. 还原数和尝试还原数相等, 则还原时间少的优先
		if s.Result1 == other.Result1 {
			if s.Result2 < other.Result2 {
				return true
			}
			return s.Result3 < other.Result3
		}
		return s.Result1 > other.Result1
	default:
		if s.Best == 0 || other.Best == 0 {
			return s.Best == 0
		}
		return s.Best < other.Best
	}
}

func (s *Score) IsBestAvgScore(other Score) bool {
	switch s.Project {
	case Cube333MBF:
		// 多盲没有平均
		return true
	default:
		if s.Avg == 0 || other.Avg == 0 {
			return s.Best == 0
		}
		return s.Avg < other.Avg
	}
}
