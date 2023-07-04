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
	staticHtmlPath  = "./static/html/"
	staticIconPath  = "./static/favicon.ico"
	staticCssPath   = "./static/css"
	staticJsPath    = "./static/js"
	staticImagePath = "./static/image/"
)

func AddWebRoute(e *gin.Engine) {
	e.Static("/css", staticCssPath)
	e.Static("/js", staticJsPath)
	e.Static("/image", staticImagePath)

	e.StaticFile("/favicon.ico", staticIconPath)

	e.StaticFile("/admin/score", path.Join(staticHtmlPath, "score_admin.html")) // 成绩记录页

	e.StaticFile("/", path.Join(staticHtmlPath, "index.html"))           // 主页
	e.StaticFile("/score", path.Join(staticHtmlPath, "best_score.html")) // 所有项目汇总主页
	e.StaticFile("/contest", path.Join(staticHtmlPath, "contest.html"))  // 所有比赛的展示
}
