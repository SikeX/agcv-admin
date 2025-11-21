# Point标签映射修复说明

## 问题描述

在原有的实现中，point结构体标签中不同的type可以对应相同的value，这导致了映射冲突。例如：
- 字段A: `point:"value:10,type:1"` (遥信)
- 字段B: `point:"value:10,type:2"` (遥测)

原实现仅使用value作为映射的key，这会导致后面的定义覆盖前面的，造成数据丢失或错误。

## 解决方案

使用 **"type-value" 复合key** 作为映射的唯一标识。

### 修改内容

#### 1. 映射Key格式变更

**修改前：**
```go
// key格式：仅使用value
pointToFieldMapping = map[string]PointInfo{
    "10": PointInfo{...},  // 会被覆盖！
    "10": PointInfo{...},  // 覆盖了上一个
}
```

**修改后：**
```go
// key格式：使用 "type-value" 复合key
pointToFieldMapping = map[string]PointInfo{
    "1-10": PointInfo{...},  // type=1, value=10
    "2-10": PointInfo{...},  // type=2, value=10
}
```

#### 2. 函数签名变更

**修改前：**
```go
func GetFieldNameByPoint(point string) (string, bool)
func GetPointInfoByPoint(point string) (PointInfo, bool)
```

**修改后：**
```go
func GetFieldNameByPoint(dataType int, point string) (string, bool)
func GetPointInfoByPoint(dataType int, point string) (PointInfo, bool)
```

#### 3. 使用示例

**修改前：**
```go
// 无法区分type，可能获取到错误的字段
if fieldName, found := agvc.GetFieldNameByPoint("10"); found {
    // 哪个type=10的字段？不确定！
}
```

**修改后：**
```go
// 明确指定type和value
if fieldName, found := agvc.GetFieldNameByPoint(2, "10"); found {
    // 明确获取 type=2, value=10 的字段
}
```

## 修改的文件列表

### Model层（定义映射）
1. `server/model/agvc/agvc_nbq.go`
   - 修改 `InitPointMapping()` - 使用复合key
   - 修改 `GetFieldNameByPoint(dataType int, point string)` - 新增dataType参数
   - 修改 `GetPointInfoByPoint(dataType int, point string)` - 新增dataType参数
   - 优化 `GetAllPointValues()` - 从复合key提取value并去重
   - 优化 `GetPointValuesByType()` - 过滤type并提取value

2. `server/model/agvc/agvc_bwd.go`
   - 同上（BWD设备）

3. `server/model/agvc/agvc_agc.go`
   - 同上（AGC设备）

4. `server/model/agvc/agvc_avc.go`
   - 同上（AVC设备）

### Service层（使用映射）
5. `server/service/agvc/agvc_his.go`
   - 修改 `getHistoryNbq()` - 调用时传递dataType参数（cons.YC）
   - 修改 `getHistoryBwd()` - 调用时传递dataType参数（cons.YC）
   - 修改 `getHistoryAgc()` - 调用时传递dataType参数（cons.YC或cons.YX）
   - 修改 `getHistoryAvc()` - 调用时传递dataType参数（cons.YC或cons.YX）

6. `server/service/agvc/agvc_real.go`
   - 修改实时数据查询逻辑
   - 尝试匹配type=1和type=2，选择dataType匹配的字段

## 技术细节

### 复合Key生成
```go
// 在InitPointMapping中
key := fmt.Sprintf("%d-%s", dataType, pointValue)
pointToFieldMapping[key] = PointInfo{
    FieldName: field.Name,
    DataType:  dataType,
}
```

### 复合Key查询
```go
// 在GetFieldNameByPoint/GetPointInfoByPoint中
key := fmt.Sprintf("%d-%s", dataType, point)
pointInfo, ok := pointToFieldMapping[key]
```

### GetAllPointValues去重逻辑
```go
func GetAllPointValues() []string {
    pointSet := make(map[string]bool)
    for key := range pointToFieldMapping {
        // key格式为 "type-value"，提取value部分
        parts := strings.Split(key, "-")
        if len(parts) == 2 {
            pointSet[parts[1]] = true  // 去重
        }
    }
    
    points := make([]string, 0, len(pointSet))
    for point := range pointSet {
        points = append(points, point)
    }
    return points
}
```

## 调用场景说明

### 场景1：历史数据查询
历史数据查询通过InfluxDB的Flux查询已经过滤了dataType，因此可以直接使用对应的type常量：

```go
// NBQ历史数据查询（只查询YC类型）
flux := `... |> filter(fn: (r) => r["dataType"] == "2") ...`
// 所以调用时直接使用 cons.YC (值为2)
if fieldName, found := agvc.GetFieldNameByPoint(cons.YC, point); found {
    setFieldValue(nbqHis, fieldName, valueFloat64)
}
```

### 场景2：实时数据查询
实时数据查询可能包含多种type的数据，需要尝试匹配：

```go
// 尝试匹配type=1
if pointInfo, found := agvc.GetPointInfoByPoint(1, point); found {
    pointDataType = pointInfo.DataType
    fieldFound = true
// 尝试匹配type=2
} else if pointInfo, found := agvc.GetPointInfoByPoint(2, point); found {
    pointDataType = pointInfo.DataType
    fieldFound = true
}

// 验证dataType是否匹配
if dataTypeStr == fmt.Sprintf("%d", pointDataType) {
    // 使用找到的dataType
    if fieldName, found := agvc.GetFieldNameByPoint(pointDataType, point); found {
        // 设置字段值
    }
}
```

## 优势

1. **准确性**：避免了不同type相同value的映射冲突
2. **类型安全**：明确区分遥信(YX, type=1)和遥测(YC, type=2)
3. **性能**：依然保持O(1)的查询效率
4. **可维护性**：代码逻辑更清晰，易于理解和维护

## 兼容性说明

这是一个**破坏性修改**，所有调用 `GetFieldNameByPoint` 和 `GetPointInfoByPoint` 的地方都需要更新，传递 `dataType` 参数。

但由于这些函数主要在service层内部使用，影响范围可控。所有相关调用已在本次修复中一并更新。

## 测试建议

1. **单元测试**：验证复合key的生成和查询逻辑
2. **集成测试**：测试历史数据和实时数据查询的正确性
3. **边界测试**：测试同一value不同type的情况
4. **性能测试**：验证修改后的查询性能

## 参考

- InfluxDB数据模型：`dataType` 标签用于区分遥信(1)和遥测(2)
- 常量定义：`server/service/agvc/cons/def.go`
  - `const YX = 1` - 遥信
  - `const YC = 2` - 遥测
