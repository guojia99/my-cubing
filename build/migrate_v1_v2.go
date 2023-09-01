/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/19 下午3:44.
 *  * Author: guojia(https://github.com/guojia99)
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	json "github.com/json-iterator/go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// v1 表
type (
	Model struct {
		ID        uint      `gorm:"primaryKey;column:id"`
		CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	}

	// Player 选手表
	Player struct {
		Model

		Name  string `json:"Name" gorm:"unique;not null;column:name"` // 选手名
		WcaID string `json:"WcaID" gorm:"column:wca_id"`              // 选手WcaID，用于查询选手WCA的成绩
	}

	// Contest 比赛表，记录某场比赛
	Contest struct {
		Model

		Name        string    `json:"Name" gorm:"unique;not null;column:name"`        // 比赛名
		Content     string    `json:"contest" gorm:"column:content;null"`             // 比赛论次描述, 为 ContestRoutes 的结构体json化
		Description string    `json:"Description" gorm:"not null;column:description"` // 描述
		IsEnd       bool      `json:"IsEnd" gorm:"null;column:is_end"`                // 是否已结束
		StartTime   time.Time `json:"StartTime" gorm:"column:start_time"`             // 开始时间
		EndTime     time.Time `json:"EndTime" gorm:"column:end_time"`                 // 结束时间
	}

	// Score 成绩表
	Score struct {
		ID           uint      `gorm:"primaryKey;column:id"`
		CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
		PlayerID     uint      `json:"PlayerID" gorm:"index;not null;column:player_id"`   // 选手的ID
		ContestID    uint      `json:"ContestID" gorm:"index;not null;column:contest_id"` // 比赛的ID
		RouteNumber  uint      `json:"RouteNumber" gorm:"not null;column:route_number"`   // 该项目的轮次
		Project      int       `json:"Project" gorm:"not null;column:project"`            // 分项目 333/222/444等
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
)

func GetToken() string {
	url := uri + "/v2/api/auth/token"
	method := "POST"

	message := map[string]interface{}{
		"user_name": "admin",
		"password":  "admin",
	}

	body, _ := json.Marshal(message)
	payload := bytes.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	var m map[string]interface{}
	_ = json.Unmarshal(data, &m)
	fmt.Println(m)
	return m["Token"].(string)
}

func AddPlayer(token string, name string, wcaID string) error {
	url := uri + "/v2/api/player"
	method := "POST"

	m := map[string]interface{}{
		"Name":  name,
		"WcaID": wcaID,
	}

	body, _ := json.Marshal(m)
	payload := bytes.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("token", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("error %s", string(body))
	}
	return nil
}

func CreateContest(token, name, description string) error {
	url := uri + "/v2/api/contest"
	method := "POST"

	m := map[string]interface{}{
		"Name":        name,
		"Description": description,
	}

	body, _ := json.Marshal(m)
	payload := bytes.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("token", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode > 400 {
		return fmt.Errorf("error")
	}
	return nil
}

func EndContest(token string, ContestID uint) error {
	url := uri + "/v2/api/score/end_contest"
	method := "PUT"

	m := map[string]interface{}{
		"ContestID": ContestID,
	}

	body, _ := json.Marshal(m)
	payload := bytes.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("token", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode > 400 {
		log.Printf("err %s\n", string(body))
		return fmt.Errorf("error")
	}
	return nil
}

func AddScore(token string, name string, contestId int, project string, num int, result []float64) error {
	url := uri + "/v2/api/score"
	method := "POST"

	m := map[string]interface{}{
		"PlayerName": name,
		"ContestID":  contestId,
		"Project":    project,
		"RouteNum":   num,
		"Results":    result,
	}

	body, _ := json.Marshal(m)
	payload := bytes.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("token", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("error %s", string(body))
	}
	return nil
}

const uri = "http://127.0.0.1:20000"

func main() {
	// 1. 查询所有角色， 并通过接口添加到新表中
	// 2. 查询到所有比赛，生成轮次数据
	// 3. 查询所有的成绩，按比赛写入不同的轮次及数据

	dbDSN := "root:my123456@tcp(127.0.0.1:3306)/mycube?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dbDSN}), &gorm.Config{})
	//
	//db, err := gorm.Open(sqlite.Open("back.db"), &gorm.Config{})
	if err != nil {
		log.Printf("error %s\n", err)
		return
	}

	token := GetToken()
	if token == "" {
		log.Printf("get token error\n")
		return
	}
	log.Printf("token %s\n", token)

	// 1. 所有角色入库
	var player []Player
	if err = db.Find(&player).Error; err != nil {
		log.Printf("error %s\n", err)
		return
	}
	var playerNames = make(map[uint]string)
	for _, p := range player {
		err = AddPlayer(token, p.Name, p.WcaID)
		if err != nil {
			log.Printf("[AddPlayer] error %s\n", err)
			return
		}
		playerNames[p.ID] = p.Name
	}

	// 2. 创建比赛， 记录成绩
	var contests []Contest
	if err = db.Find(&contests).Error; err != nil {
		log.Printf("error %s\n", err)
		return
	}
	for _, contest := range contests {
		if err = CreateContest(token, contest.Name, contest.Description); err != nil {
			log.Printf("[CreateContest] error %s\n", err)
			return
		}
		var scores []Score
		_ = db.Where("contest_id = ?", contest.ID).Find(&scores)
		for _, score := range scores {
			err = AddScore(token, playerNames[score.PlayerID], int(contest.ID), projectMap[score.Project], 1, []float64{
				score.Result1, score.Result2, score.Result3, score.Result4, score.Result5,
			})
			if err != nil {
				log.Printf("[AddScore] error %s\n", err)
				return
			}
		}

		if err = EndContest(token, contest.ID); err != nil {
			log.Printf("[EndContest] error %s\n", err)
			return
		}
	}
}

var projectMap = map[int]string{
	1:  "222",
	2:  "333",
	3:  "444",
	4:  "555",
	5:  "666",
	6:  "777",
	7:  "skewb",
	8:  "pyram",
	9:  "sq1",
	10: "minx",
	11: "clock",
	12: "333oh",
	13: "333fm",
	14: "333bf",
	15: "444bf",
	16: "555bf",
	17: "333mbf",
	18: "jhh",
	19: "o_cola",
	20: "333ft",
}
