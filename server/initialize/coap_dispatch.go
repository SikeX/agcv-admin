package initialize

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	agvcMainService "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
	"go.uber.org/zap"
)

// CoapDispatchServer 启动调度CoAP服务器（1187端口）
func CoapDispatchServer() {
	host := "0.0.0.0"
	port := 1187

	address := fmt.Sprintf("%s:%d", host, port)

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Error("解析调度CoAP UDP地址失败", zap.Error(err))
		}
		return
	}

	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Error("监听调度CoAP UDP失败", zap.Error(err))
		}
		return
	}

	global.GVA_COAP_DISPATCH_SERVER = listener
	if global.GVA_LOG != nil {
		global.GVA_LOG.Info("调度CoAP服务器已启动", zap.String("address", address))
	}

	go serveCoapDispatch(listener)
}

// serveCoapDispatch 处理调度CoAP请求
func serveCoapDispatch(conn *net.UDPConn) {
	buffer := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if !errors.Is(err, net.ErrClosed) && global.GVA_LOG != nil {
				global.GVA_LOG.Error("调度CoAP服务器读取失败", zap.Error(err))
			}
			return
		}

		packet := append([]byte(nil), buffer[:n]...)
		msg, err := parseCoapMessage(packet)
		if err != nil {
			if global.GVA_LOG != nil {
				global.GVA_LOG.Warn("无效的CoAP消息", zap.String("remote", addr.String()), zap.Error(err))
			}
			continue
		}

		responseType := coapTypeAcknowledgement
		if msg.typ != coapTypeConfirmable {
			responseType = coapTypeNonConfirmable
		}

		var (
			respCode byte = coapCodeNotFound
			payload  []byte
		)

		// 获取请求路径
		path := msg.path()

		// 仅处理 /agvc/data POST请求
		if path == "/agvc/data" && msg.code == coapCodePost {
			respCode, payload = handleDispatchData(context.Background(), msg)
		} else {
			respCode = coapCodeNotFound
			payload = []byte(`{"error":"not found"}`)
		}

		response := buildCoapResponse(msg, responseType, respCode, payload)
		if _, err = conn.WriteToUDP(response, addr); err != nil {
			if global.GVA_LOG != nil {
				global.GVA_LOG.Error("调度CoAP服务器写入失败", zap.Error(err))
			}
		}
	}
}

// handleDispatchData 处理调度数据
func handleDispatchData(ctx context.Context, msg coapMessage) (code byte, payload []byte) {
	// 解析JSON数据
	var dataBatch agvc.AgvcDataBatch
	if err := json.Unmarshal(msg.payload, &dataBatch); err != nil {
		global.GVA_LOG.Error("解析调度数据失败", zap.Error(err))
		return coapCodeBadRequest, []byte(`{"error":"invalid json"}`)
	}

	// 验证数据
	if len(dataBatch) == 0 {
		global.GVA_LOG.Warn("收到空的调度数据")
		return coapCodeBadRequest, []byte(`{"error":"empty data"}`)
	}

	// 转换并存储到调度数据存储
	now := time.Now().Unix()
	dispatchDataList := make([]*agvc_main.RealtimeData, len(dataBatch))
	for i, item := range dataBatch {
		dispatchDataList[i] = &agvc_main.RealtimeData{
			PSID:      item.Psid,
			EQID:      item.Eqid,
			EQType:    item.EqType,
			DataType:  item.DataType,
			Point:     item.Point,
			Value:     item.Value,
			Timestamp: now,
		}
	}

	// 存储到调度数据内存
	agvcMainService.DispatchStorage.StoreBatch(dispatchDataList)

	global.GVA_LOG.Info("调度数据已存储到内存",
		zap.Int("count", len(dataBatch)),
		zap.Int("总数", agvcMainService.DispatchStorage.GetDataCount()))

	return coapCodeCreated, []byte("ok")
}
