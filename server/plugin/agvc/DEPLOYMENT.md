# AGVC 插件部署指南

## 环境要求

### 必需组件

1. **Go 1.23+**
2. **MySQL 5.7+ / PostgreSQL 9.6+ / SQLite 3**
3. **InfluxDB 2.x**
4. **Redis 6.0+** (可选，用于缓存)

### 推荐配置

- **内存**: 最少2GB，推荐4GB+
- **CPU**: 2核心以上
- **磁盘**: 至少10GB可用空间（用于InfluxDB数据存储）

## 安装步骤

### 1. 确认GVA主系统已安装

AGVC是GVA的插件，需要先安装gin-vue-admin主系统。

```bash
cd /path/to/gin-vue-admin/server
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置InfluxDB

#### 安装InfluxDB 2.x

```bash
# Ubuntu/Debian
wget https://dl.influxdata.com/influxdb/releases/influxdb2-2.7.5-amd64.deb
sudo dpkg -i influxdb2-2.7.5-amd64.deb
sudo systemctl start influxdb

# Docker
docker run -d -p 8086:8086 \
  -v influxdb2:/var/lib/influxdb2 \
  influxdb:2.7
```

#### 初始化InfluxDB

访问 http://localhost:8086 进行初始化：

1. 创建组织（Organization）：如 `myorg`
2. 创建bucket：如 `agvc_data`
3. 生成token

### 4. 配置config.yaml

编辑 `server/config.yaml`，添加或修改以下配置：

```yaml
# InfluxDB配置
influxdb:
  url: http://localhost:8086
  token: YOUR_INFLUXDB_TOKEN
  org: myorg
  bucket: agvc_data

# CoAP配置（在插件配置中）
plugin:
  agvc:
    enable: true
    coapHost: "127.0.0.1"
    coapPort: 1188
```

### 5. 注册插件

插件会自动被GVA系统发现和加载，无需手动注册。

### 6. 启动服务

```bash
cd server
go run main.go
```

或者编译后运行：

```bash
go build -o gva-server
./gva-server
```

### 7. 验证安装

#### 检查数据库表

登录数据库，检查是否创建了以下表：

```sql
SHOW TABLES LIKE 'agvc_%';
```

应该看到：
- `agvc_device`
- `agvc_point_mapping`
- `agvc_agc_config`
- `agvc_agc_regulation_record`
- `agvc_inverter_regulation`
- `agvc_avc_config`
- `agvc_avc_regulation_record`
- `agvc_device_reactive_regulation`

#### 检查API

```bash
# 获取设备列表（需要先登录获取token）
curl -X GET "http://localhost:8888/api/agvc/device/list?page=1&pageSize=10" \
  -H "x-token: YOUR_TOKEN"
```

#### 检查日志

查看启动日志中是否有：

```
[INFO] AGVC插件数据库表初始化成功
[INFO] 测点映射初始化成功
[INFO] 数据存储服务初始化成功
[INFO] AGC控制服务初始化成功
[INFO] AVC控制服务初始化成功
```

## 配置CoAP设备

### 设置CoAP客户端

设备需要配置CoAP客户端，向服务器发送数据：

- **服务器地址**: 服务器IP
- **端口**: 1188（默认）
- **URL**: `/agvc/data`
- **方法**: POST
- **内容类型**: application/json

### 数据格式

```json
[{
  "psid": "001",
  "eqid": "0001",
  "eqType": "01",
  "dataType": "02",
  "point": "401",
  "value": 2.35
}]
```

## 初始化数据

### 1. 创建电站设备

```bash
curl -X POST "http://localhost:8888/api/agvc/device/create" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0000",
    "eqType": "02",
    "dataType": "02",
    "name": "1号电站并网点",
    "status": 1,
    "maxPower": 20.0,
    "maxReg": 5.0,
    "maxReact": 10.0
  }'
```

### 2. 创建逆变器

```bash
curl -X POST "http://localhost:8888/api/agvc/device/create" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0001",
    "eqType": "01",
    "dataType": "02",
    "name": "1号逆变器",
    "status": 1,
    "maxPower": 2.5,
    "maxReg": 0.5,
    "maxReact": 1.0
  }'
```

### 3. 创建AGC配置

```bash
curl -X POST "http://localhost:8888/api/agvc/agc/config" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 0,
    "controlAuth": 1,
    "runMode": 2,
    "jitterRange": 0.5,
    "regPeriod": 30,
    "regStep": 0.1,
    "powerUpperLimit": 20.0,
    "powerLowerLimit": 0.0,
    "dispatchExecValue": 10.0
  }'
```

### 4. 创建AVC配置

```bash
curl -X POST "http://localhost:8888/api/agvc/avc/config" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 0,
    "controlAuth": 1,
    "runMode": 2,
    "targetVoltageLow": 10.0,
    "targetVoltageHigh": 10.5,
    "regPeriod": 60,
    "voltageDeadZone": 0.2,
    "reactiveSensitivity": 1.0
  }'
```

## 生产环境部署

### 1. 使用Systemd管理服务

创建 `/etc/systemd/system/gva-agvc.service`:

```ini
[Unit]
Description=GVA AGVC Service
After=network.target mysql.service influxdb.service

[Service]
Type=simple
User=gva
WorkingDirectory=/opt/gva/server
ExecStart=/opt/gva/server/gva-server
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable gva-agvc
sudo systemctl start gva-agvc
sudo systemctl status gva-agvc
```

### 2. 使用Docker部署

#### Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY server/ .
RUN go mod download
RUN go build -o gva-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/gva-server .
COPY --from=builder /app/config.yaml .
EXPOSE 8888
CMD ["./gva-server"]
```

#### Docker Compose

```yaml
version: '3.8'

services:
  gva-server:
    build: .
    ports:
      - "8888:8888"
      - "1188:1188/udp"  # CoAP端口
    environment:
      - GVA_CONFIG=/app/config.yaml
    depends_on:
      - mysql
      - influxdb
      - redis
    volumes:
      - ./config.yaml:/app/config.yaml
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: gva
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  influxdb:
    image: influxdb:2.7
    ports:
      - "8086:8086"
    volumes:
      - influxdb_data:/var/lib/influxdb2
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=admin123456
      - DOCKER_INFLUXDB_INIT_ORG=myorg
      - DOCKER_INFLUXDB_INIT_BUCKET=agvc_data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  influxdb_data:
  redis_data:
```

启动：

```bash
docker-compose up -d
```

### 3. 使用Nginx反向代理

```nginx
upstream gva_backend {
    server 127.0.0.1:8888;
}

server {
    listen 80;
    server_name your-domain.com;

    # API代理
    location /api/ {
        proxy_pass http://gva_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket支持（如果需要）
    location /ws/ {
        proxy_pass http://gva_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## 监控和维护

### 1. 日志管理

日志位置：根据GVA配置，通常在 `server/log/` 目录

查看实时日志：

```bash
tail -f server/log/server.log | grep AGVC
```

### 2. 性能监控

#### 监控实时数据量

```bash
curl -X GET "http://localhost:8888/api/agvc/device/list?page=1&pageSize=1" \
  -H "x-token: YOUR_TOKEN" | jq '.data.total'
```

#### 监控AGC/AVC运行状态

```bash
# AGC状态
curl -X GET "http://localhost:8888/api/agvc/agc/config?psid=001" \
  -H "x-token: YOUR_TOKEN"

# AVC状态
curl -X GET "http://localhost:8888/api/agvc/avc/config?psid=001" \
  -H "x-token: YOUR_TOKEN"
```

### 3. 数据备份

#### 备份MySQL数据

```bash
mysqldump -u root -p gva > gva_backup_$(date +%Y%m%d).sql
```

#### 备份InfluxDB数据

```bash
influx backup /backup/influxdb_$(date +%Y%m%d) \
  --bucket agvc_data \
  --org myorg \
  --token YOUR_TOKEN
```

### 4. 故障排查

#### 检查进程

```bash
ps aux | grep gva-server
```

#### 检查端口

```bash
netstat -tuln | grep 8888  # HTTP
netstat -tuln | grep 1188  # CoAP
```

#### 检查连接

```bash
# InfluxDB
curl http://localhost:8086/health

# MySQL
mysql -u root -p -e "SELECT 1"
```

## 升级指南

### 1. 备份数据

```bash
# 备份数据库
mysqldump -u root -p gva > backup.sql

# 备份配置
cp config.yaml config.yaml.bak
```

### 2. 更新代码

```bash
git pull origin main
cd server
go mod tidy
```

### 3. 停止服务

```bash
sudo systemctl stop gva-agvc
```

### 4. 重新编译

```bash
go build -o gva-server main.go
```

### 5. 启动服务

```bash
sudo systemctl start gva-agvc
```

### 6. 验证升级

```bash
# 检查服务状态
sudo systemctl status gva-agvc

# 检查API
curl http://localhost:8888/api/agvc/device/list?page=1&pageSize=1
```

## 常见问题

### Q1: InfluxDB连接失败

**A**: 检查InfluxDB是否运行，token是否正确：

```bash
curl http://localhost:8086/health
```

### Q2: CoAP数据无法接收

**A**: 检查防火墙是否开放UDP 1188端口：

```bash
sudo ufw allow 1188/udp
```

### Q3: AGC/AVC无法启动

**A**: 检查配置是否创建，日志中是否有错误信息。

### Q4: 内存占用过高

**A**: 检查实时数据量，考虑缩短保存周期（默认5分钟）。

## 性能优化

### 1. InfluxDB优化

```toml
# influxdb.conf
[data]
  cache-max-memory-size = "1g"
  cache-snapshot-memory-size = "256m"
```

### 2. MySQL优化

```ini
# my.cnf
[mysqld]
innodb_buffer_pool_size = 2G
innodb_log_file_size = 512M
```

### 3. Go应用优化

```bash
# 设置GOMAXPROCS
export GOMAXPROCS=4

# 启用性能分析
go run main.go -cpuprofile=cpu.prof -memprofile=mem.prof
```

## 联系支持

如有问题，请提交Issue到项目仓库，或联系技术支持团队。
