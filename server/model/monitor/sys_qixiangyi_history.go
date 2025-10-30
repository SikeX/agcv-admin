
// 自动生成模板SysQixiangyiHistory
package monitor
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 气象仪监控 结构体  SysQixiangyiHistory
type SysQixiangyiHistory struct {
    global.GVA_MODEL
    Name        *string `json:"name" form:"name" gorm:"column:name;"`               //气象仪名称
    QixiangyiID int     `json:"qixiangyiID" form:"qixiangyiID" gorm:"column:qixiangyi_id;"` // 气象仪ID (1或2)
}


// TableName 气象仪监控 SysQixiangyiHistory自定义表名 sys_qixiangyi_history
func (SysQixiangyiHistory) TableName() string {
    return "sys_qixiangyi_history"
}





