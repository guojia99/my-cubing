/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/17 下午6:13.
 *  * Author: guojia(https://github.com/guojia99)
 */

package result

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/core/model"
	"github.com/guojia99/my-cubing/src/svc"
)

type (
	PlayerRequest struct {
		Id uint `uri:"player_id"`
	}
	PlayerResponse struct {
		model.Player
	}
)

func GetPlayer(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PlayerRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		out := svc.Core.GetPlayerDetail(req.Id)
		if out.ID == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		ctx.JSON(http.StatusOK, out)
	}
}

type PlayersResponse struct {
	Players []model.Player `json:"Players"`
}

func GetPlayers(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp PlayersResponse
		if err := svc.DB.Find(&resp.Players).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for i := 0; i < len(resp.Players); i++ {
			resp.Players[i].GetTitles()
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func CreatePlayer(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.Player
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var player model.Player
		if err := svc.DB.Where("name = ?", req.Name).First(&player).Error; err == nil {
			player = model.Player{
				Model:      req.Model,
				Name:       req.Name,
				WcaID:      req.WcaID,
				ActualName: req.ActualName,
				Titles:     req.Titles,
			}
			player.SetTitles(req.TitlesVal)
			svc.DB.Save(&player)
			ctx.JSON(http.StatusOK, gin.H{})
			return
		}

		if err := svc.DB.Create(&req).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func DeletePlayer(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
