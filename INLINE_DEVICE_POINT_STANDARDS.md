# 设备及测点标准数据内联化

## 变更说明

将原本需要从Excel文件（`设备及测点标准.xlsx`）读取的设备点位映射数据，转换为代码中的内联数据结构，避免对外部文件的依赖。

## 变更内容

### 1. 新增文件

#### `server/service/agvc/agcv_main/point_mapper_data.go`
- 定义了 `devicePointStandards` 变量，包含所有设备类型的点位映射数据
- 定义了 `PointMapping` 结构体，用于表示点位映射关系
- 包含以下设备类型的完整数据：
  - 逆变器（77个点位）
  - 集中逆变器（34个点位）
  - 箱变（3个点位）
  - 开关柜和保护装置（3个点位）
  - 电能表（17个点位）
  - 环境监测仪（17个点位）
  - agc（预留）
  - agv（预留）

#### `server/service/agvc/agcv_main/point_mapper_test.go`
- 添加了完整的单元测试
- 测试初始化、点位查询、反向映射等功能
- 所有测试均通过

### 2. 修改文件

#### `server/service/agvc/agcv_main/point_mapper.go`
- 移除了对 `github.com/xuri/excelize/v2` 包的依赖
- 移除了 `loadFromExcel` 方法
- 新增了 `loadFromInlineData` 方法，从内联数据加载映射关系
- 修改了 `Initialize` 方法，调用新的 `loadFromInlineData` 方法
- 功能保持不变，向后兼容

#### 其他格式化修复
修复了以下文件中的格式化错误（将 `%s` 改为 `%d` 以匹配int类型参数）：
- `server/service/agvc/agcv_main/agc.go`
- `server/service/agvc/agcv_main/avc.go`
- `server/service/agvc/agcv_main/data_storage.go`

## 优势

1. **无外部依赖**：不再需要Excel文件，减少部署复杂度
2. **性能提升**：避免每次启动时读取和解析Excel文件
3. **易于维护**：数据以Go代码形式存在，便于版本控制和代码审查
4. **类型安全**：编译时即可发现数据错误，而非运行时
5. **测试友好**：便于编写单元测试，不依赖外部文件

## 测试验证

所有测试均通过：
```bash
cd /home/engine/project/server
go test -v ./service/agvc/agcv_main -run TestPointMapper
go test -v ./service/agvc/agcv_main -run TestDevicePointStandards
```

## 数据一致性

所有数据均从原Excel文件 `设备及测点标准.xlsx` 提取，保持完全一致。

## 注意事项

1. 当Excel源文件更新时，需要同步更新 `point_mapper_data.go` 中的数据
2. 对于具有重复点位名称的设备类型（如箱变的"A相电压Ua"），只保留最后一个值
3. 反向映射（点标识→测点名称）同样会保留最后一个出现的映射关系

## 向后兼容

本次修改完全向后兼容，不影响现有功能的使用方式：
- `PointMapper.Initialize()` - 初始化点位映射
- `PointMapper.GetPointID(eqType, pointName)` - 获取点标识
- `PointMapper.GetPointName(eqType, pointID)` - 获取测点名称
- `PointMapper.GetAllPointsForType(eqType)` - 获取所有点位
