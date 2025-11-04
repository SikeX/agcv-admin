package agvc_main

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type ScheduleRouter struct{}

// InitScheduleRouter 初始化计划曲线路由
func (s *ScheduleRouter) InitScheduleRouter(Router *gin.RouterGroup) {
	scheduleRouter := Router.Group("schedule")
	scheduleApi := v1.ApiGroupApp.AgvcMainApiGroup.Schedule
	{
		scheduleRouter.POST("create", scheduleApi.CreateSchedule)       // 创建计划曲线
		scheduleRouter.PUT("update", scheduleApi.UpdateSchedule)        // 更新计划曲线
		scheduleRouter.DELETE("delete", scheduleApi.DeleteSchedule)     // 删除计划曲线
		scheduleRouter.GET("list", scheduleApi.GetScheduleList)         // 获取计划曲线列表
	}
}
