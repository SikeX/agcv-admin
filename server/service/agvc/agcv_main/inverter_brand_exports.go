package agcv_main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main/inverter_brands"
)

// 导出inverter_brands包的公共变量，保持向后兼容
var (
	// BrandFactory 品牌工厂实例
	BrandFactory = inverter_brands.BrandFactory

	// InverterBrandMapper 品牌映射器实例
	InverterBrandMapper = inverter_brands.InverterBrandMapper

	// HuaweiController 华为控制器实例
	HuaweiController = inverter_brands.HuaweiController

	// HuaweiPointInit 华为点位初始化器实例
	HuaweiPointInit = inverter_brands.HuaweiPointInit
)

// InitInverterBrands 初始化逆变器品牌模块
// 设置CoAP适配器，使inverter_brands包可以使用CoapSender
func InitInverterBrands() {
	// 注入CoAP发送器适配器
	inverter_brands.SetCoapAdapter(CoapSender)
	
	// 初始化品牌映射器
	InverterBrandMapper.Initialize()
}

// CoapSender 实现 CoapSenderInterface 接口
// 这个结构体已经在 coap_sender.go 中定义，这里只需要确保它实现了接口
func (s *coapSender) GetDefaultCoapHost() string {
	return "127.0.0.1"
}

func (s *coapSender) GetDefaultCoapPort() int {
	return 5683
}

func (s *coapSender) GetDispatchBackCoapPort() int {
	return 1189
}
