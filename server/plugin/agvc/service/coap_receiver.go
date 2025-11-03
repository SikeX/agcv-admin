package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/model"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"go.uber.org/zap"
)

type coapReceiver struct{}

var CoapReceiver = new(coapReceiver)

// RegisterHandlers 注册CoAP处理器
func (s *coapReceiver) RegisterHandlers(router *mux.Router) {
	// 注册数据接收路由
	router.Handle("/agvc/data", mux.HandlerFunc(s.handleDataReceive))
	global.GVA_LOG.Info("CoAP数据接收处理器已注册")
}

// handleDataReceive 处理接收到的CoAP数据
func (s *coapReceiver) handleDataReceive(w mux.ResponseWriter, r *mux.Message) {
	// 读取请求体
	body, err := io.ReadAll(r.Body())
	if err != nil {
		global.GVA_LOG.Error("读取CoAP请求体失败", zap.Error(err))
		s.sendErrorResponse(w, r, codes.BadRequest, "读取请求体失败")
		return
	}

	// 解析JSON数据
	var dataMessages []model.RealtimeData
	if err := json.Unmarshal(body, &dataMessages); err != nil {
		global.GVA_LOG.Error("解析CoAP数据失败", zap.Error(err), zap.String("body", string(body)))
		s.sendErrorResponse(w, r, codes.BadRequest, "数据格式错误")
		return
	}

	// 验证并存储数据
	validData := make([]*model.RealtimeData, 0)
	for _, data := range dataMessages {
		if err := s.validateData(&data); err != nil {
			global.GVA_LOG.Warn("数据验证失败",
				zap.Error(err),
				zap.String("psid", data.PSID),
				zap.String("eqid", data.EQID),
				zap.String("point", data.Point))
			continue
		}
		validData = append(validData, &data)
	}

	// 批量存储数据
	if len(validData) > 0 {
		DataStorage.StoreBatch(validData)
		global.GVA_LOG.Debug("接收并存储CoAP数据",
			zap.Int("总数", len(dataMessages)),
			zap.Int("有效数", len(validData)))
	}

	// 发送成功响应
	s.sendSuccessResponse(w, r, fmt.Sprintf("成功接收%d条数据", len(validData)))
}

// validateData 验证数据完整性
func (s *coapReceiver) validateData(data *model.RealtimeData) error {
	if data.PSID == "" {
		return fmt.Errorf("PSID不能为空")
	}
	if data.EQID == "" {
		return fmt.Errorf("EQID不能为空")
	}
	if data.EQType == "" {
		return fmt.Errorf("EQType不能为空")
	}
	if data.DataType == "" {
		return fmt.Errorf("DataType不能为空")
	}
	if data.Point == "" {
		return fmt.Errorf("Point不能为空")
	}
	if data.Value == nil {
		return fmt.Errorf("Value不能为空")
	}

	// 验证字段长度
	if len(data.PSID) != 3 {
		return fmt.Errorf("PSID必须是3位: %s", data.PSID)
	}
	if len(data.EQID) != 4 {
		return fmt.Errorf("EQID必须是4位: %s", data.EQID)
	}
	if len(data.EQType) != 2 {
		return fmt.Errorf("EQType必须是2位: %s", data.EQType)
	}
	if len(data.DataType) != 2 {
		return fmt.Errorf("DataType必须是2位: %s", data.DataType)
	}

	return nil
}

// sendSuccessResponse 发送成功响应
func (s *coapReceiver) sendSuccessResponse(w mux.ResponseWriter, r *mux.Message, msg string) {
	response := map[string]interface{}{
		"code": 0,
		"msg":  msg,
	}

	data, _ := json.Marshal(response)

	if err := w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(data)); err != nil {
		global.GVA_LOG.Error("发送CoAP响应失败", zap.Error(err))
	}
}

// sendErrorResponse 发送错误响应
func (s *coapReceiver) sendErrorResponse(w mux.ResponseWriter, r *mux.Message, code codes.Code, msg string) {
	response := map[string]interface{}{
		"code": -1,
		"msg":  msg,
	}

	data, _ := json.Marshal(response)

	if err := w.SetResponse(code, message.AppJSON, bytes.NewReader(data)); err != nil {
		global.GVA_LOG.Error("发送CoAP错误响应失败", zap.Error(err))
	}
}
