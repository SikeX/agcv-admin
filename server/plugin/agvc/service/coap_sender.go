package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/model/request"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/plgd-dev/go-coap/v3/udp/client"
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
	conn, err := udp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		global.GVA_LOG.Error("连接CoAP服务器失败", 
			zap.String("host", host),
			zap.Int("port", port),
			zap.Error(err))
		return fmt.Errorf("连接CoAP服务器失败: %v", err)
	}
	defer conn.Close()

	coapClient := client.NewClient(conn)

	// 准备请求数据
	jsonData, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	// 创建CoAP请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送POST请求
	resp, err := coapClient.Post(ctx, "/agvc/data", message.AppJSON, bytes.NewReader(jsonData))
	if err != nil {
		global.GVA_LOG.Error("发送CoAP数据失败", 
			zap.String("host", host),
			zap.Int("port", port),
			zap.Error(err))
		return fmt.Errorf("发送CoAP数据失败: %v", err)
	}

	// 检查响应
	if resp.Code() != message.Content && resp.Code() != message.Created && resp.Code() != message.Changed {
		respBody, _ := resp.ReadBody()
		global.GVA_LOG.Warn("CoAP服务器返回非成功状态", 
			zap.String("code", resp.Code().String()),
			zap.String("body", string(respBody)))
		return fmt.Errorf("CoAP服务器返回错误: %s", resp.Code().String())
	}

	global.GVA_LOG.Info("CoAP数据发送成功", 
		zap.String("host", host),
		zap.Int("port", port),
		zap.Int("count", len(messages)))

	return nil
}

// SendAGCCommand 发送AGC控制指令
func (s *coapSender) SendAGCCommand(host string, port int, psid string, commands map[string]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)
	
	// 构建AGC标准点数据
	// 遥控点：AGC投退信号(401), AGC就地远方控制模式(402), AGC开/闭环状态(404)
	// 遥调点：有功执行(401)
	
	for point, value := range commands {
		var dataType string
		switch point {
		case "401", "402", "404": // 遥控
			dataType = "03"
		case "401_exec": // 有功执行（遥调）
			dataType = "04"
			point = "401"
		default:
			continue
		}
		
		messages = append(messages, request.CoAPDataMessage{
			PSID:     psid,
			EQID:     "0000", // AGC是全站级别的，设备ID为0000
			EQType:   "02",   // 并网点
			DataType: dataType,
			Point:    point,
			Value:    value,
		})
	}
	
	if len(messages) == 0 {
		return fmt.Errorf("没有有效的AGC指令")
	}
	
	return s.SendData(host, port, messages)
}

// SendAVCCommand 发送AVC控制指令
func (s *coapSender) SendAVCCommand(host string, port int, psid string, commands map[string]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)
	
	// 构建AVC标准点数据
	// 遥控点：AVC功能投退信号(401), AVC功能就地远方控制模式(402), AVC功能当前指令状态(403), AVC功能开闭环状态(404)
	// 遥调点：电压执行(401), 无功执行(402)
	
	for point, value := range commands {
		var dataType string
		switch point {
		case "401", "402", "403", "404": // 遥控
			dataType = "03"
		case "401_voltage", "402_reactive": // 遥调
			dataType = "04"
			if point == "401_voltage" {
				point = "401"
			} else {
				point = "402"
			}
		default:
			continue
		}
		
		messages = append(messages, request.CoAPDataMessage{
			PSID:     psid,
			EQID:     "0000", // AVC是全站级别的，设备ID为0000
			EQType:   "02",   // 并网点
			DataType: dataType,
			Point:    point,
			Value:    value,
		})
	}
	
	if len(messages) == 0 {
		return fmt.Errorf("没有有效的AVC指令")
	}
	
	return s.SendData(host, port, messages)
}

// SendInverterCommand 发送逆变器控制指令
func (s *coapSender) SendInverterCommand(host string, port int, psid, eqid string, commands map[string]interface{}) error {
	messages := make([]request.CoAPDataMessage, 0)
	
	// 逆变器标准点
	// 遥控：开关机(401)
	// 遥调：有功功率降额执行值(401), 无功功率补偿执行值(402)
	
	for point, value := range commands {
		var dataType string
		switch point {
		case "401_switch": // 开关机（遥控）
			dataType = "03"
			point = "401"
		case "401_power", "402_reactive": // 遥调
			dataType = "04"
			if point == "401_power" {
				point = "401"
			} else {
				point = "402"
			}
		default:
			continue
		}
		
		messages = append(messages, request.CoAPDataMessage{
			PSID:     psid,
			EQID:     eqid,
			EQType:   "01", // 逆变器
			DataType: dataType,
			Point:    point,
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
