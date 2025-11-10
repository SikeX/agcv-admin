package cons

// 数据类型
const (
	YX = 1
	YC = 2
	YM = 3
	YK = 4
	YT = 5
)

const (
	TYPE_BWG = 5
	TYPE_NBQ = 2
	TYPE_AGC = 60
	TYPE_AVC = 61
)

// AGC调度标准点名称常量
const (
	// 遥信(YX)
	AGC_YX_SIGNAL        = "AGC投退信号"     // 点标识: 401
	AGC_YX_CONTROL_MODE  = "AGC就地远方控制模式" // 点标识: 402
	AGC_YX_LOOP_STATUS   = "AGC开/闭环状态"   // 点标识: 404
	AGC_YX_UP_REG_LOCK   = "AGC有功上调节闭锁"  // 点标识: 405
	AGC_YX_DOWN_REG_LOCK = "AGC有功下调节闭锁"  // 点标识: 406

	// 遥测(YC)
	AGC_YC_POWER_UPPER_LIMIT = "有功调节上限" // 点标识: 401
	AGC_YC_POWER_LOWER_LIMIT = "有功调节下限" // 点标识: 402
	AGC_YC_POWER_EXEC_VALUE  = "有功执行值"  // 点标识: 403

	// 遥控(YK)
	AGC_YK_SIGNAL       = "AGC投退信号"     // 点标识: 401
	AGC_YK_CONTROL_MODE = "AGC就地远方控制模式" // 点标识: 402
	AGC_YK_LOOP_STATUS  = "AGC开/闭环状态"   // 点标识: 404

	// 遥调(YT)
	AGC_YT_POWER_EXEC = "有功执行" // 点标识: 401
)

// AVC调度标准点名称常量
const (
	// 遥信(YX)
	AVC_YX_SIGNAL         = "AVC功能投退信号"     // 点标识: 401
	AVC_YX_CONTROL_MODE   = "AVC功能就地远方控制模式" // 点标识: 402
	AVC_YX_COMMAND_STATUS = "AVC功能当前指令状态"   // 点标识: 403
	AVC_YX_LOOP_STATUS    = "AVC功能开闭环状态"    // 点标识: 404
	AVC_YX_UP_REG_LOCK    = "AVC功能上调节闭锁"    // 点标识: 405
	AVC_YX_DOWN_REG_LOCK  = "AVC功能下调节闭锁"    // 点标识: 406

	// 遥测(YC)
	AVC_YC_REACTIVE_INC_CAP    = "无功可增容量" // 点标识: 401
	AVC_YC_REACTIVE_DEC_CAP    = "无功可减容量" // 点标识: 402
	AVC_YC_VOLTAGE_EXEC_VALUE  = "电压执行值"  // 点标识: 401
	AVC_YC_REACTIVE_EXEC_VALUE = "无功执行值"  // 点标识: 402

	// 遥控(YK)
	AVC_YK_SIGNAL         = "AVC功能投退信号"     // 点标识: 401
	AVC_YK_CONTROL_MODE   = "AVC功能就地远方控制模式" // 点标识: 402
	AVC_YK_COMMAND_STATUS = "AVC功能当前指令状态"   // 点标识: 403
	AVC_YK_LOOP_STATUS    = "AVC功能开闭环状态"    // 点标识: 404

	// 遥调(YT)
	AVC_YT_VOLTAGE_EXEC  = "电压执行" // 点标识: 401
	AVC_YT_REACTIVE_EXEC = "无功执行" // 点标识: 402
)

// 并网柜常用点名称常量
const (
	BWG_YC_ACTIVE_POWER   = "有功功率P(kW)"
	BWG_YC_REACTIVE_POWER = "无功功率Q(kvar)"
	BWG_YC_VOLTAGE_AB     = "AB线电压Uab(kV)"
	BWG_YC_FREQUENCY      = "F(频率Hz)"
	BWG_YC_APPARENT_POWER = "S(视在功率kVA)"
)

// 逆变器常用点名称常量
const (
	NBQ_YC_ACTIVE_POWER   = "交流功率(kW)"   // 点标识: 10
	NBQ_YC_REACTIVE_POWER = "无功功率(kVar)" // 点标识:
	NBQ_YC_APPARENT_POWER = "视在功率(kVa)"  // 点标识:
)

// 逆变器品牌常量
const (
	INVERTER_BRAND_HUAWEI  = 1 // 华为
	INVERTER_BRAND_SUNGROW = 2 // 阳光
	INVERTER_BRAND_GOODWE  = 3 // 固德威
)

// 逆变器控制点位常量
const (
	// 遥控点位
	INV_YK_POWER_SWITCH = "开关机" // 点标识: 401

	// 遥调点位
	INV_YT_ACTIVE_POWER_LIMIT  = "有功功率降额执行值" // 点标识: 401
	INV_YT_REACTIVE_POWER_COMP = "无功功率补偿执行值" // 点标识: 402
)

const (
	// 遥测点位
	HUAWEI_YG_MODE             = "501" //有功功率执行值
	HUAWEI_POWER_EXEC_VALUE    = "502" //有功功率执行值
	HUAWEI_POWER_ADJUST_VALUE  = "503" //有功调节指令值
	HUAWEI_WG_MODE             = "504" //无功调节模式
	HUAWEI_REACTIVE_EXEC_VALUE = "505" //无功执行值

)
