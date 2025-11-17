package agvc

const (
	// AGC设备类型
	AGC_DEVICE_TYPE = 60
	// AVC设备类型
	AVC_DEVICE_TYPE = 61
	// BWD设备类型
	BWD_DEVICE_TYPE = 5
)

// AGC调度标准点名称常量
const (
	// 遥信(YX)
	AGC_YX_SIGNAL        = "401" // AGC投退信号 点标识: 401
	AGC_YX_CONTROL_MODE  = "402" // AGC就地远方控制模式 点标识: 402
	AGC_YX_LOOP_STATUS   = "404" // AGC开/闭环状态 点标识: 404
	AGC_YX_UP_REG_LOCK   = "405" // AGC有功上调节闭锁 点标识: 405
	AGC_YX_DOWN_REG_LOCK = "406" // AGC有功下调节闭锁 点标识: 406

	// 遥测(YC)
	AGC_YC_POWER_UPPER_LIMIT = "401" // 有功调节上限 点标识: 401
	AGC_YC_POWER_LOWER_LIMIT = "402" // 有功调节下限 点标识: 402
	AGC_YC_POWER_EXEC_VALUE  = "403" // 有功执行值 点标识: 403

	// 遥控(YK)
	AGC_YK_SIGNAL       = "401" // AGC投退信号 点标识: 401
	AGC_YK_CONTROL_MODE = "402" // AGC就地远方控制模式 点标识: 402
	AGC_YK_LOOP_STATUS  = "404" // AGC开/闭环状态 点标识: 404

	// 遥调(YT)
	AGC_YT_POWER_EXEC = "403" // 有功执行 点标识: 403
)

// AVC调度标准点名称常量
const (
	// 遥信(YX)
	AVC_YX_SIGNAL         = "401" // AVC功能投退信号 点标识: 401
	AVC_YX_CONTROL_MODE   = "402" // AVC功能就地远方控制模式 点标识: 402
	AVC_YX_COMMAND_STATUS = "403" // AVC功能当前指令状态 点标识: 403
	AVC_YX_LOOP_STATUS    = "404" // AVC功能开闭环状态 点标识: 404
	AVC_YX_UP_REG_LOCK    = "405" // AVC功能上调节闭锁 点标识: 405
	AVC_YX_DOWN_REG_LOCK  = "406" // AVC功能下调节闭锁 点标识: 406

	// 遥测(YC)
	AVC_YC_REACTIVE_INC_CAP    = "401" // 无功可增容量 点标识: 401
	AVC_YC_REACTIVE_DEC_CAP    = "402" // 无功可减容量 点标识: 402
	AVC_YC_VOLTAGE_EXEC_VALUE  = "403" // 电压执行值 点标识: 403
	AVC_YC_REACTIVE_EXEC_VALUE = "404" // 无功执行值 点标识: 404

	// 遥控(YK)
	AVC_YK_SIGNAL         = "401" // AVC功能投退信号 点标识: 401
	AVC_YK_CONTROL_MODE   = "402" // AVC功能就地远方控制模式 点标识: 402
	AVC_YK_COMMAND_STATUS = "403" // AVC功能当前指令状态 点标识: 403
	AVC_YK_LOOP_STATUS    = "404" // AVC功能开闭环状态 点标识: 404

	// 遥调(YT)
	AVC_YT_VOLTAGE_EXEC  = "403" // 电压执行 点标识: 401
	AVC_YT_REACTIVE_EXEC = "404" // 无功执行 点标识: 402
)

// bwd
const (
	BWD_YC_POWER_REAL = "7" // 有功实际值 点标识: 403
)
