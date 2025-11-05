package agcv_main

import (
    "fmt"
    "sync"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "go.uber.org/zap"
)

// PointMapper 点位映射服务
type pointMapper struct {
    mu sync.RWMutex
    // 映射结构: 设备类型 -> 测点名称 -> 点标识
    mapping map[int]map[string]string
    // 反向映射: 设备类型 -> 点标识 -> 测点名称
    reverseMapping map[int]map[string]string
}

var PointMapper = new(pointMapper)

// Initialize 初始化点位映射
func (pm *pointMapper) Initialize() error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    pm.mapping = make(map[int]map[string]string)
    pm.reverseMapping = make(map[int]map[string]string)

    // 从内联数据加载映射关系
    pm.loadFromInlineData()

    global.GVA_LOG.Info("设备点位映射初始化成功")
    return nil
}

// loadFromInlineData 从内联数据加载映射关系
func (pm *pointMapper) loadFromInlineData() {
    // 设备类型映射
    sheetToEqType := map[string]int{
        "逆变器":      2,
        "集中逆变器":    2, // 与逆变器使用同一类型
        "箱变":       4,
        "开关柜和保护装置": 5,
        "电能表":      7,
        "环境监测仪":    8,
        "agc":      98,
        "agv":      99,
    }

    for sheetName, eqType := range sheetToEqType {
        // 从内联数据获取该设备类型的点位
        points, exists := devicePointStandards[sheetName]
        if !exists {
            continue
        }

        // 初始化该设备类型的映射
        if pm.mapping[eqType] == nil {
            pm.mapping[eqType] = make(map[string]string)
        }
        if pm.reverseMapping[eqType] == nil {
            pm.reverseMapping[eqType] = make(map[string]string)
        }

        // 加载点位数据
        for _, point := range points {
            if point.PointName != "" && point.PointID != "" {
                pm.mapping[eqType][point.PointName] = point.PointID
                pm.reverseMapping[eqType][point.PointID] = point.PointName
            }
        }

        global.GVA_LOG.Info("加载设备映射",
            zap.String("设备类型", sheetName),
            zap.Int("eqType", eqType),
            zap.Int("点位数量", len(pm.mapping[eqType])))
    }
}

// GetPointID 根据设备类型和测点名称获取点标识
func (pm *pointMapper) GetPointID(eqType int, pointName string) (string, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()

    if typeMap, ok := pm.mapping[eqType]; ok {
        if pointID, ok := typeMap[pointName]; ok {
            return pointID, nil
        }
    }

    return "", fmt.Errorf("未找到点位映射: eqType=%d, pointName=%s", eqType, pointName)
}

// GetPointName 根据设备类型和点标识获取测点名称
func (pm *pointMapper) GetPointName(eqType int, pointID string) (string, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()

    if typeMap, ok := pm.reverseMapping[eqType]; ok {
        if pointName, ok := typeMap[pointID]; ok {
            return pointName, nil
        }
    }

    return "", fmt.Errorf("未找到点位映射: eqType=%d, pointID=%s", eqType, pointID)
}

// GetAllPointsForType 获取某设备类型的所有点位
func (pm *pointMapper) GetAllPointsForType(eqType int) map[string]string {
    pm.mu.RLock()
    defer pm.mu.RUnlock()

    result := make(map[string]string)
    if typeMap, ok := pm.mapping[eqType]; ok {
        for k, v := range typeMap {
            result[k] = v
        }
    }
    return result
}
