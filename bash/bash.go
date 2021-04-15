package bash

import (
	"os"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func init() {
	err := os.Setenv("prod", "true")
	if err != nil {
		panic(err)
	}
	zapLogConf := new(zap.Config)
	err = jsoniter.Unmarshal([]byte(logConf), zapLogConf)
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
}
