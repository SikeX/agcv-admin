package agvc

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
    AgvcBwdSettingRouter
    AgvcQxySettingRouter
    AgvcNbqSettingRouter
    AgvcEventHisRouter
    AgvcNbqHisRouter
    AgvcBwdHisRouter
    AgvcQxyHisRouter
    AgvcChartRouter
}

var (
    agvcBwdSettingApi     = api.ApiGroupApp.AgvcApiGroup.AgvcBwdSettingApi
    agvcBwdSettingSyncApi = api.ApiGroupApp.AgvcApiGroup.AgvcBwdSettingSyncApi
    agvcQxySettingApi     = api.ApiGroupApp.AgvcApiGroup.AgvcQxySettingApi
    agvcNbqSettingApi     = api.ApiGroupApp.AgvcApiGroup.AgvcNbqSettingApi
    agvcEventHisApi       = api.ApiGroupApp.AgvcApiGroup.AgvcEventHisApi
    agvcNbqHisApi         = api.ApiGroupApp.AgvcApiGroup.AgvcNbqHisApi
    agvcBwdHisApi         = api.ApiGroupApp.AgvcApiGroup.AgvcBwdHisApi
    agvcQxyHisApi         = api.ApiGroupApp.AgvcApiGroup.AgvcQxyHisApi
    agvcChartApi          = api.ApiGroupApp.AgvcApiGroup.AgvcChartApi
)
