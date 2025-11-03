package router

import (
	"github.com/gin-gonic/gin"
)

type device struct{}

var Device = new(device)

func (r *device) InitDeviceRouter(public, private *gin.RouterGroup) {
	devicePublic := public.Group("device")
	devicePrivate := private.Group("device")
	{
		// 公开路由（如果需要）
		_ = devicePublic
	}
	{
		// 私有路由（需要鉴权）
		devicePrivate.POST("create", apiDevice.CreateDevice)
		// devicePrivate.DELETE("delete", apiDevice.DeleteDevice)
		devicePrivate.PUT("update", apiDevice.UpdateDevice)
		// devicePrivate.GET("find", apiDevice.GetDevice)
		devicePrivate.GET("list", apiDevice.GetDeviceList)
		devicePrivate.GET("realtimeData", apiDevice.GetDeviceRealtimeData)
		devicePrivate.GET("inverters", apiDevice.GetInverters)
	}
}
