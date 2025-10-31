package config

type Coap struct {
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Host   string `mapstructure:"host" json:"host" yaml:"host"`
	Port   int    `mapstructure:"port" json:"port" yaml:"port"`
}
