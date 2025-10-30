package agvc

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	AgvcBwdSettingRouter
	AgvcQxySettingRouter
}

var (
	agvcBwdSettingApi = api.ApiGroupApp.AgvcApiGroup.AgvcBwdSettingApi
	agvcQxySettingApi = api.ApiGroupApp.AgvcApiGroup.AgvcQxySettingApi
)
