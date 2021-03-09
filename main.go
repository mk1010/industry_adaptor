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
	"context"
	"github/mk1010/industry_adaptor/dubbo_pkg"
	"os"
	"time"

	hessian "github.com/apache/dubbo-go-hessian2"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	gxlog "github.com/dubbogo/gost/log"

	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	dubboConfig "github.com/apache/dubbo-go/config"

	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var userProvider = new(dubbo_pkg.UserProvider)

func init() {
	dubboConfig.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&dubbo_pkg.User{})
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	hessian.RegisterPOJO(&dubbo_pkg.User{})
	dubboConfig.Load()
	time.Sleep(3 * time.Second)
	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &dubbo_pkg.User{}
	err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)
}