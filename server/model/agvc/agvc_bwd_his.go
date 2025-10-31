// 自动生成模板AgvcBwdHis
package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// agvcBwdHis表 结构体  AgvcBwdHis
type AgvcBwdHis struct {
	global.GVA_MODEL
	Number *string `json:"number" form:"number" gorm:"column:number;comment:设备编号"`
	Name *string `json:"name" form:"name" gorm:"column:name;comment:并网点名称"`
}

// TableName agvcBwdHis表 AgvcBwdHis自定义表名 agvc_bwd_his
func (AgvcBwdHis) TableName() string {
	return "agvc_bwd_his"
}
