package agcv_main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap"
)

type dataStorage struct {
	mu           sync.RWMutex
	realtimeData map[string]*agvc_main.RealtimeData // key: psid_eqid_eqType_dataType_point
	saveTimer    *time.Ticker
	stopChan     chan struct{}
}

var DataStorage = new(dataStorage)

// Initialize 初始化数据存储服务
func (s *dataStorage) Initialize() {
	s.realtimeData = make(map[string]*agvc_main.RealtimeData)
	s.stopChan = make(chan struct{})

	// 启动5分钟定时保存到InfluxDB
	s.saveTimer = time.NewTicker(5 * time.Minute)
	go s.periodicSave()

	global.GVA_LOG.Info("数据存储服务初始化成功")
}

// Stop 停止数据存储服务
func (s *dataStorage) Stop() {
	if s.saveTimer != nil {
		s.saveTimer.Stop()
	}
	close(s.stopChan)
	global.GVA_LOG.Info("数据存储服务已停止")
}

// StoreData 存储实时数据到内存
func (s *dataStorage) StoreData(data *agvc_main.RealtimeData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.makeKey(data.PSID, data.EQID, data.EQType, data.DataType, data.Point)
	data.Timestamp = time.Now().Unix()
	s.realtimeData[key] = data
}

// StoreBatch 批量存储数据
func (s *dataStorage) StoreBatch(dataList []*agvc_main.RealtimeData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	for _, data := range dataList {
		key := s.makeKey(data.PSID, data.EQID, data.EQType, data.DataType, data.Point)
		data.Timestamp = now
		s.realtimeData[key] = data
	}
}

// GetData 从内存获取实时数据
func (s *dataStorage) GetData(psid, eqid, eqType, dataType, point string) (*agvc_main.RealtimeData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.makeKey(psid, eqid, eqType, dataType, point)
	data, exists := s.realtimeData[key]
	return data, exists
}

// GetDeviceData 获取设备的所有实时数据
func (s *dataStorage) GetDeviceData(psid, eqid, eqType string, dataType string) map[string]*agvc_main.RealtimeData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*agvc_main.RealtimeData)
	prefix := fmt.Sprintf("%s_%s_%s_%s_", psid, eqid, eqType, dataType)

	for key, data := range s.realtimeData {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			result[data.Point] = data
		}
	}

	return result
}

// GetDeviceAllData 获取设备所有数据类型的实时数据
func (s *dataStorage) GetDeviceAllData(psid, eqid, eqType string) map[string]*agvc_main.RealtimeData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*agvc_main.RealtimeData)
	prefix := fmt.Sprintf("%s_%s_%s_", psid, eqid, eqType)

	for key, data := range s.realtimeData {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			fullKey := fmt.Sprintf("%s_%s", data.DataType, data.Point)
			result[fullKey] = data
		}
	}

	return result
}

// periodicSave 定期保存数据到InfluxDB
func (s *dataStorage) periodicSave() {
	for {
		select {
		case <-s.saveTimer.C:
			if err := s.saveToInfluxDB(); err != nil {
				global.GVA_LOG.Error("保存数据到InfluxDB失败", zap.Error(err))
			} else {
				global.GVA_LOG.Info("数据已保存到InfluxDB", zap.Int("数据量", len(s.realtimeData)))
			}
		case <-s.stopChan:
			return
		}
	}
}

// saveToInfluxDB 保存数据到InfluxDB
func (s *dataStorage) saveToInfluxDB() error {
	if global.GVA_INFLUXDB == nil {
		return fmt.Errorf("InfluxDB客户端未初始化")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	var points []*write.Point
	now := time.Now()

	for _, data := range s.realtimeData {
		// 创建InfluxDB点
		point := influxdb2.NewPoint(
			"device_data",
			map[string]string{
				"psid":     data.PSID,
				"eqid":     data.EQID,
				"eqType":   data.EQType,
				"dataType": data.DataType,
				"point":    data.Point,
			},
			map[string]interface{}{
				"value": data.Value,
			},
			now,
		)
		points = append(points, point)
	}

	if len(points) > 0 {
		return writeAPI.WritePoint(context.Background(), points...)
	}

	return nil
}

// makeKey 生成Map键
func (s *dataStorage) makeKey(psid, eqid, eqType, dataType, point string) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s", psid, eqid, eqType, dataType, point)
}

// GetDataAsFloat64 获取数据并转换为float64（用于数值计算）
func (s *dataStorage) GetDataAsFloat64(psid, eqid, eqType, dataType, point string) (float64, error) {
	data, exists := s.GetData(psid, eqid, eqType, dataType, point)
	if !exists {
		return 0, fmt.Errorf("数据不存在")
	}

	switch v := data.Value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		var f float64
		if _, err := fmt.Sscanf(v, "%f", &f); err != nil {
			return 0, fmt.Errorf("无法转换字符串为float64: %v", err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("不支持的数据类型: %T", v)
	}
}

// GetDataAsInt 获取数据并转换为int（用于状态判断）
func (s *dataStorage) GetDataAsInt(psid, eqid, eqType, dataType, point string) (int, error) {
	data, exists := s.GetData(psid, eqid, eqType, dataType, point)
	if !exists {
		return 0, fmt.Errorf("数据不存在")
	}

	switch v := data.Value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("不支持的数据类型: %T", v)
	}
}

// ExportSnapshot 导出当前数据快照（用于调试）
func (s *dataStorage) ExportSnapshot() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, _ := json.MarshalIndent(s.realtimeData, "", "  ")
	return string(data)
}

// GetDataCount 获取当前存储的数据量
func (s *dataStorage) GetDataCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.realtimeData)
}
