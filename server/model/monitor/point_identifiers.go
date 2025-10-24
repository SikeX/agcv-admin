package monitor

// PointIdentifier 点标识配置
// 定义并网点对应的点标识范围和数据名称
const (
	// 并网点1的点标识范围 (1-20)
	GridPoint1Min = 1
	GridPoint1Max = 20
	
	// 并网点2的点标识范围 (21-25, 40-43) - 根据开关柜数据规范
	GridPoint2Min = 21
	GridPoint2Max = 43
)

// PointIdentifierInfo 点标识信息结构体
type PointIdentifierInfo struct {
	PointID   int    `json:"pointID"`   // 点标识ID
	PointName string `json:"pointName"` // 点标识名称
}

// PointNameMap 点标识与中文名称映射 - 根据开关柜（保护装置）数据规范
var PointNameMap = map[int]string{
	// 并网点1的点标识名称 (1-20)
	1:  "AB线电压Uab",
	2:  "BC线电压Ubc", 
	3:  "CA线电压Uca",
	4:  "A相电流Ia",
	5:  "B相电流Ib",
	6:  "C相电流Ic",
	7:  "有功功率P",
	8:  "无功功率Q",
	9:  "PF(COS功率因数)",
	10: "F（频率）",
	11: "S（视在功率）",
	12: "A相电压Ua",
	13: "B相电压Ub",
	14: "C相电压Uc",
	15: "零序电流（A）",
	16: "PA",
	17: "PB",
	18: "PC",
	19: "QA",
	20: "QB",
	
	// 并网点2的点标识名称 (21-43)
	21: "QC",
	22: "SA",
	23: "SB",
	24: "SC",
	25: "零序电压（V）",
	40: "正向有功",
	41: "正向无功",
	42: "反向有功",
	43: "反向无功",
}

// GetPointIdentifiersByGridPoint 根据并网点ID获取对应的点标识列表
func GetPointIdentifiersByGridPoint(gridPointID int) []int {
	var identifiers []int
	
	switch gridPointID {
	case 1:
		for i := GridPoint1Min; i <= GridPoint1Max; i++ {
			identifiers = append(identifiers, i)
		}
	case 2:
		for i := GridPoint2Min; i <= GridPoint2Max; i++ {
			identifiers = append(identifiers, i)
		}
	}
	
	return identifiers
}

// GetPointName 根据点标识获取中文名称
func GetPointName(pointID int) string {
	if name, exists := PointNameMap[pointID]; exists {
		return name
	}
	return "未知点位"
}