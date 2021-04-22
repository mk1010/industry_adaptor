/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mk1010/industry_adaptor/bash"
	"github.com/mk1010/industry_adaptor/service"
	"github.com/mk1010/industry_adaptor/task"

	"github.com/apache/dubbo-go/common/logger"
)

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	logger.Info("service start \n")
	err := service.Init()
	if err != nil {
		panic(err)
	}
	task.Init()
	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(3*time.Second), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			logger.Info("provider app exit now...")
			return
		}
	}
}
