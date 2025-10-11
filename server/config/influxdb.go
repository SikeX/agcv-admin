package config

import "fmt"

type InfluxDB struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	Token    string `mapstructure:"token" json:"token" yaml:"token"`
	Org      string `mapstructure:"org" json:"org" yaml:"org"`
	Bucket   string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

func (i InfluxDB) Dsn() string {
	return fmt.Sprintf("http://%s:%s", i.Host, i.Port)
}
