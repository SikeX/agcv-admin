package config

import "fmt"

type InfluxDB struct {
	Host          string `mapstructure:"host" json:"host" yaml:"host"`
	Port          string `mapstructure:"port" json:"port" yaml:"port"`
	Token         string `mapstructure:"token" json:"token" yaml:"token"`
	Org           string `mapstructure:"org" json:"org" yaml:"org"`
	AgvcBucket    string `mapstructure:"agvcBucket" json:"agvcBucket" yaml:"agvcBucket"`          // agvc储存桶(1天保留)
	AgvcRetention string `mapstructure:"agvcRetention" json:"agvcRetention" yaml:"agvcRetention"` // agvc储存桶保留时间(1天)
	NbqBucket     string `mapstructure:"nbqBucket" json:"nbqBucket" yaml:"nbqBucket"`             // 逆变器储存桶(7天保留)
	NbqRetention  string `mapstructure:"nbqRetention" json:"nbqRetention" yaml:"nbqRetention"`    // 逆变器储存桶保留时间(7天)
	Username      string `mapstructure:"username" json:"username" yaml:"username"`
	Password      string `mapstructure:"password" json:"password" yaml:"password"`
}

func (i InfluxDB) Dsn() string {
	return fmt.Sprintf("http://%s:%s", i.Host, i.Port)
}

// GetAgvcBucket 获取agvc储存桶名称，如果未配置则返回默认值 agvc_data
func (i InfluxDB) GetAgvcBucket() string {
	if i.AgvcBucket == "" {
		return "agvc_data"
	}
	return i.AgvcBucket
}

// GetNbqBucket 获取逆变器储存桶名称，如果未配置则返回默认值 nbq_data
func (i InfluxDB) GetNbqBucket() string {
	if i.NbqBucket == "" {
		return "nbq_data"
	}
	return i.NbqBucket
}

// GetMeasurement 获取measurement名称，统一返回 "agvc"
func (i InfluxDB) GetMeasurement() string {
	return "agvc"
}
