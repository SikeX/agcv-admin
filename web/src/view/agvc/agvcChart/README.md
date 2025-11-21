# AGVC 图表监控页面

## 功能说明

本模块提供了两个实时监控图表页面：

### 1. 电站出力监控 (powerChart.vue)
- **功能**：实时监控电站的有功功率
- **数据来源**：
  - 当前有功功率：来自 BWD (并网点) 的 CurrentActivePower 字段
  - 目标有功功率：来自 AGC 的 TargetActivePower 字段
- **特性**：
  - 5秒自动刷新
  - ECharts 折线图展示
  - 时间窗口：最近5分钟数据
  - 支持手动开启/关闭自动刷新

### 2. 电压和无功监控 (voltageReactiveChart.vue)
- **功能**：实时监控电压和无功功率
- **数据来源**：
  - 当前电压：来自 BWD 的 CurrentVoltage 字段
  - 当前无功：来自 BWD 的 CurrentReactivePower 字段
  - 目标电压：来自 AVC 的 TargetVoltage 字段
  - 目标无功：来自 AVC 的 TargetReactivePower 字段
- **特性**：
  - 5秒自动刷新
  - 双图表并排展示（电压图 + 无功图）
  - 时间窗口：最近5分钟数据
  - 支持手动开启/关闭自动刷新

## 技术实现

### 后端接口

#### 1. 获取电站出力图表数据
- **路由**：`GET /agvcChart/getPowerChartData`
- **参数**：
  - `psid`: 电站编号
  - `eqid`: 设备编号
  - `startTime`: 开始时间 (ISO 8601格式)
  - `endTime`: 结束时间 (ISO 8601格式)
- **返回**：包含时间序列数据的数组
  ```json
  [
    {
      "time": "2024-01-01T12:00:00Z",
      "currentActivePower": 1500.5,
      "targetActivePower": 1600.0
    }
  ]
  ```

#### 2. 获取电压和无功图表数据
- **路由**：`GET /agvcChart/getVoltageReactiveChartData`
- **参数**：同上
- **返回**：包含时间序列数据的数组
  ```json
  [
    {
      "time": "2024-01-01T12:00:00Z",
      "currentVoltage": 10.5,
      "currentReactivePower": 200.3,
      "targetVoltage": 10.0,
      "targetReactivePower": 180.0
    }
  ]
  ```

### InfluxDB 查询优化
- **aggregateWindow**: 5秒聚合窗口
- **数据源**：
  - BWD: eqType="bwd", dataType=2 (遥测YC)
  - AGC: eqType="agc", dataType=2 (遥测YC)
  - AVC: eqType="avc", dataType=2 (遥测YC)
- **点位映射**：
  - CurrentActivePower: point_10
  - CurrentVoltage: point_12
  - CurrentReactivePower: point_8
  - TargetActivePower (AGC): point_403
  - TargetVoltage (AVC): point_403
  - TargetReactivePower (AVC): point_404

## 使用方法

1. 在路由配置中添加以下路由：
   ```javascript
   {
     path: 'powerChart',
     name: 'powerChart',
     component: () => import('@/view/agvc/agvcChart/powerChart.vue'),
     meta: {
       title: '电站出力监控'
     }
   },
   {
     path: 'voltageReactiveChart',
     name: 'voltageReactiveChart',
     component: () => import('@/view/agvc/agvcChart/voltageReactiveChart.vue'),
     meta: {
       title: '电压无功监控'
     }
   }
   ```

2. 在菜单中添加对应的菜单项

3. 访问页面：
   - 电站出力监控：`/agvc/powerChart`
   - 电压无功监控：`/agvc/voltageReactiveChart`

## 注意事项

1. **时间范围**：默认查询最近5分钟的数据，避免数据量过大
2. **自动刷新**：每5秒自动刷新一次数据，可手动控制开关
3. **数据聚合**：使用5秒聚合窗口，减少数据点数量
4. **空数据处理**：当查询无数据时，会显示友好提示
5. **图表响应式**：图表会随窗口大小自动调整

## 扩展性

如需添加更多监控指标，可以：
1. 在服务层添加新的查询方法
2. 在API层添加新的接口
3. 创建新的前端页面和图表组件
4. 按照现有模式配置路由和菜单
