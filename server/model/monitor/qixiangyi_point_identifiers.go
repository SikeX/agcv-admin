package monitor

// QixiangyiPointIdentifierInfo 气象仪点标识信息结构体
type QixiangyiPointIdentifierInfo struct {
	PointID   int    `json:"pointID"`   // 点标识ID
	PointName string `json:"pointName"` // 点标识名称
}

// QixiangyiPointNameMap 气象仪点标识与中文名称映射
var QixiangyiPointNameMap = map[int]string{
	// 气象仪1的点标识名称 (1-20)
	2:  "环境温度(℃)",
	3:  "露点温度(℃)",
	4:  "风速(m/s)",
	5:  "2分风速(m/s)",
	6:  "10分风速(m/s)",
	7:  "风向(°)",
	8:  "水平辐射强度(W/㎡)",
	9:  "倾斜辐射强度(W/㎡)",
	11: "组件温度(℃)",
	12: "雨量(mm)",
	
	// 气象仪2的点标识名称 (21-52)
	21: "散射辐射强度(W/㎡)",
	22: "直射辐射强度(W/㎡)",
	39: "水平总辐射量(MJ/㎡)",
	40: "倾斜总辐射量(MJ/㎡)",
	46: "环境湿度(%RH)",
	51: "日照小时数(h)",
	52: "气压(Pa)",
}

// GetPointIdentifiersByQixiangyiId 根据气象仪ID获取对应的点标识列表
func GetPointIdentifiersByQixiangyiId(qixiangyiId int) []int {
	var identifiers []int
	
	// 根据气象仪ID返回对应的点标识
	// 气象仪1: ID <= 1，返回点标识 2,3,4,5,6,7,8,9,11,12
	// 气象仪2: ID > 1，返回点标识 21,22,39,40,46,51,52
	if qixiangyiId <= 1 {
		// 气象仪1支持的点标识
		identifiers = []int{2, 3, 4, 5, 6, 7, 8, 9, 11, 12}
	} else {
		// 气象仪2支持的点标识
		identifiers = []int{21, 22, 39, 40, 46, 51, 52}
	}
	
	return identifiers
}

// GetQixiangyiPointName 根据点标识获取中文名称
func GetQixiangyiPointName(pointID int) string {
	if name, exists := QixiangyiPointNameMap[pointID]; exists {
		return name
	}
	return "未知点位"
}