#!/bin/bash

# 测试脚本：测试调度控制功能

echo "========================================="
echo "AGVC调度控制功能测试脚本"
echo "========================================="

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试1: 发送调度数据到1187端口
echo -e "\n${YELLOW}测试1: 发送调度数据到1187端口${NC}"
echo "发送AGC有功功率目标值1000kW到并网点1..."

# 使用nc (netcat) 发送CoAP数据
# 注意：这需要安装coap-client或使用其他工具
# 这里提供curl风格的示例
cat > /tmp/dispatch_data.json << 'EOF'
[{
  "psid": 1,
  "eqid": 1,
  "eqType": 5,
  "dataType": 5,
  "point": "7",
  "value": 1000.5
}]
EOF

echo "数据已准备：/tmp/dispatch_data.json"
cat /tmp/dispatch_data.json

# 测试2: 创建计划曲线
echo -e "\n${YELLOW}测试2: 创建AGC本地计划曲线${NC}"
echo "创建一个10:00执行的计划..."

# 假设API服务运行在8888端口
API_HOST="http://localhost:8888"

# 注意：需要token，这里仅作示例
curl -X POST "${API_HOST}/agvcMain/schedule/create" \
  -H "Content-Type: application/json" \
  -H "x-token: YOUR_TOKEN_HERE" \
  -d '{
    "bwdNo": 1,
    "type": 1,
    "source": 1,
    "startTime": "10:00",
    "targetValue": 1000,
    "enabled": 1
  }' 2>/dev/null

echo -e "\n${GREEN}如果返回成功，计划曲线已创建${NC}"

# 测试3: 查询计划曲线列表
echo -e "\n${YELLOW}测试3: 查询计划曲线列表${NC}"
curl -X GET "${API_HOST}/agvcMain/schedule/list?bwdNo=1&type=1" \
  -H "x-token: YOUR_TOKEN_HERE" 2>/dev/null

# 测试说明
echo -e "\n${YELLOW}=========================================${NC}"
echo -e "${GREEN}测试完成！${NC}"
echo ""
echo "要使用CoAP测试，需要安装 coap-client："
echo "  Ubuntu/Debian: sudo apt-get install libcoap2-bin"
echo ""
echo "使用coap-client发送数据："
echo '  echo '\''[{"psid":1,"eqid":1,"eqType":5,"dataType":5,"point":"7","value":1000}]'\'' | \'
echo "    coap-client -m post -t application/json coap://localhost:1187/agvc/data"
echo ""
echo "查看服务日志："
echo "  tail -f server.log | grep -E '调度|计划|聚合'"
echo ""
echo "注意事项："
echo "  1. psid固定为1"
echo "  2. 确保设备及测点标准.xlsx文件存在"
echo "  3. 检查数据库中的agvc_bwd_setting表配置"
echo "  4. control_auth=1表示调度控制"
echo "========================================="
