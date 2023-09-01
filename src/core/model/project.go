/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

type Project string

type RouteType string

const (
	RouteType1rounds      RouteType = "1_r"      // 单轮项目
	RouteType3roundsBest  RouteType = "3_r_b"    // 三轮取最佳
	RouteType3roundsAvg   RouteType = "3_r_a"    // 三轮取平均
	RouteType5roundsBest  RouteType = "5_r_b"    // 五轮取最佳
	RouteType5roundsAvg   RouteType = "5_r_a"    // 五轮取平均
	RouteType5RoundsAvgHT RouteType = "5_r_a_ht" // 五轮去头尾取平均
	RouteTypeRepeatedly   RouteType = "ry"       // 单轮多次还原项目, 成绩1:还原数; 成绩2:尝试数; 成绩3:时间;
)

func AllProjectRoute() []Project       { return allProjectRoute }
func WCAProjectRoute() []Project       { return wcaProjectRoute }
func XCubeProjectRoute() []Project     { return xCubeProjectRoute }
func (p Project) Cn() string           { return projectItemsMap[p].Cn }
func (p Project) RouteType() RouteType { return projectItemsMap[p].RouteType }

func AllProjectItem() []projectItem { return projectsItems }
