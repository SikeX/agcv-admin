package inverter_brands

// CoapSenderInterface 定义CoAP发送器接口
// 用于解耦inverter_brands包与agcv_main包的直接依赖
type CoapSenderInterface interface {
	SendInverterCommand(host string, port, psid, inverterNo int, commands map[int]interface{}) error
	GetDefaultCoapHost() string
	GetDefaultCoapPort() int
	GetDispatchBackCoapPort() int
}

// coapAdapter CoAP适配器实例
var coapAdapter CoapSenderInterface

// SetCoapAdapter 设置CoAP适配器
// 在初始化时由agcv_main包调用，注入CoapSender实例
func SetCoapAdapter(adapter CoapSenderInterface) {
	coapAdapter = adapter
}

// GetCoapAdapter 获取CoAP适配器
func GetCoapAdapter() CoapSenderInterface {
	return coapAdapter
}
