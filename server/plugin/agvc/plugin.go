package agvc

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/initialize"
	interfaces "github.com/flipped-aurora/gin-vue-admin/server/utils/plugin/v2"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

// Register 注册插件
func (p *plugin) Register(group *gin.Engine) {
	ctx := context.Background()
	
	// 初始化数据库表
	initialize.Gorm(ctx)
	
	// 初始化服务
	initialize.Service(ctx)
	
	// 初始化路由
	initialize.Router(group)
}

// RouterPath 返回插件的路由前缀
func (p *plugin) RouterPath() string {
	return "agvc"
}
