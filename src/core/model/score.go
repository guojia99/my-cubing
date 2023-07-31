/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

import (
	"sort"
)

const (
	DNF = 0  // 未还原成功
	DNS = -1 // 未开始还原
)

// Score 成绩表
type Score struct {
	Model

	// 数据库字段
	PlayerID     uint    `json:"PlayerID" gorm:"index;not null;column:player_id"`   // 选手的ID
	PlayerName   string  `json:"PlayerName" gorm:"column:player_name"`              // 玩家名
	ContestID    uint    `json:"ContestID" gorm:"index;not null;column:contest_id"` // 比赛的ID
	RouteID      uint    `json:"RouteID" gorm:"index;column:route_id"`              // 轮次
	Project      Project `json:"Project" gorm:"not null;column:project"`            // 分项目 333/222/444等
	Result1      float64 `json:"R1" gorm:"column:r1;NULL"`                          // 成绩1 多盲时这个成绩是实际还原数
	Result2      float64 `json:"R2" gorm:"column:r2;NULL"`                          // 成绩2 多盲时这个成绩是尝试复原数
	Result3      float64 `json:"R3" gorm:"column:r3;NULL"`                          // 成绩3 多盲时这个成绩是计时
	Result4      float64 `json:"R4" gorm:"column:r4;NULL"`                          // 成绩4
	Result5      float64 `json:"R5" gorm:"column:r5;NULL"`                          // 成绩5
	Best         float64 `json:"Best" gorm:"column:best;NULL"`                      // 五把最好成绩
	Avg          float64 `json:"Avg" gorm:"column:avg;NULL"`                        // 五把平均成绩
	IsBestSingle bool    `json:"IsBestSingle" gorm:"column:is_best_single"`         // 这是玩家比往期最佳的还好的成绩
	IsBestAvg    bool    `json:"IsBestAvg" gorm:"column:is_best_avg"`               // 这是这个玩家比往期最佳的成绩还好的平均成绩

	RouteValue Round `json:"RouteValue" gorm:"-"` // 轮次实际内容
}

func (s *Score) SetResult(in []float64) error {
	for len(in) < 5 {
		in = append(in, DNF)
	}

	switch s.Project {
	case JuBaoHaoHao, OtherCola:
		s.Result1, s.Best, s.Avg = in[0], in[0], in[0]
	case Cube222, Cube333, Cube444, Cube555, CubeSk, CubePy, CubeSq1, CubeMinx, CubeClock, Cube333OH:
		s.Result1, s.Result2, s.Result3, s.Result4, s.Result5 = in[0], in[1], in[2], in[3], in[4]
		cache := in
		sort.Slice(cache, func(i, j int) bool { return cache[i] < cache[j] })
		for i := 0; i < len(cache); i++ {
			if cache[i] != 0 {
				s.Best = cache[i]
				break
			}
		}

		switch d := s.D(); d {
		case 1:
			s.Avg = (cache[2] + cache[3] + cache[4]) / 3 // 有一把D的情况下, 去掉最好成绩后取平均
		default:
			s.Avg = (cache[1] + cache[2] + cache[3]) / 3 // 正常去头尾
		}
	case Cube666, Cube777, Cube333FM, Cube333BF, Cube444BF, Cube555BF: // 三次的项目
		s.Result1, s.Result2, s.Result3 = in[0], in[1], in[2]
		cache := []float64{in[0], in[1], in[2]}
		sort.Slice(cache, func(i, j int) bool { return cache[i] < cache[j] })
		for i := 0; i < len(cache); i++ {
			if cache[i] > 0 {
				s.Best = cache[i]
				break
			}
		}
		if s.D() == 0 {
			s.Avg = (s.Result1 + s.Result2 + s.Result3) / 3
		}
	case Cube333MBF:
		s.Result1, s.Result2, s.Result3 = in[0], in[1], in[2]
		if s.Result1 >= 2 { // 两把以上才有最佳成绩, 该成绩才有效
			s.Best = s.Result1
		}
	}
	return nil
}

func (s *Score) GetResult() []float64 {
	switch s.Project {
	case Cube666, Cube777, Cube333FM, Cube333BF, Cube444BF, Cube555BF, Cube333MBF: // 三次的项目
		return []float64{s.Result1, s.Result2, s.Result3}
	case JuBaoHaoHao, OtherCola:
		return []float64{s.Result1}
	default:
		return []float64{s.Result1, s.Result2, s.Result3, s.Result4, s.Result5}
	}
}

func (s *Score) D() int {
	switch s.Project {
	case Cube333MBF:
		if s.Best == 0 {
			return 0
		}
		return 1
	default:
		d := 0
		for _, val := range s.GetResult() {
			if val <= 0 {
				d += 1
			}
		}
		return d
	}
}

func (s *Score) IsBestScore(other Score) bool {
	switch s.Project {
	case Cube333MBF:
		// 多盲特殊评分规则
		// 1. 还原数多的, 排名优先
		// 2. 还原数相等, 尝试还原数少的排名优先
		// 3. 还原数和尝试还原数相等, 则还原时间少的优先
		if s.Result1 == other.Result1 {
			return s.Result2 < other.Result2 || s.Result3 < other.Result3
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

// SortScores 给成绩列表按WCA名次排序
func SortScores(in []Score) {
	sort.Slice(in, func(i, j int) bool {
		switch in[i].Project {
		case Cube333MBF, Cube333BF, Cube444BF, Cube555BF: // 盲规则, 以最佳为准
			return in[i].IsBestScore(in[j])
		default:
			if in[i].Avg+in[j].Avg == 0 { // 都没有平均成绩
				return in[i].IsBestScore(in[j])
			}
			if in[i].Avg == in[j].Avg { // 平均成绩相同以最佳排名
				return in[i].IsBestScore(in[j])
			}
			return in[i].IsBestAvgScore(in[j]) // 其中一个有平均成绩的， 以平均成绩排序
		}
	})
}
