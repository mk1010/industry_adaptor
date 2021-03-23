package service

import (
	"encoding/json"
	"github/mk1010/industry_adaptor/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	resty "github.com/go-resty/resty/v2"
)

var mqttClient mqtt.Client

func mqttInit() error {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"logic_id": config.ConfInstance.LogicID,
			"password": config.ConfInstance.MQTTPassword,
		}).
		SetHeader("Accept", "application/json").
		Get("www.baidu.com/search_result")
		// todo
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Result().([]byte), resp)
	if err != nil {
		return err
	}
	mqttOpt := mqtt.NewClientOptions()
	// todo
	mqttOpt.Order = false
	mqttClient = mqtt.NewClient(mqttOpt)
	MeshClient = mqttClient
	return nil
}
