package bash

import (
	"os"

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
	err = jsoniter.Unmarshal([]byte(consumerConf), dubboConsumerConf)
	if err != nil {
		panic(err)
	}
	config.SetConsumerConfig(*dubboConsumerConf)
	consumerConf = ""
	err = jsoniter.Unmarshal([]byte(adaptor_conf), adaptor_config.ConfInstance)
	if err != nil {
		panic(err)
	}
	adaptor_conf = ""
}
