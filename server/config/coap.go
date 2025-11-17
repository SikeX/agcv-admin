package config

type Coap struct {
    Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
    Host   string `mapstructure:"host" json:"host" yaml:"host"`
    Port   int    `mapstructure:"port" json:"port" yaml:"port"`
    
    // CoAP客户端配置（用于主动轮询数据）
    ClientEnable   bool   `mapstructure:"client-enable" json:"clientEnable" yaml:"client-enable"`
    ClientHost     string `mapstructure:"client-host" json:"clientHost" yaml:"client-host"`
    ClientPort     int    `mapstructure:"client-port" json:"clientPort" yaml:"client-port"`
    PollInterval   int    `mapstructure:"poll-interval" json:"pollInterval" yaml:"poll-interval"` // 轮询间隔（秒）
}
