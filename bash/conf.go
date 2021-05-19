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
# dubbo client yaml configure file

check: true
# client
request_timeout: "3s"
# connect timeout
connect_timeout: "3s"

# application config
application:
  organization: "dubbo.io"
  name: "nCLinkGrpcServer"
  module: "nclink grpc server"
  version: "0.0.1"
  environment: "dev"

# registry config
registries:
  "zk":
    protocol: "zookeeper"
    timeout: "3s"
    address: "42.194.202.153:2181"
    zone: "guangzhou"

# reference config
references:
  "nCLinkServiceImpl":
    registry: "zk"
    protocol: "grpc"
    interface: "com.github.mk.NCLinkService"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    methods:
      - name: "NCLinkAuth"
        retries: 1
        loadbalance: "random"
      - name: "NCLinkSubscribe"
        retries: 0
        loadbalance: "random"
      - name: "NCLinkSendData"
        retries: 1
        loadbalance: "random"
      - name: "NCLinkSendBasicData"
        retries: 1
        loadbalance: "random"
      - name: "NCLinkGetMeta"
        retries: 1
        loadbalance: "random"

# protocol config
protocol_conf:
  grpc:
    reconnect_interval: 0
    connection_number: 2
    heartbeat_period: "5s"
    session_timeout: "20s"
    pool_size: 64
    pool_ttl: 600
    getty_session_param:
      compress_encoding: false
      tcp_no_delay: true
      tcp_keep_alive: true
      keep_alive_period: "120s"
      tcp_r_buf_size: 262144
      tcp_w_buf_size: 65536
      pkg_rq_size: 1024
      pkg_wq_size: 512
      tcp_read_timeout: "1s"
      tcp_write_timeout: "5s"
      wait_timeout: "1s"
      max_msg_len: 10240
      session_name: "client"
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
    "connect_method":"dubbo",
    "adaptor_id":"ADA-23"
}
`
