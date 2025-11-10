package agcv_main

import (
    "fmt"
    "sync"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
    "go.uber.org/zap"
)

// DispatchStorage 调度数据存储服务（专门用于1187端口接收的调度数据）
type dispatchStorage struct {
    mu           sync.RWMutex
    dispatchData map[string]*agvc_main.RealtimeData // key: psid_eqid_eqType_dataType_point
}

var DispatchStorage = new(dispatchStorage)

// Initialize 初始化调度数据存储
func (ds *dispatchStorage) Initialize() {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    ds.dispatchData = make(map[string]*agvc_main.RealtimeData)
    global.GVA_LOG.Info("调度数据存储服务初始化成功")
}

// StoreData 存储调度数据
func (ds *dispatchStorage) StoreData(data *agvc_main.RealtimeData) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    key := ds.makeKey(data.PSID, data.EQID, data.EQType, data.DataType, data.Point)
    ds.dispatchData[key] = data

    global.GVA_LOG.Debug("存储调度数据",
        zap.String("key", key),
        zap.Any("value", data.Value))
}

// StoreBatch 批量存储调度数据
func (ds *dispatchStorage) StoreBatch(dataList []*agvc_main.RealtimeData) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    for _, data := range dataList {
        key := ds.makeKey(data.PSID, data.EQID, data.EQType, data.DataType, data.Point)
        ds.dispatchData[key] = data
    }

    global.GVA_LOG.Debug("批量存储调度数据", zap.Int("count", len(dataList)))
}

// GetData 获取调度数据
func (ds *dispatchStorage) GetData(psid, eqid, eqType, dataType int, point string) (*agvc_main.RealtimeData, bool) {
    ds.mu.RLock()
    defer ds.mu.RUnlock()

    key := ds.makeKey(psid, eqid, eqType, dataType, point)
    data, exists := ds.dispatchData[key]
    return data, exists
}

// GetDataAsFloat64 获取调度数据并转换为float64
func (ds *dispatchStorage) GetDataAsFloat64(psid, eqid, eqType, dataType int, point string) (float64, error) {
    data, exists := ds.GetData(psid, eqid, eqType, dataType, point)
    if !exists {
        return 0, fmt.Errorf("调度数据不存在")
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

// GetDataCount 获取调度数据总数
func (ds *dispatchStorage) GetDataCount() int {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    return len(ds.dispatchData)
}

// makeKey 生成存储键
func (ds *dispatchStorage) makeKey(psid, eqid, eqType, dataType int, point string) string {
    return fmt.Sprintf("%d_%d_%d_%d_%s", psid, eqid, eqType, dataType, point)
}

// ClearData 清空调度数据
func (ds *dispatchStorage) ClearData() {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    ds.dispatchData = make(map[string]*agvc_main.RealtimeData)
    global.GVA_LOG.Info("调度数据已清空")
}

// SetData 设置指定点位的调度数据值
func (ds *dispatchStorage) SetData(psid, eqid, eqType, dataType int, point string, value interface{}) error {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    key := ds.makeKey(psid, eqid, eqType, dataType, point)
    data := &agvc_main.RealtimeData{
        PSID:     psid,
        EQID:     eqid,
        EQType:   eqType,
        DataType: dataType,
        Point:    point,
        Value:    value,
    }
    ds.dispatchData[key] = data

    global.GVA_LOG.Debug("设置调度数据",
        zap.String("key", key),
        zap.Any("value", value))
    
    return nil
}
