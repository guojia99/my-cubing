package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing/src/svc"
)

type PlayerImages struct {
	gorm.Model
	PlayerID   uint   `json:"PlayerID"  gorm:"unique;not null;column:player_id"`
	Avatar     string `json:"Avatar"`
	Background string `json:"Background"`
}

func CreatePlayerImages(svc *svc.Context) gin.HandlerFunc {
	_ = svc.DB.AutoMigrate(&PlayerImages{})
	return func(ctx *gin.Context) {
		// todo 接入图片创建
	}
}

func GetPlayerImages(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PlayerRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusOK, gin.H{})
			return
		}

		var out PlayerImages
		if err := svc.DB.Where("player_id = ?", req.Id).First(&out).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{})
			return
		}
		ctx.JSON(http.StatusOK, out)
	}
}
