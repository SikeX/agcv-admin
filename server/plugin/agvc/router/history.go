package router

import (
	"github.com/gin-gonic/gin"
)

type history struct{}

var History = new(history)

func (r *history) InitHistoryRouter(public, private *gin.RouterGroup) {
	historyPublic := public.Group("history")
	historyPrivate := private.Group("history")
	{
		// 公开路由
		_ = historyPublic
	}
	{
		// 私有路由（需要鉴权）
		historyPrivate.POST("query", apiHistory.QueryHistoryData)
	}
}
