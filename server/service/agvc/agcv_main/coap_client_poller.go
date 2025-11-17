package agcv_main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
    "github.com/plgd-dev/go-coap/v3/message"
    "github.com/plgd-dev/go-coap/v3/message/codes"
    "github.com/plgd-dev/go-coap/v3/udp"
    "go.uber.org/zap"
)

type coapClientPoller struct {
    stopChan   chan struct{}
    pollTicker *time.Ticker
    serverHost string
    serverPort int
    pollInterval time.Duration
}

var CoapClientPoller = new(coapClientPoller)

// Initialize 初始化CoAP客户端轮询服务
func (p *coapClientPoller) Initialize() {
    p.stopChan = make(chan struct{})
    
    // 从配置读取CoAP服务器地址
    cfg := global.GVA_CONFIG.Coap
    
    // 检查是否启用客户端轮询
    if !cfg.ClientEnable {
        global.GVA_LOG.Info("CoAP客户端轮询服务未启用")
        return
    }
    
    // 设置服务器地址，默认使用127.0.0.1:5589
    p.serverHost = cfg.ClientHost
    if p.serverHost == "" {
        p.serverHost = "127.0.0.1"
    }
    
    p.serverPort = cfg.ClientPort
    if p.serverPort == 0 {
        p.serverPort = 5589
    }
    
    // 轮询间隔，默认5秒
    pollSeconds := cfg.PollInterval
    if pollSeconds == 0 {
        pollSeconds = 5
    }
    p.pollInterval = time.Duration(pollSeconds) * time.Second
    
    global.GVA_LOG.Info("CoAP客户端轮询服务初始化",
        zap.String("server", fmt.Sprintf("%s:%d", p.serverHost, p.serverPort)),
        zap.Duration("interval", p.pollInterval))
    
    // 启动轮询goroutine
    go p.startPolling()
}

// Stop 停止轮询服务
func (p *coapClientPoller) Stop() {
    if p.pollTicker != nil {
        p.pollTicker.Stop()
    }
    close(p.stopChan)
    global.GVA_LOG.Info("CoAP客户端轮询服务已停止")
}

// startPolling 开始轮询
func (p *coapClientPoller) startPolling() {
    // 启动时立即执行一次
    p.pollAllGridPoints()
    
    // 创建定时器
    p.pollTicker = time.NewTicker(p.pollInterval)
    
    for {
        select {
        case <-p.pollTicker.C:
            p.pollAllGridPoints()
        case <-p.stopChan:
            return
        }
    }
}

// pollAllGridPoints 轮询所有并网点的数据
func (p *coapClientPoller) pollAllGridPoints() {
    // 从数据库获取所有并网点配置
    var settings []agvc.AgvcBwdSetting
    if err := global.GVA_DB.Find(&settings).Error; err != nil {
        global.GVA_LOG.Error("获取并网点配置失败", zap.Error(err))
        return
    }
    
    if len(settings) == 0 {
        global.GVA_LOG.Debug("没有配置并网点，跳过轮询")
        return
    }
    
    global.GVA_LOG.Debug("开始轮询并网点数据", zap.Int("并网点数量", len(settings)))
    
    // 遍历每个并网点
    for _, setting := range settings {
        if setting.Number == nil {
            continue
        }
        
        var bwdNo int
        if _, err := fmt.Sscanf(*setting.Number, "%d", &bwdNo); err != nil {
            global.GVA_LOG.Warn("并网点编号格式错误",
                zap.String("number", *setting.Number),
                zap.Error(err))
            continue
        }
        
        // 轮询该并网点的三种设备类型数据
        p.pollGridPointData(bwdNo)
    }
}

// pollGridPointData 轮询单个并网点的所有设备类型数据
func (p *coapClientPoller) pollGridPointData(eqid int) {
    // 需要轮询的设备类型：AGC(60)、AVC(61)、并网柜(5)
    eqTypes := []int{cons.TYPE_AGC, cons.TYPE_AVC, cons.TYPE_BWG}
    
    for _, eqType := range eqTypes {
        if err := p.requestData(eqid, eqType); err != nil {
            global.GVA_LOG.Error("请求数据失败",
                zap.Int("EQID", eqid),
                zap.Int("EQType", eqType),
                zap.Error(err))
        }
    }
}

// requestData 向CoAP服务器请求数据
func (p *coapClientPoller) requestData(eqid, eqType int) error {
    // 创建CoAP客户端连接
    conn, err := udp.Dial(fmt.Sprintf("%s:%d", p.serverHost, p.serverPort))
    if err != nil {
        return fmt.Errorf("连接CoAP服务器失败: %v", err)
    }
    defer conn.Close()
    
    // 构建请求参数
    requestParams := map[string]int{
        "psid":   1,
        "eqid":   eqid,
        "eqType": eqType,
    }
    
    jsonData, err := json.Marshal(requestParams)
    if err != nil {
        return fmt.Errorf("序列化请求参数失败: %v", err)
    }
    
    // 创建请求上下文
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 发送GET请求到 /agvc/data
    resp, err := conn.Post(ctx, "/agvc/data", message.AppJSON, bytes.NewReader(jsonData))
    if err != nil {
        return fmt.Errorf("发送CoAP请求失败: %v", err)
    }
    
    // 检查响应状态码
    if resp.Code() != codes.Content && resp.Code() != codes.Valid {
        return fmt.Errorf("CoAP服务器返回错误状态: %s", resp.Code().String())
    }
    
    // 读取响应数据
    respBody, err := resp.ReadBody()
    if err != nil {
        return fmt.Errorf("读取响应数据失败: %v", err)
    }
    
    // 解析响应数据
    var dataMessages []request.CoAPDataMessage
    if err := json.Unmarshal(respBody, &dataMessages); err != nil {
        return fmt.Errorf("解析响应数据失败: %v", err)
    }
    
    // 存储数据到内存和InfluxDB
    if len(dataMessages) > 0 {
        p.storeData(dataMessages)
        
        global.GVA_LOG.Debug("成功获取并存储数据",
            zap.Int("EQID", eqid),
            zap.Int("EQType", eqType),
            zap.Int("数据量", len(dataMessages)))
    }
    
    return nil
}

// storeData 存储数据到内存和InfluxDB
func (p *coapClientPoller) storeData(messages []request.CoAPDataMessage) {
    // 转换为RealtimeData格式
    now := time.Now().Unix()
    dataList := make([]*agvc_main.RealtimeData, len(messages))
    
    for i, msg := range messages {
        dataList[i] = &agvc_main.RealtimeData{
            PSID:      msg.PSID,
            EQID:      msg.EQID,
            EQType:    msg.EQType,
            DataType:  msg.DataType,
            Point:     msg.Point,
            Value:     msg.Value,
            Timestamp: now,
        }
    }
    
    // 批量存储到内存
    DataStorage.StoreBatch(dataList)
    
    // 注意：DataStorage的periodicSave会自动将数据保存到InfluxDB
}

// SetServerAddress 设置CoAP服务器地址（可选，用于配置）
func (p *coapClientPoller) SetServerAddress(host string, port int) {
    p.serverHost = host
    p.serverPort = port
    global.GVA_LOG.Info("更新CoAP服务器地址",
        zap.String("server", fmt.Sprintf("%s:%d", host, port)))
}

// SetPollInterval 设置轮询间隔（可选，用于配置）
func (p *coapClientPoller) SetPollInterval(interval time.Duration) {
    p.pollInterval = interval
    if p.pollTicker != nil {
        p.pollTicker.Stop()
        p.pollTicker = time.NewTicker(interval)
    }
    global.GVA_LOG.Info("更新轮询间隔", zap.Duration("interval", interval))
}
