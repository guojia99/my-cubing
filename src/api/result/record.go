package result

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/guojia99/my-cubing/src/core/model"
	"github.com/guojia99/my-cubing/src/svc"
)

func GetRecords(svc *svc.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(ctx.DefaultQuery("size", "1000"))

		if size > 100 {
			size = 100
		}

		offset := (page - 1) * size
		limit := size

		key := fmt.Sprintf("GetRecords_%d_%d", page, size)
		if val, ok := svc.Cache.Get(key); ok {
			ctx.JSON(http.StatusOK, val)
			return
		}

		var records []model.Record
		svc.DB.Order("created_at DESC").Order("id DESC").Offset(offset).Limit(limit).Find(&records)

		for i := 0; i < len(records); i++ {
			var score model.Score
			svc.DB.First(&score, "id = ?", records[i].ScoreId)
			records[i].ScoreValue = score
		}

		_ = svc.Cache.Add(key, records, time.Second*30)
		ctx.JSON(http.StatusOK, records)
	}
}
