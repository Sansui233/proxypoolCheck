package config

import (
	"errors"
	"github.com/Sansui233/proxypool/pkg/tool"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
)

var configFilePath = "config.yaml"

// ConfigOptions is a struct that represents config files
type ConfigOptions struct {
	ServerUrl          []string `json:"server_url" yaml:"server_url"`
	Domain             string   `json:"domain" yaml:"domain"`
	Port               string   `json:"port" yaml:"port"`
	Request            string   `json:"request" yaml:"request"`
	CronInterval       uint64   `json:"cron_interval" yaml:"cron_interval"`
	ShowRemoteSpeed    bool     `json:"show_remote_speed" yaml:"show_remote_speed"`
	HealthCheckTimeout int      `json:"healthcheck_timeout" yaml:"healthcheck_timeout"`
	HealthCheckConnection int 	`json:"healthcheck_connection" yaml:"healthcheck_connection"`
	SpeedTest          bool     `json:"speedtest" yaml:"speedtest"`
	SpeedConnection    int      `json:"speed_connection" yaml:"speed_connection"`
	SpeedTimeout       int      `json:"speed_timeout" yaml:"speed_timeout"`
}

var Config ConfigOptions

// Parse Config file
func Parse(path string) error {
	if path == "" {
		path = configFilePath
	} else {
		configFilePath = path
	}
	fileData, err := ReadFile(path)
	if err != nil {
		return err
	}
	Config = ConfigOptions{}
	err = yaml.Unmarshal(fileData, &Config)
	if err != nil {
		return err
	}
	// set default
	if Config.ServerUrl == nil{
		return errors.New("config error: no server url")
	}
	if Config.Domain == ""{
		Config.Domain = "127.0.0.1"
	}
	if Config.Port == ""{
		Config.Port = "80"
	}
	if Config.CronInterval == 0{
		Config.CronInterval = 15
	}
	if Config.Request == ""{
		Config.Request = "http"
	}
	if Config.HealthCheckTimeout == 0{
		Config.HealthCheckTimeout = 5
	}
	if Config.HealthCheckConnection == 0{
		Config.HealthCheckConnection = 100
	}
	if Config.SpeedConnection == 0{
		Config.SpeedConnection = 15
	}
	if Config.SpeedTimeout == 0 {
		Config.SpeedTimeout = 10
	}
	return nil
}


// 从本地文件或者http链接读取配置文件内容
func ReadFile(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := tool.GetHttpClient().Get(path)
		if err != nil {
			return nil, errors.New("config file http get fail")
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		}
		return ioutil.ReadFile(path)
	}
}
