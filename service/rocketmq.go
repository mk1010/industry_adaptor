package service

import (
	
)

type RocketmqConfig struct {
	Name            string `json:"name"`
	Enable          bool   `json:"enable"`
	EnableSubscribe bool   `json:"enableSubscribe"`
	SubscribeTopic  string `json:"subscribeTopic"`
	SubscribeModel  string `json:"subscribeModel"`
	SubscribeTag    string `json:"subscribeTag"`
	NameSrv         string `json:"nameSrv"`
	GroupName       string `json:"groupName"`
}
