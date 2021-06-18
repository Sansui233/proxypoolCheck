package app

import (
	"fmt"
	"github.com/Sansui233/proxypool/pkg/healthcheck"
	"github.com/Sansui233/proxypool/pkg/provider"
	"github.com/Sansui233/proxypoolCheck/config"
	"github.com/Sansui233/proxypoolCheck/internal/cache"
	"log"
	"time"
)

var location, _ = time.LoadLocation("PRC")

// Get all usable proxies from proxypool server and set app vars
func InitApp() error{
	// Get proxies from server
	proxies, err := getAllProxies()
	if err != nil {
		log.Println("Get proxies error: ", err)
		cache.LastCrawlTime = fmt.Sprint(time.Now().In(location).Format("2006-01-02 15:04:05"), err)
		return err
	}
	proxies = proxies.Derive().Deduplication()
	cache.AllProxiesCount = len(proxies)

	// set cache variables
	cache.SSProxiesCount = proxies.TypeLen("ss")
	cache.SSRProxiesCount = proxies.TypeLen("ssr")
	cache.VmessProxiesCount = proxies.TypeLen("vmess")
	cache.TrojanProxiesCount = proxies.TypeLen("trojan")
	cache.LastCrawlTime = fmt.Sprint(time.Now().In(location).Format("2006-01-02 15:04:05"))
	log.Println("Number of proxies:", cache.AllProxiesCount)

	log.Println("Now proceeding health check...")

	// healthcheck settings
	healthcheck.DelayConn = config.Config.HealthCheckConnection
	healthcheck.DelayTimeout = time.Duration(config.Config.HealthCheckTimeout) * time.Second
	healthcheck.SpeedConn = config.Config.SpeedConnection
	healthcheck.SpeedTimeout = time.Duration(config.Config.SpeedTimeout) * time.Second

	proxies = healthcheck.CleanBadProxiesWithGrpool(proxies)
	log.Println("Usable proxy count: ", len(proxies))

	// Save to cache
	cache.SetProxies("proxies", proxies)
	cache.UsableProxiesCount = len(proxies)

	if config.Config.SpeedTest == true {
		healthcheck.SpeedTestAll(proxies)
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

