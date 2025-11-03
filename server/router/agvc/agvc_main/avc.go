package agvc_main

import (
	"github.com/gin-gonic/gin"
)

type avc struct{}

var AVC = new(avc)

func (r *avc) InitAVCRouter(public, private *gin.RouterGroup) {
	avcPublic := public.Group("avc")
	avcPrivate := private.Group("avc")
	{
		// 公开路由
		_ = avcPublic
	}
	{
		// 私有路由（需要鉴权）
		avcPrivate.GET("config", apiAVC.GetAVCConfig)
		avcPrivate.POST("config", apiAVC.CreateAVCConfig)
		avcPrivate.PUT("config", apiAVC.UpdateAVCConfig)
		avcPrivate.POST("start", apiAVC.StartAVC)
		avcPrivate.POST("stop", apiAVC.StopAVC)
		avcPrivate.GET("records", apiAVC.GetAVCRecords)
	}
}
