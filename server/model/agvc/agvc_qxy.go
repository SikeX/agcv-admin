
// 自动生成模板AgvcQxyHis
package agvc
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 气象仪监控 结构体  AgvcQxyHis
type AgvcQxyHis struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"column:name;"`  //设备名称
  Number  *string `json:"number" form:"number" gorm:"column:number;"`  //设备编号
}


// TableName 气象仪监控 AgvcQxyHis自定义表名 agvc_qxy_his
func (AgvcQxyHis) TableName() string {
    return "agvc_qxy_his"
}





