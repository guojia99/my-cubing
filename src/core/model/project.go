/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

func AllProjectRoute() []Project {
	return append(WCAProjectRoute(), XCubeProjectRoute()...)
}

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
		Cube333Ft,
	}
}

func XCubeProjectRoute() []Project {
	return []Project{
		XCube222BF,
		XCube666BF,
		XCube777BF,
		XCube333Mini,
		XCube222OH,
		XCube333MiniOH,
		XCube444OH,
		XCube555OH,
		XCube666OH,
		XCube777OH,
		XCubeSkOH,
		XCubePyOH,
		XCubeSq1OH,
		XCubeMinxOH,
		XCube333Mirror,
		XCube333Mirroring,
		XCube333Multiple5,
		XCube333Multiple10,
		XCube333Multiple15,
		XCube333Multiple20,
		XCube27Relay,
		XCube345RelayBF,
		XCubeAlienRelay,
		XCube27AlienRelayAll,
		XCube333Ghost,
		XCube333ZongZi,
	}
}

type Project string

func (p Project) Cn() string { return cubeProjectCnCode[p] }

const (
	JuBaoHaoHao Project = "jhh"
	OtherCola   Project = "o_cola"

	Cube222    Project = "222"
	Cube333    Project = "333"
	Cube444    Project = "444"
	Cube555    Project = "555"
	Cube666    Project = "666"
	Cube777    Project = "777"
	CubeSk     Project = "skewb"
	CubePy     Project = "pyram"
	CubeSq1    Project = "sq1"
	CubeMinx   Project = "minx"
	CubeClock  Project = "clock"
	Cube333OH  Project = "333oh"
	Cube333FM  Project = "333fm"
	Cube333BF  Project = "333bf"
	Cube444BF  Project = "444bf"
	Cube555BF  Project = "555bf"
	Cube333MBF Project = "333mbf"
	Cube333Ft  Project = "333ft"

	XCube222BF           Project = "222bf"
	XCube666BF           Project = "666bf"
	XCube777BF           Project = "777bf"
	XCube333Mini         Project = "333mini"
	XCube222OH           Project = "222oh"
	XCube333MiniOH       Project = "333mini_oh"
	XCube444OH           Project = "444oh"
	XCube555OH           Project = "555oh"
	XCube666OH           Project = "666oh"
	XCube777OH           Project = "777oh"
	XCubeSkOH            Project = "skewb_oh"
	XCubePyOH            Project = "pyram_oh"
	XCubeSq1OH           Project = "sql_oh"
	XCubeMinxOH          Project = "minx_oh"
	XCube333Mirror       Project = "333mirror"
	XCube333Mirroring    Project = "333mirroring"
	XCube333Multiple5    Project = "333multiple5"
	XCube333Multiple10   Project = "333multiple10"
	XCube333Multiple15   Project = "333multiple15"
	XCube333Multiple20   Project = "333multiple20"
	XCube27Relay         Project = "2_7relay"
	XCube345RelayBF      Project = "345relay_bf"
	XCubeAlienRelay      Project = "alien_relay"
	XCube27AlienRelayAll Project = "27alien_relay"
	XCube333Ghost        Project = "333ghost"
	XCube333ZongZi       Project = "333Zongzi"
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
	Cube333Ft:  "脚拧",

	XCube222BF:           "二盲",
	XCube666BF:           "六盲",
	XCube777BF:           "七盲",
	XCube333Mini:         "三阶迷你",
	XCube333MiniOH:       "三阶迷你单",
	XCube222OH:           "二单",
	XCube444OH:           "四单",
	XCube555OH:           "五单",
	XCube666OH:           "六单",
	XCube777OH:           "七单",
	XCubeSkOH:            "斜转单",
	XCubePyOH:            "金字塔单",
	XCubeSq1OH:           "SQ1单",
	XCubeMinxOH:          "五魔单",
	XCube333Mirror:       "镜面魔方",
	XCube333Mirroring:    "镜向三阶",
	XCube333Multiple5:    "三阶五连",
	XCube333Multiple10:   "三阶十连",
	XCube333Multiple15:   "三阶十五连",
	XCube333Multiple20:   "三阶二十连",
	XCube27Relay:         "正阶连拧",
	XCube345RelayBF:      "盲连拧",
	XCubeAlienRelay:      "异形连拧",
	XCube27AlienRelayAll: "全项目连拧",
	XCube333Ghost:        "鬼魔",
	XCube333ZongZi:       "粽子魔方",
}

func (p Project) IsBF() bool {
	switch p {
	case Cube333BF, Cube444BF, Cube555BF, XCube222BF, XCube666BF, XCube777BF:
		return true
	}
	return false
}

func (p Project) Route() int {
	switch p {
	case JuBaoHaoHao, OtherCola, XCube333Multiple5, XCube333Multiple10, XCube333Multiple15, XCube333Multiple20,
		XCube27Relay, XCube345RelayBF, XCubeAlienRelay, XCube27AlienRelayAll:
		return 1
	case Cube222, Cube333, Cube444, Cube555,
		CubeSk, CubePy, CubeSq1, CubeMinx, CubeClock, Cube333OH, Cube333Ft,
		XCube333Mini, XCube222OH, XCube333MiniOH, XCubeSkOH, XCubePyOH, XCube444OH, XCube555OH,
		XCubeSq1OH, XCubeMinxOH, XCube333Mirror, XCube333Mirroring,
		XCube333Ghost, XCube333ZongZi:
		return 5
	case Cube666, Cube777, Cube333FM,
		XCube666OH, XCube777OH,
		Cube333BF, Cube444BF, Cube555BF, XCube222BF, XCube666BF, XCube777BF:
		return 3
	case Cube333MBF:
		return -1
	default:
		return 0
	}
}
