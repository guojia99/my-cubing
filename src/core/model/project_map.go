package model

type projectItem struct {
	Project   Project
	Cn        string
	RouteType RouteType
	IsWca     bool
}

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
	XCubePyBF            Project = "pyram_bf"
	XCubePyFm            Project = "pyram_fm"
	XCubePySk            Project = "skewb_fm"
	XCube333Mini         Project = "333mini"
	XCube222OH           Project = "222oh"
	XCube333MiniOH       Project = "333mini_oh"
	XCube444OH           Project = "444oh"
	XCube555OH           Project = "555oh"
	XCube666OH           Project = "666oh"
	XCube777OH           Project = "777oh"
	XCubeSkOH            Project = "skewb_oh"
	XCubePyOH            Project = "pyram_oh"
	XCubeSq1OH           Project = "sq1_oh"
	XCubeMinxOH          Project = "minx_oh"
	XCubeClockOH         Project = "clock_oh"
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
	Xcube333Clone        Project = "333clone"
)

var projectsItems = []projectItem{
	{Project: JuBaoHaoHao, Cn: "菊爆浩浩", RouteType: RouteType1rounds, IsWca: false},
	{Project: OtherCola, Cn: "速可乐", RouteType: RouteType1rounds, IsWca: false},

	{Project: Cube222, Cn: "二阶", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: Cube333, Cn: "三阶", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: Cube444, Cn: "四阶", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: Cube555, Cn: "五阶", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: Cube666, Cn: "六阶", RouteType: RouteType3roundsAvg, IsWca: true},
	{Project: Cube777, Cn: "七阶", RouteType: RouteType3roundsAvg, IsWca: true},
	{Project: CubeSk, Cn: "斜转", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: CubePy, Cn: "金字塔", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: CubeSq1, Cn: "SQ1", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: CubeMinx, Cn: "五魔方", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: CubeClock, Cn: "魔表", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: Cube333OH, Cn: "单手", RouteType: RouteType5RoundsAvgHT, IsWca: true},
	{Project: Cube333FM, Cn: "最少步", RouteType: RouteType3roundsAvg, IsWca: true},
	{Project: Cube333BF, Cn: "三盲", RouteType: RouteType3roundsBest, IsWca: true},
	{Project: Cube444BF, Cn: "四盲", RouteType: RouteType3roundsBest, IsWca: true},
	{Project: Cube555BF, Cn: "五盲", RouteType: RouteType3roundsBest, IsWca: true},
	{Project: Cube333MBF, Cn: "多盲", RouteType: RouteTypeRepeatedly, IsWca: true},
	{Project: Cube333Ft, Cn: "脚拧", RouteType: RouteType5RoundsAvgHT, IsWca: true},

	{Project: XCube222BF, Cn: "二盲", RouteType: RouteType5roundsBest, IsWca: false},
	{Project: XCube666BF, Cn: "六盲", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube777BF, Cn: "七盲", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCubePyBF, Cn: "塔盲", RouteType: RouteType3roundsBest, IsWca: false},
	{Project: XCubePyFm, Cn: "塔少步", RouteType: RouteType3roundsAvg, IsWca: false},
	{Project: XCubePySk, Cn: "斜少步", RouteType: RouteType3roundsAvg, IsWca: false},
	{Project: XCube333Mini, Cn: "三阶迷你", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube333MiniOH, Cn: "三阶迷你单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube222OH, Cn: "二单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube444OH, Cn: "四单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube555OH, Cn: "五单", RouteType: RouteType3roundsAvg, IsWca: false},
	{Project: XCube666OH, Cn: "六单", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube777OH, Cn: "七单", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCubeSkOH, Cn: "斜转单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCubePyOH, Cn: "金字塔单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCubeSq1OH, Cn: "SQ1单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCubeMinxOH, Cn: "五魔单", RouteType: RouteType3roundsAvg, IsWca: false},
	{Project: XCubeClockOH, Cn: "表单", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube333Mirror, Cn: "镜面魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube333Mirroring, Cn: "镜向三阶", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube333Multiple5, Cn: "三阶五连", RouteType: RouteType5roundsBest, IsWca: false},
	{Project: XCube333Multiple10, Cn: "三阶十连", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube333Multiple15, Cn: "三阶十五连", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube333Multiple20, Cn: "三阶二十连", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube27Relay, Cn: "正阶连拧", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube345RelayBF, Cn: "盲连拧", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCubeAlienRelay, Cn: "异形连拧", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube27AlienRelayAll, Cn: "全项目连拧", RouteType: RouteType1rounds, IsWca: false},
	{Project: XCube333Ghost, Cn: "鬼魔", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: XCube333ZongZi, Cn: "粽子魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false},
	{Project: Xcube333Clone, Cn: "三阶克隆", RouteType: RouteType3roundsBest, IsWca: false},
}

var projectItemsMap = make(map[Project]projectItem, len(projectsItems))
var allProjectRoute []Project
var xCubeProjectRoute []Project
var wcaProjectRoute []Project

func init() {
	for _, val := range projectsItems {
		allProjectRoute = append(allProjectRoute, val.Project)
		if val.IsWca {
			wcaProjectRoute = append(wcaProjectRoute, val.Project)
		} else {
			xCubeProjectRoute = append(xCubeProjectRoute, val.Project)
		}
		projectItemsMap[val.Project] = val
	}
}
