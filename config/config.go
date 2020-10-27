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
	ServerUrl string `json:"server_url" yaml:"server_url"`
	Request	  string `json:"request" yaml:"request"'`
	Domain    string `json:"domain" yaml:"domain"`
	Port      string `json:"port" yaml:"port"`
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
	if Config.ServerUrl == ""{
		Config.ServerUrl = "http://127.0.0.1:8080"
	}
	if Config.Domain == ""{
		Config.Domain = "127.0.0.1"
	}
	if Config.Port == ""{
		Config.Port = "8080"
	}
	if Config.Request == ""{
		Config.Request = "http"
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