package agcv_main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

type coapSender struct{}

var CoapSender = new(coapSender)

// SendData 发送数据到CoAP服务器
func (s *coapSender) SendData(host string, port int, messages []request.CoAPDataMessage) error {
	if len(messages) == 0 {
		return fmt.Errorf("没有数据需要发送")
	}

	// 创建CoAP客户端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		global.GVA_LOG.Error("连接CoAP服务器失败",
			zap.String("host", host),
			zap.Int("port", port),
			zap.Error(err))
		return fmt.Errorf("连接CoAP服务器失败: %v", err)
	}
	defer conn.Close()

	// coapClient := client.NewClient(conn)

	// 准备请求数据
	jsonData, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		global.GVA_LOG.Error("发送CoAP数据失败",
			zap.String("host", host),
			zap.Int("port", port),
			zap.Error(err))
		return fmt.Errorf("发送CoAP数据失败: %v", err)
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		global.GVA_LOG.Error("接收CoAP响应失败",
			zap.String("host", host),
			zap.Int("port", port),
			zap.Error(err))
		return fmt.Errorf("接收CoAP响应失败: %v", err)
	}
	res := string(buf[:n])
	fmt.Println(res)

	return nil
}

// SendAGCCommand 发送AGC控制指令
func (s *coapSender) SendAGCCommand(host string, port int, bwdNo int, commands map[int]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)

	// 构建AGC标准点数据
	// 遥控点：AGC投退信号(401), AGC就地远方控制模式(402), AGC开/闭环状态(404)
	// 遥调点：有功执行(401)

	for point, value := range commands {
		var dataType int
		var pointStr string
		switch point {
		case 401, 402, 404: // 遥控
			dataType = cons.YK
			pointStr = fmt.Sprintf("%d", point)
		case 403: // 有功执行（遥调）
			dataType = cons.YT
			pointStr = "401"
		default:
			continue
		}

		messages = append(messages, request.CoAPDataMessage{
			PSID:     1,
			EQID:     bwdNo,         // AGC是全站级别的，设备ID为0000
			EQType:   cons.TYPE_BWG, // 并网点
			DataType: dataType,
			Point:    pointStr,
			Value:    value,
		})
	}

	if len(messages) == 0 {
		return fmt.Errorf("没有有效的AGC指令")
	}

	return s.SendData(host, port, messages)
}

// SendAVCCommand 发送AVC控制指令
func (s *coapSender) SendAVCCommand(host string, port, bwdNo int, commands map[int]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)

	// 构建AVC标准点数据
	// 遥控点：AVC功能投退信号(401), AVC功能就地远方控制模式(402), AVC功能当前指令状态(403), AVC功能开闭环状态(404)
	// 遥调点：电压执行(401), 无功执行(402)

	for point, value := range commands {
		var dataType int
		var pointStr string
		switch point {
		case 401, 402, 403: // 遥控
			dataType = cons.YK
			pointStr = fmt.Sprintf("%d", point)
		case 404, 405: // 遥调
			dataType = cons.YT
			if point == 406 {
				pointStr = "401"
			} else {
				pointStr = "402"
			}
		default:
			continue
		}

		messages = append(messages, request.CoAPDataMessage{
			PSID:     1,
			EQID:     bwdNo,         // AVC是全站级别的，设备ID为0000
			EQType:   cons.TYPE_BWG, // 并网点
			DataType: dataType,
			Point:    pointStr,
			Value:    value,
		})
	}

	if len(messages) == 0 {
		return fmt.Errorf("没有有效的AVC指令")
	}

	return s.SendData(host, port, messages)
}

// SendInverterCommand 发送逆变器控制指令
func (s *coapSender) SendInverterCommand(host string, port, psid, eqid int, commands map[int]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)

	// 逆变器标准点
	// 遥控：开关机(401)
	// 遥调：有功功率降额执行值(401), 无功功率补偿执行值(402)

	for point, value := range commands {
		var dataType int
		var pointStr string
		switch point {
		case 501: // 有功功率固定值降额(kW)
			dataType = cons.YT
			pointStr = "501"
		case 502: // 无功功率补偿(PF)
			dataType = cons.YT
			pointStr = "502"
		case 503: // 无功功率补偿(Q/S)
			dataType = cons.YT
			pointStr = "503"
		case 504: // 有功功率百分比降额(0.1%)
			dataType = cons.YT
			pointStr = "504"
		case 505: // 有功功率固定值降额(W)
			dataType = cons.YT
			pointStr = "505"
		default:
			continue
		}

		messages = append(messages, request.CoAPDataMessage{
			PSID:     psid,
			EQID:     eqid,
			EQType:   cons.TYPE_NBQ, // 逆变器
			DataType: dataType,
			Point:    pointStr,
			Value:    value,
		})
	}

	if len(messages) == 0 {
		return fmt.Errorf("没有有效的逆变器指令")
	}

	return s.SendData(host, port, messages)
}

// GetDefaultCoapHost 获取默认CoAP主机地址（可从配置读取）
func (s *coapSender) GetDefaultCoapHost() string {
	// TODO: 从配置文件读取
	return "127.0.0.1"
}

// GetDefaultCoapPort 获取默认CoAP端口
func (s *coapSender) GetDefaultCoapPort() int {
	return 1188
}

// GetDispatchCoapHost 获取调度CoAP主机地址
func (s *coapSender) GetDispatchCoapHost() string {
	// TODO: 从配置文件读取
	return "127.0.0.1"
}

// GetDispatchCoapPort 获取调度CoAP返回端口（1189）
func (s *coapSender) GetDispatchCoapPort() int {
	return 1189
}

func (s *coapSender) GetDispatchBackCoapPort() int {
	return 1187
}

func (s *coapSender) SavaDataToMemAndDB(bwdNo int, messages []request.CoAPDataMessage) error {
	//获取并网点最新配置
	config, err := AGC.GetAGCConfig(bwdNo)
	if err != nil {
		global.GVA_LOG.Error("获取AGC配置失败", zap.Error(err))
		return err
	}

	for _, message := range messages {
		DataStorage.StoreData(&agvc_main.RealtimeData{
			PSID:     message.PSID,
			EQID:     message.EQID,
			EQType:   message.EQType,
			DataType: message.DataType,
			Point:    message.Point,
			Value:    message.Value,
		})
		if message.EQType == cons.TYPE_AGC {
			if message.DataType == cons.YX {
				switch message.Point {
				case "401":
					val := int64(message.Value.(float64))
					config.AgcIsEnabled = &val
				case "402":
					val := int64(message.Value.(float64))
					config.AgcRemoteMode = &val
				case "404":
					val := int64(message.Value.(float64))
					config.AgcLoopStatus = &val
				}
			}
			if message.DataType == cons.YC {
				switch message.Point {
				case "403":
					val := message.Value.(float64)
					config.AgcDispatchExecValue = &val
				}

			}
		}
		if message.EQType == cons.TYPE_AVC {
			if message.DataType == cons.YX {
				switch message.Point {
				case "401":
					val := int64(message.Value.(float64))
					config.AvcIsEnabled = &val
				case "402":
					val := int64(message.Value.(float64))
					config.AvcRemoteMode = &val
				case "404":
					val := int64(message.Value.(float64))
					config.AvcLoopStatus = &val
				}
			}
			if message.DataType == cons.YC {
				switch message.Point {
				case "403":
					val := message.Value.(float64)
					config.AvcVolteExecValue = &val
				case "404":
					val := message.Value.(float64)
					config.AvcWGExecValue = &val
				}
			}
		}
	}
	//更新数据库
	global.GVA_DB.Save(&config)
	return nil
}

// SendAGCResultToDispatch 发送AGC计算结果到调度（1189端口）
func (s *coapSender) SendAGCResultToDispatch(bwdNo int, results map[string]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)

	//发送给调度后更新内存和数据库
	defer func() {
		if err := s.SavaDataToMemAndDB(bwdNo, messages); err != nil {
			global.GVA_LOG.Error("保存AGC数据到内存和数据库失败", zap.Error(err))
		}
	}()

	// AGC遥信标准点
	yxPoints := map[string]string{
		"401": "agcSignal",      // AGC投退信号
		"402": "agcControlMode", // AGC就地远方控制模式
		"404": "agcLoopStatus",  // AGC开/闭环状态
		"405": "agcUpRegLock",   // AGC有功上调节闭锁
		"406": "agcDownRegLock", // AGC有功下调节闭锁
	}

	// AGC遥测标准点
	ycPoints := map[string]string{
		"401": "powerUpperLimit", // 有功调节上限
		"402": "powerLowerLimit", // 有功调节下限
		"403": "powerExecValue",  // 有功执行值
	}

	// 添加遥信点
	for point, key := range yxPoints {
		if value, ok := results[key]; ok {
			messages = append(messages, request.CoAPDataMessage{
				PSID:     1,
				EQID:     bwdNo,
				EQType:   cons.TYPE_AGC,
				DataType: cons.YX,
				Point:    point,
				Value:    value,
			})
		}
	}

	// 添加遥测点
	for point, key := range ycPoints {
		if value, ok := results[key]; ok {
			messages = append(messages, request.CoAPDataMessage{
				PSID:     1,
				EQID:     bwdNo,
				EQType:   cons.TYPE_AGC,
				DataType: cons.YC,
				Point:    point,
				Value:    value,
			})
		}
	}

	if len(messages) == 0 {
		return fmt.Errorf("没有AGC结果数据需要发送")
	}

	host := s.GetDispatchCoapHost()
	port := s.GetDispatchCoapPort()

	global.GVA_LOG.Info("发送AGC计算结果到调度",
		zap.Int("bwdNo", bwdNo),
		zap.String("host", host),
		zap.Int("port", port),
		zap.Int("dataCount", len(messages)))

	return s.SendData(host, port, messages)
}

// SendAVCResultToDispatch 发送AVC计算结果到调度（1189端口）
func (s *coapSender) SendAVCResultToDispatch(bwdNo int, results map[string]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)

	//发送给调度后更新内存和数据库
	defer func() {
		if err := s.SavaDataToMemAndDB(bwdNo, messages); err != nil {
			global.GVA_LOG.Error("保存AGC数据到内存和数据库失败", zap.Error(err))
		}
	}()

	// AVC遥信标准点
	yxPoints := map[string]string{
		"401": "avcSignal",      // AVC功能投退信号
		"402": "avcControlMode", // AVC功能就地远方控制模式
		"403": "avcCmdStatus",   // AVC功能当前指令状态
		"404": "avcLoopStatus",  // AVC功能开闭环状态
		"405": "avcUpRegLock",   // AVC功能上调节闭锁
		"406": "avcDownRegLock", // AVC功能下调节闭锁
	}

	// AVC遥测标准点
	ycPoints := map[string]string{
		// "401": "reactiveIncreaseCap", // 无功可增容量
		// "402": "reactiveDecreaseCap", // 无功可减容量
		"403": "voltageExecValue",  // 电压执行值
		"404": "reactiveExecValue", // 无功执行值
	}

	// 添加遥信点
	for point, key := range yxPoints {
		if value, ok := results[key]; ok {
			messages = append(messages, request.CoAPDataMessage{
				PSID:     1,
				EQID:     bwdNo,
				EQType:   cons.TYPE_AVC,
				DataType: cons.YX,
				Point:    point,
				Value:    value,
			})
		}
	}

	// 添加遥测点
	for point, key := range ycPoints {
		if value, ok := results[key]; ok {
			messages = append(messages, request.CoAPDataMessage{
				PSID:     1,
				EQID:     bwdNo,
				EQType:   cons.TYPE_AVC,
				DataType: cons.YC,
				Point:    point,
				Value:    value,
			})
		}
	}

	if len(messages) == 0 {
		return fmt.Errorf("没有AVC结果数据需要发送")
	}

	host := s.GetDispatchCoapHost()
	port := s.GetDispatchCoapPort()

	global.GVA_LOG.Info("发送AVC计算结果到调度",
		zap.Int("bwdNo", bwdNo),
		zap.String("host", host),
		zap.Int("port", port),
		zap.Int("dataCount", len(messages)))

	return s.SendData(host, port, messages)
}
