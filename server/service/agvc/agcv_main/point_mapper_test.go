package agcv_main

import (
    "testing"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "go.uber.org/zap"
)

func init() {
    // 初始化全局日志，避免nil pointer
    global.GVA_LOG = zap.NewNop()
}

func TestPointMapper_Initialize(t *testing.T) {
    pm := new(pointMapper)
    err := pm.Initialize()
    if err != nil {
        t.Fatalf("初始化失败: %v", err)
    }

    // 测试逆变器类型（eqType=2）
    eqType := 2
    if len(pm.mapping[eqType]) == 0 {
        t.Error("逆变器点位映射为空")
    }

    // 测试获取点位ID
    pointName := "交流功率(kW)"
    pointID, err := pm.GetPointID(eqType, pointName)
    if err != nil {
        t.Errorf("获取点位ID失败: %v", err)
    }
    if pointID != "0" {
        t.Errorf("期望点位ID为'0'，实际为'%s'", pointID)
    }

    // 测试反向映射 - 注意：相同的pointID可能对应多个点位名称，取决于最后加载的
    name, err := pm.GetPointName(eqType, pointID)
    if err != nil {
        t.Errorf("获取点位名称失败: %v", err)
    }
    if name == "" {
        t.Errorf("获取的点位名称为空")
    }
    t.Logf("点位ID '%s' 对应的名称: '%s'", pointID, name)

    // 测试获取所有点位
    allPoints := pm.GetAllPointsForType(eqType)
    if len(allPoints) == 0 {
        t.Error("获取所有点位失败，返回为空")
    }

    t.Logf("成功加载 %d 个点位映射", len(allPoints))
}

func TestPointMapper_GetPointID(t *testing.T) {
    pm := new(pointMapper)
    pm.Initialize()

    tests := []struct {
        name      string
        eqType    int
        pointName string
        wantID    string
        wantErr   bool
    }{
        {"逆变器-总发电量", 2, "总发电量(kWh)", "大于0", false},
        {"逆变器-日发电量", 2, "日发电量(kWh)", "0", false},
        {"逆变器-不存在的点位", 2, "不存在的点位", "", true},
        {"环境监测仪-环境温度", 8, "环境温度(℃)", "-20", false},
        // 注意：箱变的"A相电压Ua"有两个值，最后加载的是"低压侧"
        {"箱变-A相电压", 4, "A相电压Ua", "低压侧", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotID, err := pm.GetPointID(tt.eqType, tt.pointName)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetPointID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if gotID != tt.wantID {
                t.Errorf("GetPointID() = %v, want %v", gotID, tt.wantID)
            }
        })
    }
}

func TestDevicePointStandards(t *testing.T) {
    // 验证数据结构是否正确
    expectedSheets := []string{"逆变器", "集中逆变器", "箱变", "开关柜和保护装置", "电能表", "环境监测仪"}
    
    for _, sheet := range expectedSheets {
        points, exists := devicePointStandards[sheet]
        if !exists {
            t.Errorf("缺少设备类型: %s", sheet)
            continue
        }
        
        if len(points) == 0 {
            t.Errorf("设备类型 %s 的点位数据为空", sheet)
        }
        
        t.Logf("设备类型 %s: %d 个点位", sheet, len(points))
    }
}
