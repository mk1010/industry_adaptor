package service

import (
	"github/mk1010/industry_adaptor/config"
)

func Init() error {
	switch config.ConfInstance.ConnectMethod {
	case "MQTT":
		err := mqttInit()
		if err != nil {
			return err
		}
	case "HTTP2":
	default:
		panic("Error ConnectMethod")
	}
	return nil
}
