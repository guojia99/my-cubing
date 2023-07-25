/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/17 下午5:26.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/util/cache"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func NewRateMiddleware(secondRete int) gin.HandlerFunc {
	c := *cache.NewLRUExpireCache(10000)
	return func(ctx *gin.Context) {
		clientIp := ctx.ClientIP()

		var r *rate.Limiter
		if val, ok := c.Get(clientIp); ok {
			r = val.(*rate.Limiter)
		} else {
			r = rate.NewLimiter(rate.Limit(1), secondRete)
		}

		if !r.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many req"})
			return
		}
		ctx.Next()
	}
}
