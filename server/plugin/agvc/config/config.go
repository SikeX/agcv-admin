package config

// Config AGVC插件配置
type Config struct {
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`       // 是否启用插件
	CoapHost string `mapstructure:"coapHost" json:"coapHost" yaml:"coapHost"` // CoAP服务器地址
	CoapPort int    `mapstructure:"coapPort" json:"coapPort" yaml:"coapPort"` // CoAP服务器端口
}
