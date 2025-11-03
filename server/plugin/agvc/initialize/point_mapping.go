package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/model"
	"go.uber.org/zap"
)

// InitPointMappings 初始化测点映射数据
func InitPointMappings(ctx context.Context) {
	// 检查是否已经初始化过
	var count int64
	global.GVA_DB.Model(&model.PointMapping{}).Count(&count)
	if count > 0 {
		global.GVA_LOG.Info("测点映射已存在，跳过初始化")
		return
	}

	mappings := []model.PointMapping{
		// AGC遥信点
		{EQType: "02", DataType: "01", Point: "401", PointName: "AGC投退信号", Unit: "", Description: "0:未投入 1:投入", Category: "AGC"},
		{EQType: "02", DataType: "01", Point: "402", PointName: "AGC就地远方控制模式", Unit: "", Description: "1:调度 2:站内", Category: "AGC"},
		{EQType: "02", DataType: "01", Point: "404", PointName: "AGC开/闭环状态", Unit: "", Description: "1:开环 2:闭环", Category: "AGC"},
		{EQType: "02", DataType: "01", Point: "405", PointName: "AGC有功上调节闭锁", Unit: "", Description: "0:未闭锁 1:闭锁", Category: "AGC"},
		{EQType: "02", DataType: "01", Point: "406", PointName: "AGC有功下调节闭锁", Unit: "", Description: "0:未闭锁 1:闭锁", Category: "AGC"},

		// AGC遥测点
		{EQType: "02", DataType: "02", Point: "401", PointName: "有功调节上限", Unit: "MW", Description: "有功调节上限", Category: "AGC"},
		{EQType: "02", DataType: "02", Point: "402", PointName: "有功调节下限", Unit: "MW", Description: "有功调节下限", Category: "AGC"},
		{EQType: "02", DataType: "02", Point: "403", PointName: "有功执行值", Unit: "MW", Description: "有功执行值", Category: "AGC"},

		// AGC遥控点
		{EQType: "02", DataType: "03", Point: "401", PointName: "AGC投退信号", Unit: "", Description: "控制AGC投退", Category: "AGC"},
		{EQType: "02", DataType: "03", Point: "402", PointName: "AGC就地远方控制模式", Unit: "", Description: "切换控制模式", Category: "AGC"},
		{EQType: "02", DataType: "03", Point: "404", PointName: "AGC开/闭环状态", Unit: "", Description: "切换开闭环", Category: "AGC"},

		// AGC遥调点
		{EQType: "02", DataType: "04", Point: "401", PointName: "有功执行", Unit: "MW", Description: "有功执行值下发", Category: "AGC"},

		// AVC遥信点
		{EQType: "02", DataType: "01", Point: "411", PointName: "AVC功能投退信号", Unit: "", Description: "0:未投入 1:投入", Category: "AVC"},
		{EQType: "02", DataType: "01", Point: "412", PointName: "AVC功能就地远方控制模式", Unit: "", Description: "1:调度 2:站内", Category: "AVC"},
		{EQType: "02", DataType: "01", Point: "413", PointName: "AVC功能当前指令状态", Unit: "", Description: "0:无指令 1:有指令", Category: "AVC"},
		{EQType: "02", DataType: "01", Point: "414", PointName: "AVC功能开闭环状态", Unit: "", Description: "1:开环 2:闭环", Category: "AVC"},
		{EQType: "02", DataType: "01", Point: "415", PointName: "AVC功能上调节闭锁", Unit: "", Description: "0:未闭锁 1:闭锁", Category: "AVC"},
		{EQType: "02", DataType: "01", Point: "416", PointName: "AVC功能下调节闭锁", Unit: "", Description: "0:未闭锁 1:闭锁", Category: "AVC"},

		// AVC遥测点
		{EQType: "02", DataType: "02", Point: "411", PointName: "无功可增容量", Unit: "MVar", Description: "无功可增容量", Category: "AVC"},
		{EQType: "02", DataType: "02", Point: "412", PointName: "无功可减容量", Unit: "MVar", Description: "无功可减容量", Category: "AVC"},
		{EQType: "02", DataType: "02", Point: "413", PointName: "电压执行值", Unit: "kV", Description: "电压执行值", Category: "AVC"},
		{EQType: "02", DataType: "02", Point: "414", PointName: "无功执行值", Unit: "MVar", Description: "无功执行值", Category: "AVC"},

		// AVC遥控点
		{EQType: "02", DataType: "03", Point: "411", PointName: "AVC功能投退信号", Unit: "", Description: "控制AVC投退", Category: "AVC"},
		{EQType: "02", DataType: "03", Point: "412", PointName: "AVC功能就地远方控制模式", Unit: "", Description: "切换控制模式", Category: "AVC"},
		{EQType: "02", DataType: "03", Point: "413", PointName: "AVC功能当前指令状态", Unit: "", Description: "指令状态", Category: "AVC"},
		{EQType: "02", DataType: "03", Point: "414", PointName: "AVC功能开闭环状态", Unit: "", Description: "切换开闭环", Category: "AVC"},

		// AVC遥调点
		{EQType: "02", DataType: "04", Point: "411", PointName: "电压执行", Unit: "kV", Description: "电压执行值下发", Category: "AVC"},
		{EQType: "02", DataType: "04", Point: "412", PointName: "无功执行", Unit: "MVar", Description: "无功执行值下发", Category: "AVC"},

		// 逆变器遥测点
		{EQType: "01", DataType: "02", Point: "401", PointName: "有功功率", Unit: "MW", Description: "逆变器有功功率", Category: "INVERTER"},
		{EQType: "01", DataType: "02", Point: "402", PointName: "无功功率", Unit: "MVar", Description: "逆变器无功功率", Category: "INVERTER"},
		{EQType: "01", DataType: "02", Point: "403", PointName: "运行状态", Unit: "", Description: "0:停机 1:运行", Category: "INVERTER"},

		// 逆变器遥控点
		{EQType: "01", DataType: "03", Point: "401", PointName: "开关机", Unit: "", Description: "0:停机 1:开机", Category: "INVERTER"},

		// 逆变器遥调点
		{EQType: "01", DataType: "04", Point: "401", PointName: "有功功率降额执行值", Unit: "MW", Description: "有功功率设定值", Category: "INVERTER"},
		{EQType: "01", DataType: "04", Point: "402", PointName: "无功功率补偿执行值", Unit: "MVar", Description: "无功功率设定值", Category: "INVERTER"},

		// 并网点遥测点（用于采集）
		{EQType: "02", DataType: "02", Point: "401", PointName: "并网点电压", Unit: "kV", Description: "并网点电压", Category: "POINT"},
		{EQType: "02", DataType: "02", Point: "402", PointName: "并网点有功", Unit: "MW", Description: "并网点有功功率", Category: "POINT"},
		{EQType: "02", DataType: "02", Point: "403", PointName: "并网点无功", Unit: "MVar", Description: "并网点无功功率", Category: "POINT"},
		{EQType: "02", DataType: "02", Point: "404", PointName: "系统频率", Unit: "Hz", Description: "系统频率", Category: "POINT"},
	}

	if err := global.GVA_DB.CreateInBatches(mappings, 100).Error; err != nil {
		global.GVA_LOG.Error("初始化测点映射失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("测点映射初始化成功", zap.Int("数量", len(mappings)))
	}
}
