# ![rein](https://note.youdao.com/yws/public/resource/8bd89fcf7e10c7a878881b71865dcae4/xmlnote/E959C106CE854D6D825AD3E77B4AEB9F/27450) rein

| [English](https://github.com/firstboot/rein/blob/master/README.md) |  [中文](https://github.com/firstboot/rein/blob/master/README_zh.md) | 

| [Commit Log](https://github.com/firstboot/rein/blob/master/README_commit_log.md) |

#### Introduction

This program is used to efficiently redirect connections from one IP address/port combination to another. 

It is useful when operating virtual servers, dockers, firewalls and so on. 

It creates a tunnel from a public endpoint to a locally running service (version >= 1.0.5). It was developed by golang.



* `rein` support mode: 

| mode        | tip                                                          |
| ----------- | ------------------------------------------------------------ |
| `upstream`  | It is used to efficiently redirect connections from one IP address/port combination to another. |
| `fileshare` | local files publishing.                                     |
| `inps` | It creates a tunnel from a public endpoint to a locally running service. Put `inps` onto your server with public IP. (version >= 1.0.5) |
| `inpc` | It creates a tunnel from a public endpoint to a locally running service. Put `inpc` onto your server in LAN (that can't be connected from public Internet). (version >= 1.0.5) |



* `inps` and `inpc`  mode illustration: 

 ![mode-inps-inc](https://note.youdao.com/yws/public/resource/8bd89fcf7e10c7a878881b71865dcae4/xmlnote/32F23E8C5AF2447BBA7C124547326B17/27445)



----



### 1. Download by your OS

All release download: https://note.youdao.com/ynoteshare1/index.html?id=e11547282e63ce5920c5c2755a5cd89a&type=note

#### 1.1 Download

* *1.1.1 Download by CentOS/RHEL/Ubuntu*

```shell
cd ~
wget 	
http://note.youdao.com/yws/public/resource/e11547282e63ce5920c5c2755a5cd89a/xmlnote/18F3E51677BC41B3B1FE0F6B7DE359F5/27478 -O rein.zip
unzip rein.zip
chmod +x rein
```



* *1.1.2  Download by Windows*

  * download  file

    http://note.youdao.com/yws/public/resource/e11547282e63ce5920c5c2755a5cd89a/xmlnote/63FA93DEBA63475BA3DB18CB3574662B/27486

  * decompress `rein-x.x.x-amd64-windows.zip`



----

### 2. Simple deployment 

#### 2.1 Specify required functions

* *Which mode do you want to use ?*

| mode        | tip                                                          |
| ----------- | ------------------------------------------------------------ |
| `upstream`  | It is used to efficiently redirect connections from one IP address/port combination to another. |
| `fileshare` | local files publishing.                                      |
| `inps/inpc` | It creates a tunnel from a public endpoint to a locally running service. (server/client endpoint)  (version >= 1.0.5) |



#### 2.2 Function mode tip 

*Tip*:  OS description of operating differences, `CentOS/RHEL/Ubuntu` will be used by **default** in subsequent instructions.  *2.2.1* will still introduce them separately.

```shell
# show all mode, CentOS/RHEL/Ubuntu, eg: upstream
./rein -e-detail
./rein -e-detail-upstream

# generate config file, CentOS/RHEL/Ubuntu, eg: upstream
./rein -e-detail-upstream > rein.json

#####################################################

# show all mode, Windows, eg: upstream
./rein.exe -e-detail
./rein.exe -e-detail-upstream

# generate a mode conf file, eg: upstream
# Windows cmd 
./rein.exe -e-detail-upstream > rein.json

# Windows powershell, eg: upstream
./rein.exe -e-detail-upstream | out-file -encoding ascii rein.json

```



* *2.2.1 How to use mode `upstream`?*
  
  Use `-e-detail-xxx`  option, generate config and running.
  
  * *CentOS/RHEL/Ubuntu*
  
    * ```shell
      # show all mode
      ./rein -e-detail
      Enter a mode, show specific example, as follow:
      -e-detail-upstream
      -e-detail-inps
      -e-detail-inpc
      -e-detail-fileshare
      
      # generate a mode conf file, eg: upstream
      ./rein -e-detail-upstream > rein.json
      
      # rein.json, modify rein.json for you
      {
      	"upstream": [
      		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9990"}
      	]
      }
      
      # running
      ./rein -c rein.json
      ```
  
  * *Windows*
  
    * ```powershell
      # show all mode
      ./rein.exe -e-detail
      Enter a mode, show specific example, as follow:
      -e-detail-upstream
      -e-detail-inps
      -e-detail-inpc
      -e-detail-fileshare
      
      # generate a mode conf file, eg: upstream
      # windows cmd 
      ./rein.exe -e-detail-upstream > rein.json
      
      # windows powershell
      ./rein.exe -e-detail-upstream | out-file -encoding ascii rein.json
      
      # rein.json, modify rein.json for you
      {
      	"upstream": [
      		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9990"}
      	]
      }
      
    # running
      ./rein.exe -c rein.json
      ```
      
  
  
  
  
* *2.2.2 How to use mode `fileshare`?*

  ```shell
  # show default config
  ./rein -e-detail-fileshare
  {
  	"fileshare": [
  		{"port": "9990", "path": "."}
  	]
  }
  
  # generate a mode conf file
  ./rein -e-detail-fileshare > rein.json
  
  # rein.json, modify rein.json for you
  {
  	"fileshare": [
  		{"port": "9990", "path": "."}
  	]
  }
  
  # running
  ./rein -c rein.json
  ```

  

* *2.2.3 How to use mode `inps/inpc`?*

  If you have a server(A) public IP is `52.74.223.119` ,  and you have a server(B) private IP is `192.168.1.122`.

  Put `inps` onto your server A with public IP. Put `inpc` onto your server B in LAN (that can't be connected from public Internet). 

  *illustration*:

  ![example-inps-inpc](https://note.youdao.com/yws/public/resource/8bd89fcf7e10c7a878881b71865dcae4/xmlnote/877F17E2DC6C478892E82AD9BB29C0B2/27498)

  

  **deploy: inps**

  Put `inps` onto your server A with public IP.

  ```shell
  # generate a mode conf file
  ./rein -e-detail-inps > rein.json
  
  # rein.json, modify rein.json for you
  {
  	"inps": [
  		{"ctrl": "0.0.0.0:17500"}
  	]
  }
  
  # running
  ./rein -c rein.json
  ```

  

  **deploy: inpc**

  Put `inpc` onto your server B in LAN (that can't be connected from public Internet). 

  ```shell
  # generate a mode conf file
  ./rein -e-detail-inpc > rein.json
  
  # rein.json, modify rein.json for you
  # port 17500 is 'inps' server port
  # port 22 is local host service port, this port is ssh
  # port 9800 is 'inps' server open port
  {
  	"inpc": [
  		{
  			"ctrl": "52.74.223.119:17500",
  			"source": "0.0.0.0:9800",
  			"target": "127.0.0.1:22"
  		}
  	]
  }
  
  # running
  ./rein -c rein.json
  ```

  