package main

/*
rein
date: 2019-6-17
author: lz
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func acceptDeal(netListen net.Listener, destServerAddr string) {
	for {
		fmt.Println("netListen.Accept ...")
		conn, err := netListen.Accept()
		if err != nil {
			fmt.Println("netListen.Accept error!")
			os.Exit(1)
		}
		fmt.Println("netListen.Accept ok!, conn id: ", fmt.Sprintf("%0x", &conn))
		go communicationDeal(conn, 2048, destServerAddr)
	}
}

type chanEle struct {
	flag int
	n    int
	bf   string
}

func orgiConnReadProducter(orgiConn net.Conn, bufferLen int, channel chan<- chanEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		ce := chanEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			fmt.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func clientConnReadProducter(clientConn *net.TCPConn, bufferLen int, channel chan<- chanEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			fmt.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func consumerDeal(orgiConn net.Conn, clientConn *net.TCPConn, channel <-chan chanEle) {
	for {
		ce := <-channel
		strLen := strconv.Itoa(strings.Count(ce.bf, ""))
		fmt.Println("flag: ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen)
		if ce.flag == 0 {
			clientConn.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == 1 {
			orgiConn.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == -1 {
			clientConn.Close()
			orgiConn.Close()
			fmt.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			fmt.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			return
		}
	}
}

func communicationDeal(orgiConn net.Conn, bufferLen int, destServerAddr string) {
	clientConn := getClientConn(destServerAddr)

	channel := make(chan chanEle)
	go orgiConnReadProducter(orgiConn, bufferLen, channel)
	go clientConnReadProducter(clientConn, bufferLen, channel)
	go consumerDeal(orgiConn, clientConn, channel)
}

func getClientConn(destServer string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", destServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	fmt.Println("getClientConn ok")
	return conn
}

func isTail(buffer []byte, len int) bool {
	for inx := 0; inx < len; inx++ {
		// fmt.Printf("%0x ", buffer[inx])
		if buffer[inx] == 0 {
			return true
		}
	}
	return false
}

func proxyServer(sourceAddr string, targetAddr string) {
	destServerAddr := targetAddr
	fmt.Println("rein server start ...")
	fmt.Println("net.Listen ...")
	netListen, err := net.Listen("tcp", sourceAddr)
	if err != nil {
		fmt.Println("net.Listen error!")
		os.Exit(1)
	}
	fmt.Println("net.Listen ok")
	acceptDeal(netListen, destServerAddr)
}

func getConf(confPath string) map[string]interface{} {
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// fmt.Println(b)
	confStr := string(b)
	fmt.Println(confStr)

	var c map[string]interface{}
	jsonErr := json.Unmarshal([]byte(confStr), &c)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		os.Exit(1)
	}

	return c

}

func isExistKeyOfMap(key string, confMap map[string]interface{}) bool {
	for confKey, _ := range confMap {
		if confKey == key {
			return true
		}
	}
	return false
}

func runFileShare(port string, sharePath string) {
	http.Handle("/", http.FileServer(http.Dir(sharePath)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	//var confMap map[string]interface{}

	// eg: rein.json
	exampleStr := `
{
	"upstream": [
		{"source": "0.0.0.0:8150", "target": "127.0.0.1:9991"}
	],
	"fileshare": [
		{"port": "9991", "path": "."}
	]
}`

	if len(os.Args) > 1 {
		if os.Args[1] == "-v" {
			fmt.Println("version: rein 1.0.2")
			return
		}

		if os.Args[1] == "-h" {
			helpStr := `-h: help(detail see: https://github.com/firstboot/rein)
-v: version
-c: conf(eg: rein -c rein.json)
-e: conf example`
			fmt.Println(helpStr)
			return
		}

		if os.Args[1] == "-e" {
			fmt.Println(exampleStr)
			return
		}

		if os.Args[1] == "-c" {
			confPath := os.Args[2]
			confMap := getConf(confPath)

			if isExistKeyOfMap("fileshare", confMap) == true {
				for k, v := range confMap["fileshare"].([]interface{}) {
					// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
					fmt.Println(k, v)
					port := v.(map[string]interface{})["port"].(string)
					path := v.(map[string]interface{})["path"].(string)
					fmt.Println(port, path)
					go runFileShare(port, path)
				}
			}

			if isExistKeyOfMap("upstream", confMap) == true {
				for k, v := range confMap["upstream"].([]interface{}) {
					// fmt.Println(k, v.(map[string]interface{})["source"], v.(map[string]interface{})["target"])
					fmt.Println(k, v)
					source := v.(map[string]interface{})["source"].(string)
					target := v.(map[string]interface{})["target"].(string)
					fmt.Println(source, target)
					go proxyServer(source, target)
				}
			}
		}

	} else {
		infoStr := `author:  lz
e-mail:  linzhanggeorge@gmail.com
index:   https://github.com/firstboot/rein
help:    -h`
		fmt.Println(infoStr)
		return
	}

	for {
		time.Sleep(time.Second * 1)
	}
}
