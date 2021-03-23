package service

import (
	"github/mk1010/industry_adaptor/config"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var serviceInitOnce sync.Once

// 接口类型，将dubbo封装成MQTT的方法
var MeshClient MQTT.Client

func Init() (err error) {
	serviceInitOnce.Do(func() {
		switch config.ConfInstance.ConnectMethod {
		case "MQTT":
			err = mqttInit()
			if err != nil {
				return
			}
		case "DUBBO":
			// todo
		default:
			panic("Error ConnectMethod")
		}
	})
	return
}
