package app

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/Sansui233/proxypool/pkg/proxy"
	"github.com/Sansui233/proxypoolCheck/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func getAllProxies() (proxy.ProxyList, error) {
	var proxylist proxy.ProxyList
	var errs []error // collect errors

	for _, value := range config.Config.ServerUrl {
		url := formatURL(value)
		pjson, err := getProxies(url)
		if err != nil {
			log.Printf("Error when fetch %s: %s\n", url, err.Error())
			errs = append(errs, err)
			continue
		}
		log.Printf("Get %s line count: %d\n", url, len(pjson))

		for i, p := range pjson {
			if i == 0 || len(p) < 2 {
				continue
			}
			p = p[2:] // remove "- "

			if pp, ok := convert2Proxy(p); ok {
				if i == 1 && pp.BaseInfo().Name == "NULL" {
					log.Println("no proxy on " + url)
					errs = append(errs, errors.New("no proxy on "+url))
					continue
				}
				if config.Config.ShowRemoteSpeed == true {
					name := strings.Replace(pp.BaseInfo().Name, " |", "_", 1)
					pp.SetName(name)
				}
				proxylist = append(proxylist, pp)
			}
		}
	}

	if proxylist == nil {
		if errs != nil {
			errInfo := "\n"
			for _, e := range errs {
				errInfo = errInfo + e.Error() + ";\n"
			}
			return nil, errors.New(errInfo)
		}
		return nil, errors.New("no proxy")
	}
	return proxylist, nil
}

func formatURL(value string) string {
	url := "http://127.0.0.1:8080"
	if value != "http://127.0.0.1:8080" {
		url = value
		if url[len(url)-1] == '/' {
			url = url[:len(url)-1]
		}
	}
	urls := strings.Split(url, "/")
	if urls[len(urls)-2] != "clash" {
		url = url + "/clash/proxies"
	}
	return url
}

// get proxy strings from url
func getProxies(url string) ([]string, error) {
	//resp, err := http.Get(url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: tr,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	proxyJson := strings.Split(string(body), "\n")
	if len(proxyJson) < 2 {
		return nil, errors.New("no proxy on " + url)
	}
	return proxyJson, nil
}

// Convert json string(clash format) to proxy
func convert2Proxy(pjson string) (proxy.Proxy, bool) {
	var f interface{}
	err := json.Unmarshal([]byte(pjson), &f)
	if err != nil {
		return nil, false
	}
	jsnMap := f.(interface{}).(map[string]interface{})

	switch jsnMap["type"].(string) {
	case "ss":
		var p proxy.Shadowsocks
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil {
			return nil, false
		}
		return &p, true
	case "ssr":
		var p proxy.ShadowsocksR
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil {
			return nil, false
		}
		return &p, true
	case "vmess":
		var p proxy.Vmess
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil {
			return nil, false
		}
		return &p, true
	case "trojan":
		var p proxy.Trojan
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil {
			return nil, false
		}
		return &p, true
	}
	return nil, false
}
