package service

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
