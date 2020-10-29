package app

import (
	"github.com/Sansui233/proxypool/pkg/provider"
	"github.com/Sansui233/proxypoolCheck/config"
	appcache "github.com/Sansui233/proxypoolCheck/internal/cache"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const version = "v0.3.11"

var router *gin.Engine

func setupRouter(){
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()          // 没有任何中间件的路由
	router.Use(gin.Recovery())  // 加上处理panic的中间件，防止遇到panic退出程序
	temp, err := loadHTMLTemplate() // 加载模板，模板源存放于html.go中的类似_assetsHtmlSurgeHtml的变量
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(temp) // 应用模板

	router.StaticFile("/css/index.css", "assets/css/index.css")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"domain":               config.Config.Domain,
			"request":              config.Config.Request,
			"port":                 config.Config.Port,
			"all_proxies_count":    appcache.AllProxiesCount,
			"ss_proxies_count":     appcache.SSProxiesCount,
			"ssr_proxies_count":    appcache.SSRProxiesCount,
			"vmess_proxies_count":  appcache.VmessProxiesCount,
			"trojan_proxies_count": appcache.TrojanProxiesCount,
			"useful_proxies_count": appcache.UsableProxiesCount,
			"last_crawl_time":      appcache.LastCrawlTime,
			"version":              version,
		})
	})
	router.GET("/clash", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clash.html", gin.H{
			"domain": config.Config.Domain,
			"port": config.Config.Port,
			"request": config.Config.Request,
		})
	})

	router.GET("/surge", func(c *gin.Context) {
		c.HTML(http.StatusOK, "surge.html", gin.H{
			"domain": config.Config.Domain,
			"request": config.Config.Request,
			"port": config.Config.Port,
		})
	})

	router.GET("/clash/config", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clash-config.yaml", gin.H{
			"domain": config.Config.Domain,
			"request": config.Config.Request,
			"port": config.Config.Port,
		})
	})
	router.GET("/clash/localconfig", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clash-config-local.yaml", gin.H{
			"port": config.Config.Port,
		})
	})
	router.GET("/clash/proxies", func(c *gin.Context) {
		proxyTypes := c.DefaultQuery("type", "")
		proxyCountry := c.DefaultQuery("c", "")
		proxyNotCountry := c.DefaultQuery("nc", "")
		text := ""
		if proxyTypes == "" && proxyCountry == "" && proxyNotCountry == "" {
			text = appcache.GetString("clashproxies")
			if text == "" {
				proxies := appcache.GetProxies("proxies")
				clash := provider.Clash{
					provider.Base{
						Proxies: &proxies,
					},
				}
				text = clash.Provide() // 根据Query筛选节点
				appcache.SetString("clashproxies", text)
			}
		} else {
			proxies := appcache.GetProxies("proxies")
			clash := provider.Clash{
				provider.Base{
					Proxies:    &proxies,
					Types:      proxyTypes,
					Country:    proxyCountry,
					NotCountry: proxyNotCountry,
				},
			}
			text = clash.Provide() // 根据Query筛选节点
		}
		c.String(200, text)
	})
	router.GET("/surge/proxies", func(c *gin.Context) {
		proxyTypes := c.DefaultQuery("type", "")
		proxyCountry := c.DefaultQuery("c", "")
		proxyNotCountry := c.DefaultQuery("nc", "")
		text := ""
		if proxyTypes == "" && proxyCountry == "" && proxyNotCountry == "" {
			text = appcache.GetString("surgeproxies")
			if text == "" {
				proxies := appcache.GetProxies("proxies")
				surge := provider.Surge{
					provider.Base{
						Proxies: &proxies,
					},
				}
				text = surge.Provide()
				appcache.SetString("surgeproxies", text)
			}
		} else {
			proxies := appcache.GetProxies("proxies")
			surge := provider.Surge{
				provider.Base{
					Proxies:    &proxies,
					Types:      proxyTypes,
					Country:    proxyCountry,
					NotCountry: proxyNotCountry,
				},
			}
			text = surge.Provide()
		}
		c.String(200, text)
	})
}

func Run() {
	setupRouter()
	port := config.Config.Port
	if port == "" {
		port = "8080"
	}
	// Run on this server
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err, "\n[router.go] Web server starting failed. Exit")
	}
}

// 返回页面templates
func loadHTMLTemplate() (t *template.Template, err error) {
	/* 使用本地模板文件 */
	filePaths, err := GetAllFilePaths("assets" + string(os.PathSeparator) + "html")
	if err != nil {
		log.Fatal("[router.go] Fail to load web templates: ", err)
		return nil, err;
	}
	for _, filePath := range filePaths {
		t, _ = t.ParseFiles(filePath) // Parsefile后的模板无路径前缀
		if err != nil {
			log.Panic("[router.go] ", err)
		}
	}
	return t, nil
}

// unix directory format
// TODO: This function shouldn't be here
func GetAllFilePaths(pathname string) (filenames []string,err error) {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFilePaths(pathname + string(os.PathSeparator) + fi.Name())
		} else {
			filename := pathname + string(os.PathSeparator) + fi.Name()
			filenames = append(filenames, filename)
		}
	}
	return filenames,err
}