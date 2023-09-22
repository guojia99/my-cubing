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
	DNF = -10000 - iota // 未还原成功
	DNS                 // 未开始还原
)

// Score 成绩表
type Score struct {
	Model

	// 数据库字段
	PlayerID     uint    `json:"PlayerID,omitempty" gorm:"index;not null;column:player_id"`   // 选手的ID
	PlayerName   string  `json:"PlayerName,omitempty" gorm:"column:player_name"`              // 玩家名
	ContestID    uint    `json:"ContestID,omitempty" gorm:"index;not null;column:contest_id"` // 比赛的ID
	RouteID      uint    `json:"RouteID,omitempty" gorm:"index;column:route_id"`              // 轮次
	Project      Project `json:"Project,omitempty" gorm:"not null;column:project"`            // 分项目 333/222/444等
	Result1      float64 `json:"R1,omitempty" gorm:"column:r1;NULL"`                          // 成绩1 多盲时这个成绩是实际还原数
	Result2      float64 `json:"R2,omitempty" gorm:"column:r2;NULL"`                          // 成绩2 多盲时这个成绩是尝试复原数
	Result3      float64 `json:"R3,omitempty" gorm:"column:r3;NULL"`                          // 成绩3 多盲时这个成绩是计时
	Result4      float64 `json:"R4,omitempty" gorm:"column:r4;NULL"`                          // 成绩4
	Result5      float64 `json:"R5,omitempty" gorm:"column:r5;NULL"`                          // 成绩5
	Best         float64 `json:"Best,omitempty" gorm:"column:best;NULL"`                      // 五把最好成绩
	Avg          float64 `json:"Avg,omitempty" gorm:"column:avg;NULL"`                        // 五把平均成绩
	IsBestSingle bool    `json:"IsBestSingle,omitempty" gorm:"column:is_best_single"`         // 这是玩家比往期最佳的还好的成绩
	IsBestAvg    bool    `json:"IsBestAvg,omitempty" gorm:"column:is_best_avg"`               // 这是这个玩家比往期最佳的成绩还好的平均成绩
	Penalty      string  `json:"Penalty,omitempty" grom:"column:penalty"`                     // 判罚 ScorePenalty
	Rank         int     `json:"Rank,omitempty" gorm:"column:rank"`                           // 排名

	RouteValue Round `json:"RouteValue,omitempty" gorm:"-"` // 轮次实际内容
}

type ScorePenalty struct {
	R1 []int `json:"R1,omitempty"`
	R2 []int `json:"R2,omitempty"`
	R3 []int `json:"R3,omitempty"`
	R4 []int `json:"R4,omitempty"`
	R5 []int `json:"R5,omitempty"`
}

func (s *Score) DBest() bool { return s.Best <= DNF }
func (s *Score) DAvg() bool  { return s.Avg <= DNF }

func (s *Score) SetResult(in []float64, penalty ScorePenalty) {
	for len(in) < 5 {
		in = append(in, DNF)
	}

	s.Result1, s.Result2, s.Result3, s.Result4, s.Result5 = in[0], in[1], in[2], in[3], in[4]
	s.Best, s.Avg = DNF, DNF

	switch s.Project.RouteType() {
	case RouteType1rounds:
		s.Result1 = s.Result1 + float64(len(penalty.R1)*2)
		s.Result1, s.Best, s.Avg = in[0], in[0], in[0]
	case RouteType3roundsAvg, RouteType3roundsBest:
		s.Result1, s.Result2, s.Result3 =
			s.Result1+float64(len(penalty.R1)*2),
			s.Result2+float64(len(penalty.R2)*2),
			s.Result3+float64(len(penalty.R3)*2)

		cache := []float64{in[0], in[1], in[2]}
		sort.Slice(cache, func(i, j int) bool {
			if cache[i] <= DNF {
				return false
			}
			if cache[j] <= DNF {
				return true
			}
			return cache[i] < cache[j]
		})
		for i := 0; i < len(cache); i++ {
			if cache[i] <= DNF {
				continue
			}
			s.Best = cache[i]
			break
		}
		if s.D() == 0 {
			s.Avg = (s.Result1 + s.Result2 + s.Result3) / 3
		}
	case RouteType5roundsAvg, RouteType5roundsBest, RouteType5RoundsAvgHT:
		s.Result1, s.Result2, s.Result3, s.Result4, s.Result5 =
			s.Result1+float64(len(penalty.R1)*2),
			s.Result2+float64(len(penalty.R2)*2),
			s.Result3+float64(len(penalty.R3)*2),
			s.Result4+float64(len(penalty.R4)*2),
			s.Result5+float64(len(penalty.R5)*2)

		cache := in
		sort.Slice(cache, func(i, j int) bool {
			if cache[i] <= DNF {
				return false
			}
			if cache[j] <= DNF {
				return true
			}
			return cache[i] < cache[j]
		})

		for i := 0; i < len(cache); i++ {
			if cache[i] <= DNF {
				continue
			}
			s.Best = cache[i]
			break
		}

		if s.Project.RouteType() == RouteType5RoundsAvgHT || s.Project.RouteType() == RouteType5roundsBest {
			switch d := s.D(); d {
			case 0, 1:
				s.Avg = (cache[1] + cache[2] + cache[3]) / 3 // 正常去头尾
			}
			break
		}

		s.Avg = (s.Result1 + s.Result2 + s.Result3 + s.Result4 + s.Result5) / 5
		if s.D() > 0 {
			s.Avg = DNF
		}
	case RouteTypeRepeatedly:
		s.Result3 += float64(len(penalty.R3) * 2)
		// 还原低于两把
		if s.Result1 < 2 {
			break
		}
		// 如果还原次数没有大于失败数
		if (s.Result2 - s.Result1) > s.Result1 {
			break
		}
		s.Best = s.Result1 - (s.Result2 - s.Result1)
	}
}

func (s *Score) GetResult() []float64 {
	switch s.Project.RouteType() {
	case RouteType1rounds:
		return []float64{s.Result1}
	case RouteType3roundsAvg, RouteType3roundsBest, RouteTypeRepeatedly:
		return []float64{s.Result1, s.Result2, s.Result3}
	default:
		return []float64{s.Result1, s.Result2, s.Result3, s.Result4, s.Result5}
	}
}

func (s *Score) D() int {
	d := 0
	for _, val := range s.GetResult() {
		if val <= DNF {
			d += 1
		}
	}
	return d
}

func (s *Score) IsBestScore(other Score) bool {
	switch s.Project.RouteType() {
	case RouteTypeRepeatedly:
		// blind cube special rules:
		// - the result1 is number of successful recovery.
		// - the result2 is number of attempts to recover.
		// - the result3 is use times, (max back row).
		// - sort priority： r1 > r2 > r3
		// - like: if r1 and r2 equal, the best r3 is rank the top.
		if s.Result1 == other.Result1 {
			return s.Result2 < other.Result2 || s.Result3 < other.Result3
		}
		return s.Result1 > other.Result1
	default:
		if s.DBest() || other.DBest() {
			return !s.DBest()
		}
		if s.Best == other.Best {
			return s.Avg < other.Avg
		}
		return s.Best < other.Best
	}
}

func (s *Score) IsBestAvgScore(other Score) bool {
	switch s.Project.RouteType() {
	case RouteTypeRepeatedly:
		return true
	default:
		if s.DAvg() || other.DAvg() {
			return !s.DAvg()
		}
		if s.DAvg() && other.DAvg() {
			return s.IsBestScore(other)
		}
		return s.Avg < other.Avg
	}
}

// SortScores sort scores by WCA.
func SortScores(in []Score) {
	if len(in) == 0 {
		return
	}

	ro := in[0].Project.RouteType()
	sort.Slice(in, func(i, j int) bool {
		switch ro {
		case RouteType1rounds, RouteType3roundsBest, RouteType5roundsBest, RouteTypeRepeatedly:
			return in[i].IsBestScore(in[j])
		case RouteType3roundsAvg, RouteType5roundsAvg, RouteType5RoundsAvgHT:
			if in[i].DAvg() && in[j].DAvg() {
				return in[i].IsBestScore(in[j])
			}
			if in[i].DAvg() || in[j].DAvg() {
				return !in[i].DAvg()
			}
			return in[i].IsBestAvgScore(in[j])
		}
		return true
	})

	// add rank in scores, the identical score rank number equal.
	in[0].Rank = 1
	prev := in[0]
	for i := 1; i < len(in); i++ {
		switch ro {
		case RouteTypeRepeatedly:
			if in[i].Best == prev.Best && in[i].Result3 == prev.Result3 {
				in[i].Rank = prev.Rank
				break
			}
			in[i].Rank = prev.Rank + 1
		case RouteType1rounds:
			if in[i].Best == prev.Best {
				in[i].Rank = prev.Rank
				break
			}
			in[i].Rank = prev.Rank + 1
		case RouteType3roundsBest, RouteType5roundsBest, RouteType3roundsAvg, RouteType5roundsAvg, RouteType5RoundsAvgHT:
			if in[i].Best == prev.Best && in[i].Avg == prev.Avg {
				in[i].Rank = prev.Rank
				break
			}
			in[i].Rank = prev.Rank + 1
		}
		prev = in[i]
	}
}

func SortByBest(in []Score) {
	if len(in) == 0 {
		return
	}
	sort.Slice(in, func(i, j int) bool { return in[i].IsBestScore(in[j]) })

	in[0].Rank = 1
	prev := in[0]
	for i := 1; i < len(in); i++ {
		if in[i].Best == prev.Best {
			in[i].Rank = prev.Rank
		} else {
			in[i].Rank = prev.Rank + 1
		}
		prev = in[i]
	}
}

func SortByAvg(in []Score) {
	if len(in) == 0 {
		return
	}
	sort.Slice(in, func(i, j int) bool { return in[i].IsBestAvgScore(in[j]) })

	in[0].Rank = 1
	prev := in[0]
	for i := 1; i < len(in); i++ {
		if in[i].Avg == prev.Avg {
			in[i].Rank = prev.Rank
		} else {
			in[i].Rank = prev.Rank + 1
		}
		prev = in[i]
	}
}
