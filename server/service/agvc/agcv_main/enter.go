package agcv_main

var Service = new(service)

type service struct {
    DataStorage  dataStorage
    CoapReceiver coapReceiver
    CoapSender   coapSender
    Device       device
    History      history
    AGC          agc
    AVC          avc
}

// init 初始化函数
func init() {
    // 初始化逆变器品牌模块
    InitInverterBrands()
}
