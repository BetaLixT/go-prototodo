package server

import (
	"embed"
	"net/http"
	"prototodo/pkg/app/server/common"
	"prototodo/pkg/app/server/contracts"
	"strconv"

	"github.com/betalixt/gorr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed static/swagger/*
var swaggerFiles embed.FS

func (a *app) startHTTP(portStr string) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	router.SetTrustedProxies(nil)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "alive",
		})
	})

	fileServer := http.FileServer(http.FS(swaggerFiles))
	swaggerGroup := router.Group("/swagger")
	swaggerGroup.Any("/*all", func(ctx *gin.Context) {
		defer func(old string) {
			ctx.Request.URL.Path = old
		}(ctx.Request.URL.Path)

		ctx.Request.URL.Path = "/static" + ctx.Request.URL.Path
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
	})

	application := router.Group("")
	application.Use(func(ctx *gin.Context) {
		traceparent := ctx.GetHeader("traceparent")
		c := a.ctxf.Create(traceparent)
		ctx.Set(contracts.InternalContextKey, c)
		ctx.Next()

		er := ctx.Errors.Last()
		if er != nil {
			if err, ok := er.Err.(*gorr.Error); ok {
				ctx.JSON(err.StatusCode, err)
			} else {
				ctx.JSON(500, gorr.NewUnexpectedError(er))
			}
		}
	})

	port, err := strconv.Atoi(portStr)
	if err != nil {
		a.lgr.Warn(
			"unable to parse provided port, setting port to default",
			zap.String("portConfig", portStr),
		)
		port = common.HTTPDefaultPort
	}
	if port < 0 {
		a.lgr.Warn(
			"port was configured was invalid, setting port to default",
		)
		port = common.HTTPDefaultPort
	}

	srv := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	a.registerHTTPHandlers(application)
	a.registerCloser(func() {
		if err := srv.Close(); err != nil {
			a.lgr.Error("failed while closing http server", zap.Error(err))
		}
	})

	a.lgr.Info("http listening", zap.Int("port", port))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.lgr.Error(
			"failed running router",
			zap.Error(err),
		)
	}
}
