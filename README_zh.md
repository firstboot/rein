# ![rein](https://note.youdao.com/yws/public/resource/8bd89fcf7e10c7a878881b71865dcae4/xmlnote/E959C106CE854D6D825AD3E77B4AEB9F/27450) rein

| [English](https://github.com/firstboot/rein/blob/master/README.md) |  [中文](https://github.com/firstboot/rein/blob/master/README_zh.md) | 

| [Commit Log](https://github.com/firstboot/rein/blob/master/README_commit_log.md) |

#### 介绍

本程序主要用于进行反向代理IP地址和端口，功能类似于 `nginx` 的 `stream` 模式和`rinetd` 的功能；在（1.0.5）版本开始，`rein`支持内网穿透，这一功能类似于`frp` 和`ngrok`。由于`rein`使用了`golang`语言开发，并且提供已经编译好的可下载版本，在部署配置方面比它们要方便些。

功能列表：

1. 反向代理`IP`和端口（`upstream` 模式）。
2. 提供本地文件的快速网络（`http`协议）分享（`fileshare` 模式）。
3. 内网穿透（`inps / inpc` 模式）



* `rein` 支持模式: 

| 模式        | 说明                                                         |
| ----------- | ------------------------------------------------------------ |
| `upstream`  | 反向代理模式                                                 |
| `fileshare` | 提供本地文件的快速网络（`http`协议）分享                     |
| `inps`      | 内网穿透的服务器端，`inps` 需要部署在有公网地址服务器上（版本 >= 1.0.5） |
| `inpc`      | 内网穿透的客户机端，`inpc` 部署在能访问互联网，没有公网 IP 地址的 PC 或服务器上（版本 >= 1.0.5） |



* `inps` 和 `inpc`  模式图解: 

 ![mode-inps-inc](https://note.youdao.com/yws/public/resource/8bd89fcf7e10c7a878881b71865dcae4/xmlnote/32F23E8C5AF2447BBA7C124547326B17/27445)



----

### 1. 根据你的操作系统下载指定可执行文件

所有操作系统版本的压缩包下载网址:  https://github.com/firstboot/rein/releases

#### 1.1 根据操作系统下载

* *1.1.1  在 CentOS/RHEL/Ubuntu 平台下载（64位）*

```shell
cd ~
wget \
https://github.com/firstboot/rein/releases/download/v1.0.6-bin/rein-1.0.6-amd64-linux.zip
unzip rein-1.0.6-amd64-linux.zip
mv rein-1.0.6-amd64-linux rein
chmod +x rein
```

* *1.1.2  在 Windows 平台下载（64位）*

  * 打开浏览器，在地址栏输入下面的地址后下载

    https://github.com/firstboot/rein/releases/download/v1.0.6-bin/rein-1.0.6-amd64-win.exe.zip

  * 解压压缩包 `rein-x.x.x-amd64-windows.zip`

----

### 2. 简单部署

#### 2.1 确定你需要的功能进行部署

* *三种模式说明*

| 功能模式        | 说明                                                          |
| ----------- | ------------------------------------------------------------ |
| `upstream`  | 反向代理 |
| `fileshare` | 提供本地文件的快速网络（`http`协议）分享                                     |
| `inps/inpc` | 内网穿透 (服务器端/客户端)  (版本 >= 1.0.5) |



#### 2.2 具体功能模式说明

*说明*:  不同操作系统的操作存在差异,  在后续的介绍中，`CentOS/RHEL/Ubuntu`，将会被作为默认操作。  不过 *2.2.1* 将仍然分开介绍。

下面是主要的差异操作部分：

```shell
# 在 CentOS/RHEL/Ubuntu下， 显示所有模式, 以 upstream 模式为例
./rein -e-detail
./rein -e-detail-upstream

# 在  CentOS/RHEL/Ubuntu 下，生成配置文件，以 upstream 模式为例
./rein -e-detail-upstream > rein.json

#####################################################

# 在 Windows 模式下， 显示所有模式, 以 upstream 模式为例
./rein.exe -e-detail
./rein.exe -e-detail-upstream

# 在  CentOS/RHEL/Ubuntu 下，生成配置文件，以 upstream 模式为例
# 使用 Windows cmd 
./rein.exe -e-detail-upstream > rein.json

# 使用 Windows powershell
./rein.exe -e-detail-upstream | out-file -encoding ascii rein.json

```



* *2.2.1 使用 `upstream` 模式*
  
  使用 `-e-detail-xxx`  选项，产生配置文件并运行。
  
  * *在 CentOS/RHEL/Ubuntu 下部署*
  
    * ```shell
      # 显示所有模式
      ./rein -e-detail
      Enter a mode, show specific example, as follow:
      -e-detail-upstream
      -e-detail-inps
      -e-detail-inpc
      -e-detail-fileshare
      
      # 生成配置文件，以 upstream 模式为例
      ./rein -e-detail-upstream > rein.json
      
      # 修改 rein.json 配置文件以符合你的需要
      {
      	"upstream": [
      		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9990"}
      	]
      }
      
      # 运行 rein
      ./rein -c rein.json
      ```
  
  * *在 Windows 下运行*
  
    * ```powershell
      # 显示所有模式
      ./rein.exe -e-detail
      Enter a mode, show specific example, as follow:
      -e-detail-upstream
      -e-detail-inps
      -e-detail-inpc
      -e-detail-fileshare
      
      # 生成配置文件，以 upstream 模式为例
      # 使用 windows cmd 
      ./rein.exe -e-detail-upstream > rein.json
      
      # 使用 windows powershell
      ./rein.exe -e-detail-upstream | out-file -encoding ascii rein.json
      
      # 修改 rein.json 配置文件以符合你的需要
      {
      	"upstream": [
      		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9990"}
      	]
      }
      
      # 运行 rein
      ./rein.exe -c rein.json
      ```
  
  
  
* *2.2.2 使用 `fileshare` 模式*

  ```shell
  # 显示默认的示例配置
  ./rein -e-detail-fileshare
  {
  	"fileshare": [
  		{"port": "9990", "path": "."}
  	]
  }
  
  # 生成配置文件
  ./rein -e-detail-fileshare > rein.json
  
  # 修改 rein.json 配置文件以符合你的需要
  {
  	"fileshare": [
  		{"port": "9990", "path": "."}
  	]
  }
  
  # 运行 rein
  ./rein -c rein.json
  ```

  

* *2.2.3 使用 `inps/inpc` 模式*

  举例：假设你有一台服务器（A），有一个公网IP地址是 `52.74.223.119` , 然后你还有一台服务器或PC是（B），它的IP地址是`192.168.1.122`（这是一个局域网IP，它没有公网IP，不过它可以访问互联网）.

  把 `inps` 部署在服务器（A），把 `inpc` 部署在服务器或PC（B）. 

  *用图来解释下*:

  ![example-inps-inpc](https://note.youdao.com/yws/public/resource/8bd89fcf7e10c7a878881b71865dcae4/xmlnote/877F17E2DC6C478892E82AD9BB29C0B2/27498)

  

  * **inps的部署命令**

  Put `inps` onto your server A with public IP.

  ```shell
  # 生成配置文件
  ./rein -e-detail-inps > rein.json
  
  # 修改 rein.json 配置文件以符合你的需要
  {
  	"inps": [
  		{"ctrl": "0.0.0.0:17500"}
  	]
  }
  
  # 运行 rein
  ./rein -c rein.json
  ```

  

  * **inpc的部署命令**

  ```shell
  # 生成配置文件
  ./rein -e-detail-inpc > rein.json
  
  # 修改 rein.json 配置文件以符合你的需要
  # 端口 17500 是部署在服务器（A）上开放的 'inps' 服务端口
  # 端口 22 是服务器（B）本地端口
  # 端口 9800 是通过 'inps' 服务开放的端口，直接映射到服务器（B）的22端口
  {
  	"inpc": [
  		{
  			"ctrl": "52.74.223.119:17500",
  			"source": "0.0.0.0:9800",
  			"target": "127.0.0.1:22"
  		}
  	]
  }
  
  # 运行 rein
  ./rein -c rein.json
  ```
  
  * **inpq的使用（1.0.6版本及以上）**
  
  ```shell
  # 通过 inpq 可以获取 inps 的连接状态
  ./rein -inpq x.x.x.x:17500
  0.0.0.0:9800/127.0.0.1:22, online
  ```
  
  

