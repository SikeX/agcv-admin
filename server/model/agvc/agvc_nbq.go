// 自动生成模板AgvcNbqHis
package agvc

import (
    "reflect"
    "strings"
    "sync"
)

// agvcNbqHis表 结构体  AgvcNbq
// 用于存储和展示逆变器实时数据和历史数据
type AgvcNbq struct {
    Ctime      *string `json:"ctime" form:"ctime" point:"name:采集时间"`
    Psid       *int    `json:"psid" form:"psid" point:"name:电站编号"`
    InverterNo *int    `json:"inverterNo" form:"inverterNo" point:"name:逆变器编号"`
    Name       *string `json:"name" form:"name" point:"name:逆变器名称"`

    // 发电量相关
    TotalPowerGeneration   *float64 `json:"totalPowerGeneration" form:"totalPowerGeneration" point:"name:总发电量(kWh),value:1,type:2"`
    DailyPowerGeneration   *float64 `json:"dailyPowerGeneration" form:"dailyPowerGeneration" point:"name:日发电量(kWh),value:2,type:2"`
    MonthlyPowerGeneration *float64 `json:"monthlyPowerGeneration" form:"monthlyPowerGeneration" point:"name:月发电量(kWh),value:3,type:2"`
    AnnualPowerGeneration  *float64 `json:"annualPowerGeneration" form:"annualPowerGeneration" point:"name:年发电量(kWh),value:4,type:2"`

    // 功率相关
    ACPower         *float64 `json:"acPower" form:"acPower" point:"name:交流功率(kW),value:10,type:2"`
    DCPower         *float64 `json:"dcPower" form:"dcPower" point:"name:直流功率(kW),value:202,type:2"`
    ReactivePower   *float64 `json:"reactivePower" form:"reactivePower" point:"name:无功功率(kVar),value:27,type:2"`
    ApparentPower   *float64 `json:"apparentPower" form:"apparentPower" point:"name:视在功率(kVa),value:28,type:2"`
    GridFrequency   *float64 `json:"gridFrequency" form:"gridFrequency" point:"name:电网频率(Hz),value:15,type:2"`
    TotalDCCurrent  *float64 `json:"totalDCCurrent" form:"totalDCCurrent" point:"name:总直流电流(A),value:16,type:2"`
    PeakActivePower *float64 `json:"peakActivePower" form:"peakActivePower" point:"name:当天峰值有功功率(kW),value:503,type:2"`

    // 相电压
    PhaseAVoltage *float64 `json:"phaseAVoltage" form:"phaseAVoltage" point:"name:A相电压Ua(V),value:94,type:2"`
    PhaseBVoltage *float64 `json:"phaseBVoltage" form:"phaseBVoltage" point:"name:B相电压Ub(V),value:95,type:2"`
    PhaseCVoltage *float64 `json:"phaseCVoltage" form:"phaseCVoltage" point:"name:C相电压Uc(V),value:96,type:2"`

    // 线电压
    LineABVoltage *float64 `json:"lineABVoltage" form:"lineABVoltage" point:"name:AB线电压Uab(V),value:7,type:2"`
    LineBCVoltage *float64 `json:"lineBCVoltage" form:"lineBCVoltage" point:"name:BC线电压Ubc(V),value:8,type:2"`
    LineCAVoltage *float64 `json:"lineCAVoltage" form:"lineCAVoltage" point:"name:CA线电压Uca(V),value:9,type:2"`

    // 相电流
    PhaseACurrent *float64 `json:"phaseACurrent" form:"phaseACurrent" point:"name:A相电流Ia(A),value:12,type:2"`
    PhaseBCurrent *float64 `json:"phaseBCurrent" form:"phaseBCurrent" point:"name:B相电流Ib(A),value:13,type:2"`
    PhaseCCurrent *float64 `json:"phaseCCurrent" form:"phaseCCurrent" point:"name:C相电流Ic(A),value:14,type:2"`

    // 功率因数和温度
    PowerFactor          *float64 `json:"powerFactor" form:"powerFactor" point:"name:功率因数,value:21,type:2"`
    DeviceTemperature    *float64 `json:"deviceTemperature" form:"deviceTemperature" point:"name:设备温度(℃),value:22,type:2"`
    ConversionEfficiency *float64 `json:"conversionEfficiency" form:"conversionEfficiency" point:"name:转换效率(%),value:102,type:2"`

    // 支路电压 (20路)
    BranchVoltage1  *float64 `json:"branchVoltage1" form:"branchVoltage1" point:"name:支路电压1(V),value:301,type:2"`
    BranchVoltage2  *float64 `json:"branchVoltage2" form:"branchVoltage2" point:"name:支路电压2(V),value:302,type:2"`
    BranchVoltage3  *float64 `json:"branchVoltage3" form:"branchVoltage3" point:"name:支路电压3(V),value:303,type:2"`
    BranchVoltage4  *float64 `json:"branchVoltage4" form:"branchVoltage4" point:"name:支路电压4(V),value:304,type:2"`
    BranchVoltage5  *float64 `json:"branchVoltage5" form:"branchVoltage5" point:"name:支路电压5(V),value:305,type:2"`
    BranchVoltage6  *float64 `json:"branchVoltage6" form:"branchVoltage6" point:"name:支路电压6(V),value:306,type:2"`
    BranchVoltage7  *float64 `json:"branchVoltage7" form:"branchVoltage7" point:"name:支路电压7(V),value:307,type:2"`
    BranchVoltage8  *float64 `json:"branchVoltage8" form:"branchVoltage8" point:"name:支路电压8(V),value:308,type:2"`
    BranchVoltage9  *float64 `json:"branchVoltage9" form:"branchVoltage9" point:"name:支路电压9(V),value:309,type:2"`
    BranchVoltage10 *float64 `json:"branchVoltage10" form:"branchVoltage10" point:"name:支路电压10(V),value:310,type:2"`
    BranchVoltage11 *float64 `json:"branchVoltage11" form:"branchVoltage11" point:"name:支路电压11(V),value:311,type:2"`
    BranchVoltage12 *float64 `json:"branchVoltage12" form:"branchVoltage12" point:"name:支路电压12(V),value:312,type:2"`
    BranchVoltage13 *float64 `json:"branchVoltage13" form:"branchVoltage13" point:"name:支路电压13(V),value:313,type:2"`
    BranchVoltage14 *float64 `json:"branchVoltage14" form:"branchVoltage14" point:"name:支路电压14(V),value:314,type:2"`
    BranchVoltage15 *float64 `json:"branchVoltage15" form:"branchVoltage15" point:"name:支路电压15(V),value:315,type:2"`
    BranchVoltage16 *float64 `json:"branchVoltage16" form:"branchVoltage16" point:"name:支路电压16(V),value:316,type:2"`
    BranchVoltage17 *float64 `json:"branchVoltage17" form:"branchVoltage17" point:"name:支路电压17(V),value:317,type:2"`
    BranchVoltage18 *float64 `json:"branchVoltage18" form:"branchVoltage18" point:"name:支路电压18(V),value:318,type:2"`
    BranchVoltage19 *float64 `json:"branchVoltage19" form:"branchVoltage19" point:"name:支路电压19(V),value:319,type:2"`
    BranchVoltage20 *float64 `json:"branchVoltage20" form:"branchVoltage20" point:"name:支路电压20(V),value:320,type:2"`

    // 支路电流 (20路)
    BranchCurrent1  *float64 `json:"branchCurrent1" form:"branchCurrent1" point:"name:支路电流1(A),value:36,type:2"`
    BranchCurrent2  *float64 `json:"branchCurrent2" form:"branchCurrent2" point:"name:支路电流2(A),value:37,type:2"`
    BranchCurrent3  *float64 `json:"branchCurrent3" form:"branchCurrent3" point:"name:支路电流3(A),value:38,type:2"`
    BranchCurrent4  *float64 `json:"branchCurrent4" form:"branchCurrent4" point:"name:支路电流4(A),value:39,type:2"`
    BranchCurrent5  *float64 `json:"branchCurrent5" form:"branchCurrent5" point:"name:支路电流5(A),value:40,type:2"`
    BranchCurrent6  *float64 `json:"branchCurrent6" form:"branchCurrent6" point:"name:支路电流6(A),value:41,type:2"`
    BranchCurrent7  *float64 `json:"branchCurrent7" form:"branchCurrent7" point:"name:支路电流7(A),value:42,type:2"`
    BranchCurrent8  *float64 `json:"branchCurrent8" form:"branchCurrent8" point:"name:支路电流8(A),value:43,type:2"`
    BranchCurrent9  *float64 `json:"branchCurrent9" form:"branchCurrent9" point:"name:支路电流9(A),value:44,type:2"`
    BranchCurrent10 *float64 `json:"branchCurrent10" form:"branchCurrent10" point:"name:支路电流10(A),value:45,type:2"`
    BranchCurrent11 *float64 `json:"branchCurrent11" form:"branchCurrent11" point:"name:支路电流11(A),value:46,type:2"`
    BranchCurrent12 *float64 `json:"branchCurrent12" form:"branchCurrent12" point:"name:支路电流12(A),value:47,type:2"`
    BranchCurrent13 *float64 `json:"branchCurrent13" form:"branchCurrent13" point:"name:支路电流13(A),value:48,type:2"`
    BranchCurrent14 *float64 `json:"branchCurrent14" form:"branchCurrent14" point:"name:支路电流14(A),value:49,type:2"`
    BranchCurrent15 *float64 `json:"branchCurrent15" form:"branchCurrent15" point:"name:支路电流15(A),value:50,type:2"`
    BranchCurrent16 *float64 `json:"branchCurrent16" form:"branchCurrent16" point:"name:支路电流16(A),value:51,type:2"`
    BranchCurrent17 *float64 `json:"branchCurrent17" form:"branchCurrent17" point:"name:支路电流17(A),value:52,type:2"`
    BranchCurrent18 *float64 `json:"branchCurrent18" form:"branchCurrent18" point:"name:支路电流18(A),value:53,type:2"`
    BranchCurrent19 *float64 `json:"branchCurrent19" form:"branchCurrent19" point:"name:支路电流19(A),value:54,type:2"`
    BranchCurrent20 *float64 `json:"branchCurrent20" form:"branchCurrent20" point:"name:支路电流20(A),value:55,type:2"`

    // MPPT
    MPPT1 *float64 `json:"mppt1" form:"mppt1" point:"name:mppt1,value:401,type:2"`
    MPPT2 *float64 `json:"mppt2" form:"mppt2" point:"name:mppt2,value:402,type:2"`
    MPPT3 *float64 `json:"mppt3" form:"mppt3" point:"name:mppt3,value:403,type:2"`
    MPPT4 *float64 `json:"mppt4" form:"mppt4" point:"name:mppt4,value:404,type:2"`

    // 设备状态和其他
    DeviceStatusCode    *float64 `json:"deviceStatusCode" form:"deviceStatusCode" point:"name:设备状态码,value:501,type:1"`
    InsulationImpedance *float64 `json:"insulationImpedance" form:"insulationImpedance" point:"name:绝缘阻抗,value:502,type:2"`
    StartupTime         *float64 `json:"startupTime" form:"startupTime" point:"name:开机时间,value:511,type:2"`
    ShutdownTime        *float64 `json:"shutdownTime" form:"shutdownTime" point:"name:关机时间,value:512,type:2"`
    PhaseATemperature   *float64 `json:"phaseATemperature" form:"phaseATemperature" point:"name:A相温度,value:251,type:2"`
    PhaseBTemperature   *float64 `json:"phaseBTemperature" form:"phaseBTemperature" point:"name:B相温度,value:252,type:2"`
    PhaseCTemperature   *float64 `json:"phaseCTemperature" form:"phaseCTemperature" point:"name:C相温度,value:253,type:2"`
}

// TableName agvcNbqHis表 AgvcNbqHis自定义表名 agvc_nbq_his
func (AgvcNbq) TableName() string {
    return "agvc_nbq_his"
}

// PointInfo point信息结构体
type PointInfo struct {
    FieldName string
    DataType  int // 1=遥信(YX), 2=遥测(YC)
}

// PointToFieldMapping point值到字段信息的映射缓存
var pointToFieldMapping map[string]PointInfo
var pointMappingOnce sync.Once

// InitPointMapping 初始化point标签映射，在系统启动时调用
func InitPointMapping() {
    pointMappingOnce.Do(func() {
        pointToFieldMapping = make(map[string]PointInfo)

        // 使用反射解析AgvcNbqHis结构体的point标签
        t := reflect.TypeOf(AgvcNbq{})
        for i := 0; i < t.NumField(); i++ {
            field := t.Field(i)
            tag := field.Tag.Get("point")

            if tag == "" {
                continue
            }

            // 解析point标签，提取value和type值
            if strings.Contains(tag, "value:") {
                parts := strings.Split(tag, ",")
                var pointValue string
                var dataType int = 2 // 默认为遥测

                for _, part := range parts {
                    part = strings.TrimSpace(part)
                    if strings.HasPrefix(part, "value:") {
                        pointValue = strings.TrimPrefix(part, "value:")
                    } else if strings.HasPrefix(part, "type:") {
                        typeStr := strings.TrimPrefix(part, "type:")
                        if typeStr == "1" {
                            dataType = 1
                        } else if typeStr == "2" {
                            dataType = 2
                        }
                    }
                }

                if pointValue != "" {
                    pointToFieldMapping[pointValue] = PointInfo{
                        FieldName: field.Name,
                        DataType:  dataType,
                    }
                }
            }
        }
    })
}

// GetFieldNameByPoint 根据point值获取字段名
func GetFieldNameByPoint(point string) (string, bool) {
    pointInfo, ok := pointToFieldMapping[point]
    return pointInfo.FieldName, ok
}

// GetPointInfoByPoint 根据point值获取点位信息
func GetPointInfoByPoint(point string) (PointInfo, bool) {
    pointInfo, ok := pointToFieldMapping[point]
    return pointInfo, ok
}

// GetAllPointValues 获取所有point值列表
func GetAllPointValues() []string {
    points := make([]string, 0, len(pointToFieldMapping))
    for point := range pointToFieldMapping {
        points = append(points, point)
    }
    return points
}

// GetPointValuesByType 根据数据类型获取point值列表
func GetPointValuesByType(dataType int) []string {
    points := make([]string, 0)
    for point, info := range pointToFieldMapping {
        if info.DataType == dataType {
            points = append(points, point)
        }
    }
    return points
}
