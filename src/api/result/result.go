/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/19 下午10:26.
 *  * Author: guojia(https://github.com/guojia99)
 */

package result

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/core/model"
	"github.com/guojia99/my-cubing/src/svc"
)

type GetProjectMapResponse struct {
	Projects []model.Project          `json:"Projects"`
	En       map[model.Project]string `json:"En"`
	Cn       map[model.Project]string `json:"Cn"`
}

func GetProjectMap(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, GetProjectMapResponse{
			Projects: model.WCAProjectRoute(),
			En:       model.GetEnMap(),
			Cn:       model.GetCnMap(),
		})
	}
}
