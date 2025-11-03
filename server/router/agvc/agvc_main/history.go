package agvc_main

import (
	"github.com/gin-gonic/gin"
)

type History struct{}

// var History = new(history)

func (r *History) InitHistoryRouter(public, private *gin.RouterGroup) {
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
