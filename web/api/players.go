/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:28.
 * Author:  guojia(https://github.com/guojia99)
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"my-cubing/db"
)

type (
	CreatePlayerRequest struct {
		Name  string `json:"Name"`
		WcaID string `json:"WcaID"`
	}

	CreatePlayerResponse struct {
		ID uint `json:"ID"`
	}

	player struct {
		Id    uint   `json:"ID"`
		Name  string `json:"Name"`
		WcaId string `json:"WcaId"`
	}

	ReadPlayersResponse struct {
		Data []player `json:"Data"`
	}
)

func CreatePlayer(ctx *gin.Context) {
	var req CreatePlayerRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var p *db.Player
	if err := db.DB.Model(&db.Player{}).Where("name = ?", req.Name).First(&p).Error; err == nil {
		p.WcaID = req.WcaID
		if err = db.DB.Save(&p).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			ctx.JSON(http.StatusOK, CreatePlayerResponse{ID: p.ID})
		}
		return
	}

	p = &db.Player{
		Name:  req.Name,
		WcaID: req.WcaID,
	}
	err := db.DB.Create(p).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, CreatePlayerResponse{ID: p.ID})
}

func ReadPlayers(ctx *gin.Context) {
	var players []db.Player
	if err := db.DB.Find(&players).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var resp = new(ReadPlayersResponse)
	for _, val := range players {
		resp.Data = append(resp.Data, player{
			Id:    val.ID,
			Name:  val.Name,
			WcaId: val.WcaID,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

func UpdatePlayers(ctx *gin.Context) {
	var req CreatePlayerRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var p *db.Player
	if err := db.DB.Model(&db.Player{}).Where("name = ?", req.Name).First(p).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	p.WcaID = req.WcaID
	p.Name = req.Name

	if err := db.DB.Updates(p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, CreatePlayerResponse{ID: p.ID})
}
