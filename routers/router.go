package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"webhook/middleware/prometheus"
	"webhook/pkg/setting"
	v1 "webhook/routers/api/v1"
)

func InitRouterHandler() *gin.Engine {
	g := gin.New()
	g.Use(prometheus.TimeDuration)
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	gin.SetMode(setting.Svc.App.Mode)
	apiV1 := g.Group("/api/v1")
	{
		apiV1.GET("/hooks", v1.GetWebHooks)
		apiV1.GET("/hook", v1.GetWebHook)
		apiV1.POST("/hooks", v1.ReceiveWebhook)
	}
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return g
}
