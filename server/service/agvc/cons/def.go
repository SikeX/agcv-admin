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
	TYPE_AGC = 98
	TYPE_AVC = 99
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
	AGC_YC_POWER_EXEC_VALUE  = "有功执行值"  // 点标识: 401

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
	NBQ_YC_ACTIVE_POWER   = "交流功率(kW)"
	NBQ_YC_REACTIVE_POWER = "无功功率(kVar)"
	NBQ_YC_APPARENT_POWER = "视在功率(kVa)"
)
