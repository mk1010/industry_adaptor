package bash

import (
	"net"
	"os"
	"strings"

	adaptor_config "github.com/mk1010/industry_adaptor/config"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func init() {
	err := os.Setenv("prod", "true")
	if err != nil {
		panic(err)
	}
	zapLogConf := new(zap.Config)
	err = yaml.Unmarshal([]byte(logConf), zapLogConf)
	if err != nil {
		panic(err)
	}
	logger.InitLogger(zapLogConf)
	logConf = "" // 释放内存
	dubboConsumerConf := new(config.ConsumerConfig)
	err = yaml.Unmarshal([]byte(consumerConf), dubboConsumerConf)
	if err != nil {
		panic(err)
	}
	if m, ok := dubboConsumerConf.ProtocolConf.(map[string]interface{}); ok {
		intMap := make(map[interface{}]interface{})
		for k, v := range m {
			intMap[k] = v
		}
		dubboConsumerConf.ProtocolConf = intMap
	}
	config.SetConsumerConfig(*dubboConsumerConf)
	consumerConf = ""
	err = jsoniter.Unmarshal([]byte(adaptor_conf), adaptor_config.ConfInstance)
	if err != nil {
		panic(err)
	}
	adaptor_conf = ""
}

func getPublicIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		logger.Error(err)
		return ""
	}
	defer conn.Close()
	s := conn.LocalAddr().String()
	if index := strings.LastIndex(s, ":"); index != -1 {
		s = s[:index]
	}
	return s
}
