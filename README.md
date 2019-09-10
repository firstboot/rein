# rein

![rein](https://raw.githubusercontent.com/firstboot/rein/master/rein-logo.png)

This program is used to efficiently redirect connections from one IP address/port combination to another. 

It is useful when operating virtual servers, dockers, firewalls and the like.  It was developed by golang.



### Alter tip:

#### 1.0.3

1. rein support `stream` mode.
2. rein support `fileshare` mode.



#### 1.0.4

1. The `fileshare` mode support multiple path, 1.0.3 only support single path.

   eg:

   ```
   {
   	"stream": [
   		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
   	],
   	"fileshare": [
   		{"port": "9991", "path": "/home/user/dir1"},
   		{"port": "9992", "path": "/home/user/dir2"}
   	]
   }
   ```

   

----



### 1. Simple deployment 

#### 1.1 CentOS/RHEL/Ubuntu

release download: https://note.youdao.com/ynoteshare1/index.html?id=b1e1ad270ba1b1af97ebdf3e2c8b7403&type=note

```shell
cd ~
wget http://note.youdao.com/yws/public/resource/b1e1ad270ba1b1af97ebdf3e2c8b7403/xmlnote/82E2CC3FF2744238B6AF36346298E5E5/27082 -O rein.zip
unzip rein.zip
chmod +x rein
./rein -e > rein.json
# modify rein.json for you
./rein -c rein.json
```

#### 1.2 Windows

- download  file

  https://note.youdao.com/ynoteshare1/index.html?id=b1e1ad270ba1b1af97ebdf3e2c8b7403&type=note

  `rein-amd64-windows.zip`

- decompress `rein-amd64-windows.zip`

- generating and modifying conf `rein.json` 

  ```powershell
  # generate default conf 'rein.json'
  # windows cmd 
  ./rein.exe -e > rein.json
  
  # windows powershell
  ./rein.exe -e | out-file -encoding ascii rein.json
  ```

  ```json
  {
  	"stream": [
  		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
  	],
  	"fileshare": [
  		{"port": "9991", "path": "."}
  	]
  }
  ```

- running

  ```powershell
  ./rein.exe -c rein.json
  ```

  
### 2. Configuration tip

#### 2.1 stream

`stream ` include keywords `source` and `target`. 

  `source`  open port  to listen, `target`  is data stream destination.

eg:

```json
{
	"stream": [
		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
	]
}
```

#### 2.2 file share

This function looks like `ftp`.

`fileshare` include keywords `port` and `path`.

eg: 

```
{
	"fileshare": [
		{"port": "9991", "path": "/home/lz"}
	]
}
```



----

`rein` 中文版说明



本程序主要用于进行反向代理IP地址和端口，功能类似于 `nginx` 的 `stream` 模式和`rinetd` 的功能，由于`rein`使用了`golang`语言开发，并且提供已经编译好的可下载版本，在部署配置方面比它们要方便些。

功能列表：

1. 反向代理`IP`和端口。
2. 提供本地文件的快速网络（`http`模式）分享。



### 修改说明:

#### 1.0.3

1. rein 支持 `stream` 模式。
2. rein 支持`fileshare` 模式。



#### 1.0.4

1. rein 的 `fileshare` 模式支持多路径分享, 在1.0.3 版本中支持一条路径。

   eg:

   ```
   {
   	"stream": [
   		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
   	],
   	"fileshare": [
   		{"port": "9991", "path": "/home/user/dir1"},
   		{"port": "9992", "path": "/home/user/dir2"}
   	]
   }
   ```

   



### 1. 简单快速部署

#### 1.1 CentOS/RHEL/Ubuntu 平台

已经编译好的版本下载地址： <https://note.youdao.com/ynoteshare1/index.html?id=b1e1ad270ba1b1af97ebdf3e2c8b7403&type=note>

下载 `rein-amd64-linux-x.x.x.zip `

如果您的 Linux 具备公网下载功能，可以直接通过下面的命令进行下载使用：

```shell
cd ~
wget http://note.youdao.com/yws/public/resource/f3c6a039b3a7ccee868fa50601663b44/xmlnote/D46BC1F68A334753AB615B3049D09F39/27313 -O rein.zip
# 需要安装 unzip 
unzip rein.zip
chmod +x rein
./rein -e > rein.json
# modify rein.json for you
./rein -c rein.json
```

#### 1.2 Windows 平台 

使用您的浏览器下载 <https://note.youdao.com/ynoteshare1/index.html?id=b1e1ad270ba1b1af97ebdf3e2c8b7403&type=note> 

`rein-amd64-windows-x.x.x.zip`并解压它。

使用下面的命令生成并修改 `rein.json` 配置文件

```powershell
# generate default conf 'rein.json'
# 使用 cmd 时
./rein.exe -e > rein.json

# 使用 powershell 时
./rein.exe -e | out-file -encoding ascii rein.json
```

生成的默认配置文件如下：

```json
{
	"upstream": [
		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
	],
	"fileshare": [
		{"port": "9991", "path": "."}
	]
}
```

根据您的需要进行修改配置文件后，运行：

```powershell
./rein.exe -c rein.json
```



### 2. 配置文件说明

#### 2.1 stream 模式

`stream` 模式主要由 `source` 和 `target` 构成，实现的功能就是将主机上的某个IP地址与端口，映射到其他的主机（本机）和端口上。在 `stream` 模式下，支持多组由 `source` 和 `target` 构成的映射对。`source` 是监听 IP 和端口，`target`是需要转发到的 IP 和端口。

举例说明：

```json
{
	"upstream": [
		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
	]
}
```

#### 2.2 file share 模式

此模式类似于`ftp`功能，能快速将本地资源进行网络（`http`方式）发布，它由`port`和`path`构成。`port`是要开放的端口，`path`是本地资源的路径。类似地，这个功能也支持多组。

举例说明：

```json
{
	"fileshare": [
		{"port": "9991", "path": "/home/lz"}
	]
}
```










