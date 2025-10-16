package initialize

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

// MQTT 初始化MQTT客户端
func MQTT() {
	mqttConfig := global.GVA_CONFIG.MQTT
	if mqttConfig.Broker == "" || mqttConfig.Port == "" {
		global.GVA_LOG.Info("MQTT配置不完整，跳过初始化")
		return
	}

	// 生成默认客户端ID（如果未设置）
	clientId := mqttConfig.ClientID
	if clientId == "" {
		// 生成随机客户端ID
		rand.Seed(time.Now().UnixNano())
		clientId = fmt.Sprintf("gin-vue-admin-client-%d", rand.Intn(1000000))
	}

	// 创建MQTT客户端选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttConfig.GetBrokerUrl())
	opts.SetClientID(clientId)

	if mqttConfig.Username != "" {
		opts.SetUsername(mqttConfig.Username)
	}

	if mqttConfig.Password != "" {
		opts.SetPassword(mqttConfig.Password)
	}

	opts.SetKeepAlive(time.Duration(mqttConfig.KeepAlive) * time.Second)
	opts.SetCleanSession(mqttConfig.CleanSession)
	opts.SetOrderMatters(mqttConfig.Order)

	// 设置遗嘱消息
	if mqttConfig.WillTopic != "" {
		opts.SetWill(mqttConfig.WillTopic, mqttConfig.WillPayload, mqttConfig.WillQos, mqttConfig.WillRetained)
	}

	// 创建MQTT客户端
	client := mqtt.NewClient(opts)

	// 连接到MQTT Broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		global.GVA_LOG.Error("MQTT连接失败", zap.Error(token.Error()))
		return
	}

	global.GVA_MQTT = client
	global.GVA_LOG.Info("MQTT连接成功", zap.String("broker", mqttConfig.GetBrokerUrl()))
}

// GetMqttClient 获取MQTT客户端实例
func GetMqttClient() mqtt.Client {
	if client, ok := global.GVA_MQTT.(mqtt.Client); ok {
		return client
	}
	return nil
}

// Subscribe 订阅MQTT主题
func Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	client := GetMqttClient()
	if client == nil {
		return fmt.Errorf("MQTT客户端未初始化")
	}

	if token := client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	global.GVA_LOG.Info("MQTT订阅主题成功", zap.String("topic", topic))
	return nil
}

// Publish 发布MQTT消息
func Publish(topic string, qos byte, retained bool, payload interface{}) error {
	client := GetMqttClient()
	if client == nil {
		return fmt.Errorf("MQTT客户端未初始化")
	}

	token := client.Publish(topic, qos, retained, payload)
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	global.GVA_LOG.Info("MQTT发布消息成功", zap.String("topic", topic))
	return nil
}

// Unsubscribe 取消订阅MQTT主题
func Unsubscribe(topics ...string) error {
	client := GetMqttClient()
	if client == nil {
		return fmt.Errorf("MQTT客户端未初始化")
	}

	if token := client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	global.GVA_LOG.Info("MQTT取消订阅主题成功", zap.Strings("topics", topics))
	return nil
}

func SaveRealData(topic string) {
	// client := global.GVA_INFLUXDB
	Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// 处理MQTT消息
		// fmt.Println("Received message:", string(msg.Payload()))
		//储存到influxDB
		// 解析JSON数据
		var dataList []MQTTData
		if err := json.Unmarshal(msg.Payload(), &dataList); err != nil {
			global.GVA_LOG.Error("解析MQTT数据失败", zap.Error(err))
			return
		}
		if err := StoreDataToInfluxDB(dataList); err != nil {
			global.GVA_LOG.Error("存储数据到InfluxDB失败", zap.Error(err))
			return
		}

	})
}

type MQTTData struct {
	Code  string `json:"code"`
	Time  int64  `json:"time"`
	Value int    `json:"value"`
}

// StoreDataToInfluxDB 将数据存储到InfluxDB
func StoreDataToInfluxDB(dataList []MQTTData) error {
	// 获取InfluxDB写API
	writeAPI := GetInfluxDBWriteAPI()
	if writeAPI == nil {
		return fmt.Errorf("InfluxDB写API未初始化")
	}

	// 批量写入数据
	for _, data := range dataList {
		// 创建数据点
		point := influxdb2.NewPoint(
			"mqtt", // measurement名称
			map[string]string{
				"code": data.Code, // tags
			},
			map[string]interface{}{
				"value": data.Value, // fields
			},
			time.UnixMilli(data.Time), // timestamp
		)

		// 写入数据点
		writeAPI.WritePoint(point)
	}

	// 确保数据被刷新到数据库
	writeAPI.Flush()

	return nil
}
