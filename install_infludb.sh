#!/bin/bash
set -euo pipefail

# ==============================================
# 配置参数（可根据实际tar包调整）
# ==============================================
# 预期的InfluxDB tar包前缀（无需包含版本号和后缀）
TAR_PREFIX="influxdb2"
# 目标安装目录
INSTALL_DIR="/usr/local/influxdb"
# 系统服务名
SERVICE_NAME="influxdb"
# 初始化配置
ORG="test"
BUCKET="test"
USERNAME="admin"
PASSWORD="12345678"
# InfluxDB监听端口（默认8086）
LISTEN_PORT="8086"
# 环境变量持久化文件（当前用户生效）
ENV_FILE="$HOME/.bashrc"

# ==============================================
# 1. 检查依赖工具
# ==============================================
check_dependency() {
    local cmd=$1
    if ! command -v $cmd &> /dev/null; then
        echo "错误：未找到必需工具 $cmd，请先安装"
        exit 1
    fi
}

echo "=== 检查依赖工具 ==="
check_dependency "tar"
check_dependency "curl"
check_dependency "systemctl"
check_dependency "grep"
check_dependency "sed"

# ==============================================
# 2. 检查tar包是否存在
# ==============================================
echo -e "\n=== 检查InfluxDB tar包 ==="
TAR_FILE=$(ls -1 "${TAR_PREFIX}"-*.tar.gz 2>/dev/null | head -n 1)
if [ -z "$TAR_FILE" ]; then
    echo "错误：当前目录未找到 ${TAR_PREFIX}-*.tar.gz 格式的tar包"
    echo "请将InfluxDB tar包放在当前目录，命名格式如：${TAR_PREFIX}-2.7.1-linux-amd64.tar.gz"
    exit 1
fi
echo "找到tar包：$TAR_FILE"

# ==============================================
# 3. 解压安装InfluxDB
# ==============================================
echo -e "\n=== 安装InfluxDB ==="
# 创建安装目录
sudo mkdir -p "$INSTALL_DIR"
# 解压tar包（忽略顶层目录，直接解压到安装目录）
sudo tar -zxf "$TAR_FILE" --strip-components=1 -C "$INSTALL_DIR"
# 设置权限
sudo chown -R $USER:$USER "$INSTALL_DIR"
# 将influx CLI加入系统PATH（临时生效，后续会持久化）
export PATH="$INSTALL_DIR:$PATH"
if ! command -v influx &> /dev/null; then
    echo "错误：InfluxDB解压失败，未找到influx命令"
    exit 1
fi
echo "InfluxDB安装完成：$INSTALL_DIR"
echo "当前influx版本：$(influx version)"

# ==============================================
# 4. 配置系统服务（systemd）
# ==============================================
echo -e "\n=== 配置系统服务 ==="
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
if [ -f "$SERVICE_FILE" ]; then
    echo "系统服务已存在，跳过创建"
else
    sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=InfluxDB v2 Service
Documentation=https://docs.influxdata.com/influxdb/
After=network.target

[Service]
User=$USER
Group=$USER
ExecStart=$INSTALL_DIR/influxd --http-bind-address :$LISTEN_PORT
Restart=on-failure
RestartSec=5s
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
    # 重载systemd配置
    sudo systemctl daemon-reload
    # 设置开机自启
    sudo systemctl enable "$SERVICE_NAME"
    echo "系统服务创建完成：$SERVICE_FILE"
fi

# 启动InfluxDB服务
if ! sudo systemctl is-active --quiet "$SERVICE_NAME"; then
    echo "启动InfluxDB服务..."
    sudo systemctl start "$SERVICE_NAME"
    # 等待服务启动（最多等待30秒）
    echo "等待服务就绪（监听端口$LISTEN_PORT）..."
    for ((i=0; i<30; i++)); do
        if curl -s "http://localhost:$LISTEN_PORT/health" | grep -q "ready"; then
            echo "InfluxDB服务启动成功"
            break
        fi
        sleep 1
    done
    if ! curl -s "http://localhost:$LISTEN_PORT/health" | grep -q "ready"; then
        echo "错误：InfluxDB服务启动失败，无法访问健康检查接口"
        exit 1
    fi
else
    echo "InfluxDB服务已运行"
fi

# ==============================================
# 5. 初始化InfluxDB（创建用户、org、bucket）
# ==============================================
echo -e "\n=== 初始化InfluxDB ==="
# 检查是否已初始化（通过查看默认配置文件）
CONFIG_FILE="$HOME/.influxdbv2/configs"
if [ -f "$CONFIG_FILE" ] && influx config list | grep -q "$ORG"; then
    echo "InfluxDB已初始化，跳过创建"
    # 提取已存在的token
    INFLUX_TOKEN=$(influx config get --org "$ORG" --format csv | grep "$ORG" | cut -d',' -f3)
else
    echo "开始初始化：org=$ORG, bucket=$BUCKET, username=$USERNAME"
    # 执行初始化命令，捕获输出中的token
    INIT_OUTPUT=$(influx setup \
        --username "$USERNAME" \
        --password "$PASSWORD" \
        --org "$ORG" \
        --bucket "$BUCKET" \
        --retention 0 \
        --force)
    # 从输出中提取token（正则匹配）
    INFLUX_TOKEN=$(echo "$INIT_OUTPUT" | grep -oP 'Your API token is: \K\S+')
    if [ -z "$INFLUX_TOKEN" ]; then
        echo "错误：初始化失败，未提取到API Token"
        echo "初始化输出：$INIT_OUTPUT"
        exit 1
    fi
    echo "初始化成功，Token已生成"
fi

# ==============================================
# 6. 持久化环境变量
# ==============================================
echo -e "\n=== 配置环境变量 ==="
# 检查环境变量是否已存在
if grep -q "INFLUX_TOKEN" "$ENV_FILE" && grep -q "INFLUX_ORG" "$ENV_FILE"; then
    echo "环境变量已存在，更新值..."
    # 更新已存在的环境变量
    sudo sed -i "/^export INFLUX_TOKEN=/d" "$ENV_FILE"
    sudo sed -i "/^export INFLUX_ORG=/d" "$ENV_FILE"
    sudo sed -i "/^export INFLUX_BUCKET=/d" "$ENV_FILE"
    sudo sed -i "/^export INFLUX_URL=/d" "$ENV_FILE"
fi

# 追加环境变量到配置文件
cat <<EOF >> "$ENV_FILE"
# InfluxDB Environment Variables
export INFLUX_URL="http://localhost:$LISTEN_PORT"
export INFLUX_TOKEN="$INFLUX_TOKEN"
export INFLUX_ORG="$ORG"
export INFLUX_BUCKET="$BUCKET"
EOF

# 使环境变量立即生效（当前终端）
source "$ENV_FILE"

echo "环境变量已持久化到：$ENV_FILE"
echo "当前终端已生效，新终端需重新登录或执行：source $ENV_FILE"

# ==============================================
# 7. 验证配置
# ==============================================
echo -e "\n=== 验证配置 ==="
echo "InfluxDB URL: $INFLUX_URL"
echo "InfluxDB Org: $INFLUX_ORG"
echo "InfluxDB Bucket: $INFLUX_BUCKET"
echo "InfluxDB Token: $INFLUX_TOKEN"
echo -e "\n验证bucket是否存在："
if influx bucket list --org "$ORG" | grep -q "$BUCKET"; then
    echo "✅ Bucket '$BUCKET' 存在"
fi
echo -e "\n验证用户是否存在："
if influx user list --org "$ORG" | grep -q "$USERNAME"; then
    echo "✅ User '$USERNAME' 存在"
fi

# ==============================================
# 安装完成提示
# ==============================================
echo -e "\n========================================"
echo "InfluxDB安装配置完成！"
echo "========================================"
echo "Web UI地址：$INFLUX_URL"
echo "登录信息："
echo "  用户名：$USERNAME"
echo "  密码：$PASSWORD"
echo "  Org：$ORG"
echo "  Bucket：$BUCKET"
echo "环境变量：已写入 $ENV_FILE，包含INFLUX_TOKEN等关键信息"
echo "使用示例：influx query 'from(bucket:\"$BUCKET\") |> range(start:-1h)'"
echo "========================================"