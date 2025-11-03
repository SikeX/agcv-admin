package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/router/agvc/agvc_main"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System          system.RouterGroup
	Example         example.RouterGroup
	Agvc            agvc.RouterGroup
	AgvcMainHistory agvc_main.History
	AgvcMainDevice  agvc_main.Device
}
