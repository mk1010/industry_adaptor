package bash

var logConf = `
level: "debug"
development: true
disableCaller: false
disableStacktrace: false
sampling:
encoding: "console"

# encoder
encoderConfig:
  messageKey: "message"
  levelKey: "level"
  timeKey: "time"
  nameKey: "logger"
  callerKey: "caller"
  stacktraceKey: "stacktrace"
  lineEnding: ""
  levelEncoder: "capital"
  timeEncoder: "iso8601"
  durationEncoder: "millis"
  callerEncoder: "short"
  nameEncoder: ""

outputPaths:
  - "stderr"
errorOutputPaths:
  - "stderr"
initialFields:
`

var consumerConf = `
{

}
`

var adaptor_conf = `
{
    "dbs":{
        "industry_identification_center":{
            "database":"device_info",
            "settings":"charset=utf8mb4&parseTime=True&loc=Local&timeout=2s&readTimeout=30s&writeTimeout=3s",
            "write":{
                "consul":"industry_identification_center.mysql.hust_cs_344_write",
                "username":"hust_cs_344_w",
                "password":"6CP5wkGGor0JhACza8UPIBnUUO82JM3B",
                "default_host_port":"42.194.213.217:3306"
            },
            "read":[
                {
                    "consul":"industry_identification_center.mysql.hust_cs_344_write",
                    "username":"hust_cs_344_w",
                    "password":"6CP5wkGGor0JhACza8UPIBnUUO82JM3B",
                    "default_host_port":"42.194.213.217:3306"
                }
            ]
        }
    },
    "redis_cluster_name":"",
    "redis_hosts":[

    ],
    "env":"prod",
    "connect_method":"dubbo"
}
`
