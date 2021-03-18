package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sansui233/proxypool/pkg/healthcheck"
	"github.com/Sansui233/proxypool/pkg/provider"
	"github.com/Sansui233/proxypool/pkg/proxy"
	"github.com/Sansui233/proxypoolCheck/config"
	"github.com/Sansui233/proxypoolCheck/internal/cache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var location, _ = time.LoadLocation("PRC")

// Get all usable proxies from proxypool server and set app vars
func InitApp() error{
	// Get proxies from server
	proxies, err := getAllProxies()
	if err != nil {
		log.Println("Get proxies error: ", err)
		cache.LastCrawlTime = fmt.Sprint(time.Now().In(location).Format("2006-01-02 15:04:05"),err)
		return err
	}
	cache.AllProxiesCount = len(proxies)

	// set cache variables
	cache.SSProxiesCount = proxies.TypeLen("ss")
	cache.SSRProxiesCount = proxies.TypeLen("ssr")
	cache.VmessProxiesCount = proxies.TypeLen("vmess")
	cache.TrojanProxiesCount = proxies.TypeLen("trojan")
	cache.LastCrawlTime = fmt.Sprint(time.Now().In(location).Format("2006-01-02 15:04:05"))

	log.Println("Number of proxies:", cache.AllProxiesCount)
	log.Println("Now proceeding health check...")
	proxies = healthcheck.CleanBadProxiesWithGrpool(proxies)
	log.Println("Usable proxy count: ", len(proxies))
	// Save to app cache
	cache.SetProxies("proxies", proxies)
	cache.UsableProxiesCount = len(proxies)

	// speedtest
	if config.Config.SpeedTest == true{
		if config.Config.Timeout >= 0 {
			healthcheck.SpeedTimeout = time.Duration(config.Config.Timeout) * time.Second
		}
		healthcheck.SpeedTestAll(proxies, config.Config.Connection)
	}
	cache.SetString("clashproxies", provider.Clash{
		provider.Base{
			Proxies: &proxies,
		},
	}.Provide())
	cache.SetString("surgeproxies", provider.Surge{
		provider.Base{
			Proxies: &proxies,
		},
	}.Provide())

	fmt.Println("Open", config.Config.Domain+":"+config.Config.Port, "to check.")
	return nil
}

func getAllProxies() (proxy.ProxyList, error){
	var proxylist proxy.ProxyList
	tempMap := map[string]byte{}
	tempMapLen:=len(tempMap)

	for _,value:=range config.Config.ServerUrl{
		url := "http://127.0.0.1:8080"
		if value != "http://127.0.0.1:8080"{
			url = value
			if url[len(url)-1] == '/' {
				url = url[:len(url)-1]
			}
		}
		urls := strings.Split(url,"/")
		if urls[len(urls)-2] != "clash" {
			url = url + "/clash/proxies"
		}
		log.Printf("Requesting %v",url)
		resp, err := http.Get(url)
		if err != nil{
			log.Printf("warning: %v\n", err)
			continue
			//return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		pjson := strings.Split(string(body),"\n")
		if len(pjson) < 2{
			log.Println("warning: No proxy on remote server")
			continue
			//return nil, errors.New("No proxy on remote server")
		}


		for i, pstr := range pjson {
			if i == 0 || len(pstr)<2{
				continue
			}
			tempMap[pstr]=0
			if len(tempMap) ==tempMapLen{
				continue
			}
			tempMapLen++
			pstr = pstr[2:]
			if pp, ok := convert2Proxy(pstr); ok{
				if i == 1 && pp.BaseInfo().Name == "NULL" {
					log.Println("warning: No proxy on remote server")
					continue
					//return nil, errors.New("No proxy on remote server")
				}
				if config.Config.ShowRemoteSpeed == true {
					name := strings.Replace(pp.BaseInfo().Name, " |", "_",1)
					pp.SetName(name)
				}
				proxylist = append(proxylist, pp)
			}
		}
	}

	if proxylist == nil {
		return nil, errors.New("No Proxy")
	}
	return proxylist, nil
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
		if err != nil{
			return nil, false
		}
		return &p, true
	case "ssr":
		var p proxy.ShadowsocksR
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil{
			return nil, false
		}
		return &p, true
	case "vmess":
		var p proxy.Vmess
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil{
			return nil, false
		}
		return &p, true
	case "trojan":
		var p proxy.Trojan
		err := json.Unmarshal([]byte(pjson), &p)
		if err != nil{
			return nil, false
		}
		return &p, true
	}
	return nil, false
}
