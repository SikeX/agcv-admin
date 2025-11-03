package router

import (
	"github.com/gin-gonic/gin"
)

type agc struct{}

var AGC = new(agc)

func (r *agc) InitAGCRouter(public, private *gin.RouterGroup) {
	agcPublic := public.Group("agc")
	agcPrivate := private.Group("agc")
	{
		// 公开路由
		_ = agcPublic
	}
	{
		// 私有路由（需要鉴权）
		agcPrivate.GET("config", apiAGC.GetAGCConfig)
		agcPrivate.POST("config", apiAGC.CreateAGCConfig)
		agcPrivate.PUT("config", apiAGC.UpdateAGCConfig)
		agcPrivate.POST("start", apiAGC.StartAGC)
		agcPrivate.POST("stop", apiAGC.StopAGC)
		agcPrivate.GET("records", apiAGC.GetAGCRecords)
	}
}
