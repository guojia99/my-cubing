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
	s.SetResult([]float64{14.74, 11.67, 11.3, 12.92, 14.5}, ScorePenalty{})
	fmt.Println(Cube333.RouteType())
	fmt.Println(s.Best, s.Avg)
}

func TestScore_SetResult1(t *testing.T) {
	type fields struct {
		Project Project
	}
	type args struct {
		in      []float64
		penalty ScorePenalty
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantBest float64
		wantAvg  float64
	}{
		{
			name: "ok",
			fields: fields{
				Project: Cube333,
			},
			args: args{
				in:      []float64{1, 2, 3, 4, 5},
				penalty: ScorePenalty{},
			},
			wantBest: 1,
			wantAvg:  3,
		},
		{
			name: "has dnf",
			fields: fields{
				Project: Cube333,
			},
			args: args{
				in:      []float64{2, 3, 5, 4, DNF},
				penalty: ScorePenalty{},
			},
			wantBest: 2,
			wantAvg:  4,
		},
		{
			name: "v1",
			fields: fields{
				Project: CubePy,
			},
			args: args{
				in:      []float64{4.66, 5.28, 4.68, -10000, 5.05},
				penalty: ScorePenalty{},
			},
			wantBest: 0,
			wantAvg:  0,
		},
		{
			name: "v2",
			fields: fields{
				Project: Cube666,
			},
			args: args{
				in:      []float64{10, 20, 30, DNF, DNF},
				penalty: ScorePenalty{},
			},
			wantBest: 10,
			wantAvg:  20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Score{
				Project: tt.fields.Project,
			}
			s.SetResult(tt.args.in, tt.args.penalty)

			if s.Avg != tt.wantAvg {
				t.Errorf("got Avg %f error, want %f", s.Avg, tt.wantAvg)
			}
			if s.Best != tt.wantBest {
				t.Errorf("got Best %f error, want %f", s.Best, tt.wantBest)
			}
		})
	}
}

func TestSortScores(t *testing.T) {
	scores := []Score{
		{
			Project: Cube333,
			Best:    1,
			Avg:     DNF,
		},
		{
			Project: Cube333,
			Best:    DNF,
			Avg:     DNF,
		},
		{
			Project: Cube333,
			Best:    1,
			Avg:     3,
		},
		{
			Project: Cube333,
			Best:    1,
			Avg:     1,
		},
	}
	SortScores(scores)

	for _, val := range scores {
		fmt.Println(val.Best, val.Avg)
	}
}
