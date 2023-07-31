/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

func WCAProjectRoute() []Project {
	return []Project{
		Cube333,
		Cube222,
		Cube444,
		Cube555,
		Cube666,
		Cube777,
		Cube333BF,
		Cube333FM,
		Cube333OH,
		CubeClock,
		CubeMinx,
		CubePy,
		CubeSk,
		CubeSq1,
		Cube444BF,
		Cube555BF,
		Cube333MBF,
	}
}

type Project string

func (p Project) Cn() string { return cubeProjectCnCode[p] }

const (
	JuBaoHaoHao Project = "jhh"
	OtherCola   Project = "o_cola"
	Cube222     Project = "222"
	Cube333     Project = "333"
	Cube444     Project = "444"
	Cube555     Project = "555"
	Cube666     Project = "666"
	Cube777     Project = "777"
	CubeSk      Project = "skewb"
	CubePy      Project = "pyram"
	CubeSq1     Project = "sq1"
	CubeMinx    Project = "minx"
	CubeClock   Project = "clock"
	Cube333OH   Project = "333oh"
	Cube333FM   Project = "333fm"
	Cube333BF   Project = "333bf"
	Cube444BF   Project = "444bf"
	Cube555BF   Project = "555bf"
	Cube333MBF  Project = "333mbf"
)

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
