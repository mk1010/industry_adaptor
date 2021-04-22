package config

import (
	"fmt"
	"io/ioutil"
	"os"

	json "github.com/json-iterator/go"
)

type Config struct {
	DBConfigs        map[string]DBConfig `json:"dbs"`
	Env              string              `json:"env"`
	RedisClusterName string              `json:"redis_cluster_name"`
	RedisHosts       []string            `json:"redis_hosts"`
	ConnectMethod    string              `json:"connect_method"`
}

type DBConfig struct {
	Database string          `json:"database"`
	Settings string          `json:"settings"`
	WriteDB  DBConnectInfo   `json:"write"`
	ReadDB   []DBConnectInfo `json:"read"` // attention
}

type DBConnectInfo struct {
	AuthKey         string `json:"auth_key"` // gorm支持的另外一种认证方式
	Consul          string `json:"consul"`
	UserName        string `json:"username"`
	Password        string `json:"password"`
	DefaultHostPort string `json:"default_host_port"`
}

var ConfInstance = new(Config)

func isProduct() bool {
	return ConfInstance.Env == "prod"
}

func NewConfig(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func Init(file string) error {
	if ConfInstance != nil {
		return nil
	}
	conf, err := NewConfig(file)
	if err != nil {
		return err
	}
	if len(conf.Env) == 0 {
		conf.Env = "dev"
	}
	ConfInstance = conf
	return nil
}

func CheckEnv() string {
	if os.Getenv("PRODUCT_ENV") != "" {
		return "prod"
	}
	return "dev"
}

func GetConfigPath() string {
	return fmt.Sprintf("./conf/industry_identification_center_%s.json", CheckEnv())
}

// config.Input_ConfDir + "/" + fmt.Sprintf("industry_identification_center_%s.json", curEnv)
