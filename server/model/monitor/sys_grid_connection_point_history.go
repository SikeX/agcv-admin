
// 自动生成模板SysGridConnectionPointHistory
package monitor
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 并网点监控 结构体  SysGridConnectionPointHistory
type SysGridConnectionPointHistory struct {
    global.GVA_MODEL
    Name  *string `json:"name" form:"name" gorm:"column:name;"`  //并网点名字
}

// 历史数据查询请求结构体
type GridPointHistoryQuery struct {
    GridPointID int    `json:"gridPointID" form:"gridPointID"` // 并网点ID
    PointID     int    `json:"pointID" form:"pointID"`         // 点标识ID
    StartTime   string `json:"startTime" form:"startTime"`     // 开始时间
    EndTime     string `json:"endTime" form:"endTime"`         // 结束时间
}

// 历史数据响应结构体
type GridPointHistoryData struct {
    Time  string  `json:"time"`  // 时间戳
    Value float64 `json:"value"` // 数值
}

// 历史数据查询响应
type GridPointHistoryResponse struct {
    PointID   int                     `json:"pointID"`   // 点标识ID
    PointName string                  `json:"pointName"` // 点标识名称
    Data      []GridPointHistoryData  `json:"data"`      // 历史数据
}

// TableName 并网点监控 SysGridConnectionPointHistory自定义表名 sys_grid_connection_point_history
func (SysGridConnectionPointHistory) TableName() string {
    return "sys_grid_connection_point_history"
}





