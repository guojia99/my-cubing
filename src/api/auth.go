/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/17 下午5:06.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	UserName string    `gorm:"column:user_name"`
	Password string    `gorm:"column:password"`
	Token    string    `gorm:"column:token"`
	Timeout  time.Time `gorm:"column:timeout"`
}

func (c *Client) initAuth() {
	_ = c.svc.DB.AutoMigrate(&Admin{})
}

type (
	GetTokenRequest struct {
		UserName string `json:"user_name"`
		PassWord string `json:"password"`
	}

	GetTokenResponse struct {
		Ts    int64  `json:"Ts"`
		Token string `json:"Token"`
	}
)

// ValidToken 获取合法token
func (c *Client) ValidToken(ctx *gin.Context) {
	var req GetTokenRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": err.Error()})
		return
	}
	if req.UserName == "" || req.PassWord == "" {
		ctx.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": "request has empty"})
		return
	}

	var admin Admin
	if err := c.svc.DB.Where("user_name = ?", req.UserName).First(&admin).Error; err != nil {
		ctx.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": err.Error()})
		return
	}
	if admin.Password != req.PassWord {
		ctx.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": "password error"})
		return
	}
	if time.Now().Sub(admin.Timeout) > 0 || admin.Token == "" {
		admin.Token = generateUniqueToken(req.UserName, time.Now().Unix())
		admin.Timeout = time.Now().Add(time.Hour * 48)
		c.svc.DB.Save(&admin)
	}

	resp := GetTokenResponse{
		Ts:    admin.Timeout.Unix(),
		Token: admin.Token,
	}
	ctx.JSON(http.StatusOK, resp)
}

func generateUniqueToken(username string, ts int64) string {
	uniqueStr := username + strconv.FormatInt(ts, 10)
	uniqueBytes := []byte(uniqueStr)
	randomBytes := make([]byte, 255-len(uniqueBytes))
	_, _ = rand.Read(randomBytes)
	tokenBytes := append(uniqueBytes, randomBytes...)
	return base64.URLEncoding.EncodeToString(tokenBytes)
}

// AuthMiddleware 授权中间件
func (c *Client) AuthMiddleware(ctx *gin.Context) {
	log.Printf("auth in %s", ctx.ClientIP())
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth error"})
		return
	}

	// token 可用，直接通过
	if _, ok := c.svc.Cache.Get(token); ok {
		ctx.Next()
		return
	}

	// 查token 是否合法
	var admin Admin
	if err := c.svc.DB.Where("token = ?", token).First(&admin).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("auth error token %s", err)})
		return
	}
	if time.Now().Sub(admin.Timeout) > 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token timeout"})
		return
	}

	_ = c.svc.Cache.Add(admin.Token, admin.UserName, time.Minute)
	ctx.Next()
}
