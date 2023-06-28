/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:48.
 *  * Author: guojia(https://github.com/guojia99)
 */

package db

import "fmt"

func StrToProject(str string) Project {
	for key, val := range cubeProjectEnCode {
		if val == str {
			return key
		}
	}
	for key, val := range cubeProjectCnCode {
		if val == str {
			return key
		}
	}
	return 0
}

type Project int

func (p Project) String() string {
	return cubeProjectEnCode[p]
}

func (p Project) Cn() string {
	return cubeProjectCnCode[p]
}

func (p Project) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", p.Cn())), nil
}

const (
	JuBaoHaoHao Project = iota + 1
	Cube222
	Cube333
	Cube444
	Cube555
	Cube666
	Cube777
	CubeSk
	CubePy
	CubeSq1
	CubeMinx
	CubeClock
	Cube333OH
	Cube333FM
	Cube333BF
	Cube444BF
	Cube555BF
	Cube333MBF

	OtherCola
)

var cubeProjectEnCode = map[Project]string{
	JuBaoHaoHao: "jhh",
	OtherCola:   "o_cola",

	Cube222:    "222",
	Cube333:    "333",
	Cube444:    "444",
	Cube555:    "555",
	Cube666:    "666",
	Cube777:    "777",
	CubeSk:     "skewb",
	CubePy:     "pyram",
	CubeSq1:    "sq1",
	CubeMinx:   "minx",
	CubeClock:  "clock",
	Cube333OH:  "333oh",
	Cube333FM:  "333fm",
	Cube333BF:  "333bf",
	Cube444BF:  "444bf",
	Cube555BF:  "555bf",
	Cube333MBF: "333mbf",
}

var cubeProjectCnCode = map[Project]string{
	JuBaoHaoHao: "菊爆浩浩",
	OtherCola:   "速可乐",

	Cube222:    "二阶",
	Cube333:    "三阶",
	Cube444:    "四阶",
	Cube555:    "五阶",
	Cube666:    "六阶",
	Cube777:    "七阶",
	CubeSk:     "斜转",
	CubePy:     "金字塔",
	CubeSq1:    "SQ1",
	CubeMinx:   "五魔方",
	CubeClock:  "魔表",
	Cube333OH:  "单手",
	Cube333FM:  "最少步",
	Cube333BF:  "三盲",
	Cube444BF:  "四盲",
	Cube555BF:  "五盲",
	Cube333MBF: "多盲",
}
