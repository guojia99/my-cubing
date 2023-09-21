package xlog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing/src/svc"
)

type XLog struct {
	gorm.Model

	Title       string `json:"Title,omitempty"`
	CreatedTime string `json:"CreatedTime,omitempty"`
	Area        string `json:"Area,omitempty"`
	Messages    string `json:"Messages,omitempty"`
}

func AddXLog(svc *svc.Context) gin.HandlerFunc {
	_ = svc.DB.AutoMigrate(&XLog{})

	return func(ctx *gin.Context) {
		var req XLog
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		svc.DB.Save(&XLog{
			Title:       req.Title,
			CreatedTime: req.CreatedTime,
			Area:        req.Area,
			Messages:    req.Messages,
		})
	}
}

func GetXLogs(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var out []XLog
		svc.DB.Find(&out)
		ctx.JSON(http.StatusOK, out)
	}
}

type DeleteXLogReq struct {
	ID uint `uri:"x_id"`
}

func DeleteXLog(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req DeleteXLogReq
		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		svc.DB.Delete(&XLog{}, "id = ?", req.ID)
	}
}
