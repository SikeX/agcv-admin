#!/bin/bash

# 实现检查脚本

echo "========================================="
echo "AGVC调度控制功能实现检查"
echo "========================================="

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✓${NC} $1"
        return 0
    else
        echo -e "${RED}✗${NC} $1 (缺失)"
        return 1
    fi
}

check_dir() {
    if [ -d "$1" ]; then
        echo -e "${GREEN}✓${NC} $1/"
        return 0
    else
        echo -e "${RED}✗${NC} $1/ (缺失)"
        return 1
    fi
}

echo -e "\n${YELLOW}=== 检查新增文件 ===${NC}"

echo -e "\n模型层:"
check_file "model/agvc/agvc_schedule.go"

echo -e "\n服务层:"
check_file "service/agvc/agcv_main/point_mapper.go"
check_file "service/agvc/agcv_main/dispatch_storage.go"
check_file "service/agvc/agcv_main/schedule_service.go"
check_file "service/agvc/agcv_main/power_aggregator.go"

echo -e "\n初始化层:"
check_file "initialize/coap_dispatch.go"

echo -e "\nAPI层:"
check_file "api/v1/agvc_main/schedule.go"

echo -e "\n路由层:"
check_file "router/agvc_main/enter.go"
check_file "router/agvc_main/schedule.go"

echo -e "\n文档:"
check_file "AGVC_DISPATCH_CONTROL.md"
check_file "test_dispatch.sh"
check_file "../IMPLEMENTATION_SUMMARY.md"

echo -e "\n${YELLOW}=== 检查关键修改 ===${NC}"

echo -e "\n检查core/server.go中的初始化..."
if grep -q "CoapDispatchServer" core/server.go; then
    echo -e "${GREEN}✓${NC} CoapDispatchServer初始化已添加"
else
    echo -e "${RED}✗${NC} CoapDispatchServer初始化缺失"
fi

if grep -q "DispatchStorage.Initialize" core/server.go; then
    echo -e "${GREEN}✓${NC} DispatchStorage初始化已添加"
else
    echo -e "${RED}✗${NC} DispatchStorage初始化缺失"
fi

if grep -q "PointMapper.Initialize" core/server.go; then
    echo -e "${GREEN}✓${NC} PointMapper初始化已添加"
else
    echo -e "${RED}✗${NC} PointMapper初始化缺失"
fi

if grep -q "ScheduleService.Initialize" core/server.go; then
    echo -e "${GREEN}✓${NC} ScheduleService初始化已添加"
else
    echo -e "${RED}✗${NC} ScheduleService初始化缺失"
fi

echo -e "\n检查global/global.go..."
if grep -q "GVA_COAP_DISPATCH_SERVER" global/global.go; then
    echo -e "${GREEN}✓${NC} GVA_COAP_DISPATCH_SERVER已添加"
else
    echo -e "${RED}✗${NC} GVA_COAP_DISPATCH_SERVER缺失"
fi

echo -e "\n检查gorm_biz.go..."
if grep -q "AgvcScheduleCurve" initialize/gorm_biz.go; then
    echo -e "${GREEN}✓${NC} AgvcScheduleCurve表迁移已添加"
else
    echo -e "${RED}✗${NC} AgvcScheduleCurve表迁移缺失"
fi

echo -e "\n检查路由注册..."
if grep -q "ScheduleRouter" initialize/router_biz.go; then
    echo -e "${GREEN}✓${NC} Schedule路由已注册"
else
    echo -e "${RED}✗${NC} Schedule路由未注册"
fi

echo -e "\n${YELLOW}=== 检查Excel文件 ===${NC}"
check_file "设备及测点标准.xlsx"

echo -e "\n${YELLOW}=== 编译测试 ===${NC}"
echo "开始编译..."
if go build -o /tmp/agvc_test 2>/dev/null; then
    echo -e "${GREEN}✓${NC} 编译成功"
    rm -f /tmp/agvc_test
else
    echo -e "${RED}✗${NC} 编译失败，请检查错误"
    go build 2>&1 | head -20
fi

echo -e "\n${YELLOW}=== 功能概览 ===${NC}"
echo ""
echo "1. CoAP 1187端口 - 接收调度命令"
echo "2. 调度数据存储 - DispatchStorage内存存储"
echo "3. 设备点位映射 - 从Excel自动加载"
echo "4. AGC/AVC控制 - 支持调度控制"
echo "5. 功率聚合 - 并网柜缺失时从逆变器聚合"
echo "6. 计划曲线 - 定时执行本地和调度曲线"
echo "7. CoAP发送 - 将控制值发送回采集端"

echo -e "\n${YELLOW}=== 下一步 ===${NC}"
echo ""
echo "1. 确保Excel文件 '设备及测点标准.xlsx' 存在"
echo "2. 配置数据库连接"
echo "3. 启动服务: go run main.go"
echo "4. 运行测试: ./test_dispatch.sh"
echo "5. 查看文档: cat AGVC_DISPATCH_CONTROL.md"
echo ""
echo "========================================="
