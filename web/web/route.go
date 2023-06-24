/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:30.
 * Author:  guojia(https://github.com/guojia99)
 */

package web

import (
	"path"

	"github.com/gin-gonic/gin"
)

const (
	staticHtmlPath = "./static/html/"
	staticIconPath = "./static/favicon.ico"
	staticCssPath  = "./static/css"
	staticJsPath   = "./static/js"
)

func AddWebRoute(e *gin.Engine) {
	e.Static("/css", staticCssPath)
	e.Static("/js", staticJsPath)

	e.StaticFile("/favicon.ico", staticIconPath)
	e.StaticFile("/", path.Join(staticHtmlPath, "report.html"))
	e.StaticFile("/score", path.Join(staticHtmlPath, "score.html"))   // 成绩记录页
	e.StaticFile("/report", path.Join(staticHtmlPath, "report.html")) // 报告页
}
