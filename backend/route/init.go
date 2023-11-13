package route

import (
	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/middleware/db"
	"github.com/apicat/apicat/backend/route/middleware/log"
	"github.com/apicat/apicat/frontend"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func Init() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.ContextWithFallback = true

	r.Use(
		log.RequestIDLog("/assets/", "/static/"),
		db.CheckDBConnStatus("/assets/", "/static/"),
		translator.UseValidatori18n(),
		gin.Recovery(),
	)

	t, _ := template.ParseFS(frontend.Dist, "dist/templates/*.tmpl")
	r.SetHTMLTemplate(t)

	// 前端单页路由以及静态资源
	frontendHandle(r)

	r.GET("/", func(ctx *gin.Context) {
		ctx.FileFromFS("dist/", http.FS(frontend.Dist))
	})

	registerMock(r)
	registerGetConfig(r)

	apiRouter := r.Group("/api")

	registerSetConfig(apiRouter)

	InitApiRouter(apiRouter)

	registerNoRoute(r)

	r.Run(config.GetSysConfig().App.Host.Value + ":" + config.GetSysConfig().App.Port.Value)
}
