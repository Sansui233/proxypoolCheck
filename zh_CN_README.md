# Proxypool 健康检查

## 信息

这是为[proxypool](https://github.com/sansui233/proxypool)的代理节点检查，并提供检查后可用的代理节点
所以，您应该有一个（或知道一个）可用的[proxypool](https://github.com/sansui233/proxypool)服务器。

Proxypool 健康检查最好是在本地（即您家里）部署，也可以在自己的中国大陆服务器上运行。

## 安装和运行

二选一

### 1. 用构建好的

从[releases](https://github.com/Sansui233/proxypoolCheck/releases)中下载

将下载的文件重命名为proxypoolcheck（可选）

不要忘了将文件添加755权限，否则无法运行

```
chmod +775 proxypoolcheck
```

您可以将config.yaml放在与proxypoolcheck文件同一文件夹内，或者-c 指定某一路径的config

```shell
./proxypoolCheck
# or
./proxypoolCheck -c PathToConfig
```

### 2. 自行构建

确保安装golang，然后下载源码
```sh
$ go get -u -v github.com/Sansui233/proxypoolCheck
```

运行
```shell script
$ go run main.go -c ./config/config.yaml
```

## 配置

基本的配置

```yaml
# proxypool远程服务器的地址，空白为http://127.0.0.1:8080
server_url:
  - https://example.proxypoolserver.com
  - https://example.proxypoolserver.com/clash/proxies?type=vmess


# 对于您的本地服务器
request: http   # http / https
domain:         # default: 127.0.0.1
port:           # default: 80

cron_interval: 15 # default: 15  minutes
show_remote_speed: true # default false

speedtest:      # default false
connection:     # default 5
timout:         # default 10
```

If your web server port is not the same as proxypoolCheck serving port, you should put web server port in configuration, and set an environment variable `PORT` for proxypoolCheck to serve. This will be really helpful when you are doing frp.

如果您的Web服务器端口与proxypoolCheck服务端口不同，应该将web服务器端口放在配置中，并且设置环境变量`PORT`以供proxypoolCheck服务。当您使用frp时，这将非常有帮助。

```
export PORT=ppcheckport
```

## 声明

本项目遵循 GNU General Public License v3.0 开源，在此基础上，所有使用本项目提供服务者都必须在网站首页保留指向本项目的链接

本项目仅限个人自己使用，**禁止使用本项目进行营利**和**做其他违法事情**，产生的一切后果本项目概不负责。

## Screenshots

![](doc/1.png)

![](doc/2.png)

by [wangwang-code](https://github.com/wangwang-code)
