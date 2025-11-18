package config

import "fmt"

type InfluxDB struct {
	Host        string `mapstructure:"host" json:"host" yaml:"host"`
	Port        string `mapstructure:"port" json:"port" yaml:"port"`
	Token       string `mapstructure:"token" json:"token" yaml:"token"`
	Org         string `mapstructure:"org" json:"org" yaml:"org"`
	Bucket      string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Username    string `mapstructure:"username" json:"username" yaml:"username"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	Measurement string `mapstructure:"measurement" json:"measurement" yaml:"measurement"` // 数据测量表名，默认 agvc_data

}

func (i InfluxDB) Dsn() string {
	return fmt.Sprintf("http://%s:%s", i.Host, i.Port)
}

// GetMeasurement 获取measurement名称，如果未配置则返回默认值 agvc_data
func (i InfluxDB) GetMeasurement() string {
	if i.Measurement == "" {
		return "agvc_data"
	}
	return i.Measurement
}
