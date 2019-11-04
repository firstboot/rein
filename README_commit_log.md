# rein alter tip



### Alter tip:

#### 1.0.5

*release: 20191104*

1. rein support `inps` and `inpc` mode.

2. what is `inps` and `inpc` mode?

   `inp`  means `internal network penetration` . It reverse proxy to help you expose a local server behind a NAT or firewall to the Internet. It creates a tunnel from a public endpoint to a locally running service.

   

----

#### 1.0.4

*release:  20190726*

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

#### 1.0.3

*release:  20190625*

1. rein support `stream` mode.
2. rein support `fileshare` mode.



----

#### 1.0.2  

*release: 20190617*

beta version



----

#### 1.0.1

*release: 20190520*

test version







