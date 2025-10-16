package config

type MQTT struct {
	Broker       string `mapstructure:"broker" json:"broker" yaml:"broker"`                         // MQTT broker地址
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               // MQTT端口
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 密码
	ClientID     string `mapstructure:"client-id" json:"client-id" yaml:"client-id"`               // 客户端ID
	Qos          byte   `mapstructure:"qos" json:"qos" yaml:"qos"`                                 // QoS级别
	KeepAlive    int    `mapstructure:"keep-alive" json:"keep-alive" yaml:"keep-alive"`            // 保持连接时间（秒）
	CleanSession bool   `mapstructure:"clean-session" json:"clean-session" yaml:"clean-session"`   // 清理会话
	Order        bool   `mapstructure:"order" json:"order" yaml:"order"`                           // 消息顺序
	WillTopic    string `mapstructure:"will-topic" json:"will-topic" yaml:"will-topic"`           // 遗嘱消息主题
	WillPayload  string `mapstructure:"will-payload" json:"will-payload" yaml:"will-payload"`     // 遗嘱消息内容
	WillQos      byte   `mapstructure:"will-qos" json:"will-qos" yaml:"will-qos"`                 // 遗嘱消息QoS
	WillRetained bool   `mapstructure:"will-retained" json:"will-retained" yaml:"will-retained"` // 遗嘱消息保留标志
}

func (m *MQTT) GetBrokerUrl() string {
	return "tcp://" + m.Broker + ":" + m.Port
}