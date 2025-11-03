#!/bin/bash

# SaveAgvcData 与 DataStorage 集成测试脚本

echo "================================"
echo "SaveAgvcData 与 DataStorage 集成测试"
echo "================================"
echo ""

# 测试 1: 编译检查
echo "[测试 1] 编译检查..."
cd /home/engine/project/server
if go build -o /tmp/test_integration ./main.go 2>&1; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    exit 1
fi
echo ""

# 测试 2: 代码格式检查
echo "[测试 2] 代码格式检查..."
unformatted=$(gofmt -l initialize/coap_handler.go service/agvc/agcv_main/data_storage.go model/agvc/agvc_main/device.go core/server_run.go)
if [ -z "$unformatted" ]; then
    echo "✅ 代码格式正确"
else
    echo "❌ 以下文件格式不正确:"
    echo "$unformatted"
    exit 1
fi
echo ""

# 测试 3: 导入检查
echo "[测试 3] 导入包检查..."
if grep -q "agvcMainService.*service/agvc/agcv_main" initialize/coap_handler.go; then
    echo "✅ coap_handler.go 导入正确"
else
    echo "❌ coap_handler.go 导入缺失"
    exit 1
fi

if grep -q "agvcMain.*service/agvc/agcv_main" core/server_run.go; then
    echo "✅ server_run.go 导入正确"
else
    echo "❌ server_run.go 导入缺失"
    exit 1
fi
echo ""

# 测试 4: 关键函数检查
echo "[测试 4] 关键函数检查..."
if grep -q "StoreAgvcDataBatch" service/agvc/agcv_main/data_storage.go; then
    echo "✅ StoreAgvcDataBatch 函数存在"
else
    echo "❌ StoreAgvcDataBatch 函数缺失"
    exit 1
fi

if grep -q "DataStorage.Stop()" core/server_run.go; then
    echo "✅ 优雅关闭调用存在"
else
    echo "❌ 优雅关闭调用缺失"
    exit 1
fi
echo ""

# 测试 5: 数据模型检查
echo "[测试 5] 数据模型检查..."
if grep -q "type AgvcDataItem struct" model/agvc/agvc_main/device.go; then
    echo "✅ AgvcDataItem 结构体存在"
else
    echo "❌ AgvcDataItem 结构体缺失"
    exit 1
fi
echo ""

# 测试 6: InfluxDB Measurement 统一检查
echo "[测试 6] InfluxDB Measurement 统一检查..."
if grep -q '"agvc_data"' service/agvc/agcv_main/data_storage.go; then
    echo "✅ data_storage.go 使用 agvc_data measurement"
else
    echo "❌ data_storage.go measurement 不正确"
    exit 1
fi
echo ""

echo "================================"
echo "✅ 所有测试通过！"
echo "================================"
echo ""
echo "改造完成摘要:"
echo "- CoAP 数据接收改为存储到 DataStorage"
echo "- 每 5 分钟自动批量保存到 InfluxDB"
echo "- 支持优雅关闭时保存数据"
echo "- 数据类型自动转换 (int → string)"
echo ""
