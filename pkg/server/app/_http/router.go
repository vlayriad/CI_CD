package _http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"kaffein/pkg/server/app/_http/handler"
)

type RouterHTTP struct {
	engine *gin.Engine
}

func NewRouterHTTP(engine *gin.Engine) *RouterHTTP {
	return &RouterHTTP{engine: engine}
}

func (r *RouterHTTP) ConfigRouterHTTP() {
	log.Info().Msg("Configuring HTTP routes...")

	pingHandler := handler.NewPingHandler()

	api := r.engine.Group("/api")
	{
		api.GET("/ping", pingHandler.Ping)
	}
}
