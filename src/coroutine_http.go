package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"./netopt"
)

// channel element struct
type chanHTTPEle struct {
	flag int
	n    int
	bf   string
}

type connChanEle struct {
	clientConn *net.TCPConn
}

type coroutineHTTPObj struct {
	bufferLen     int
	rootTargetMap map[string]string
	// rootFilterOldMap map[string]string
	// rootFilterNewMap map[string]string
}

// func coroutineHTTP(
// 	rootTargetMap map[string]string,
// 	rootFilterOldMap map[string]string,
// 	rootFilterNewMap map[string]string
// 	) coroutineHTTPObj {
// 	return coroutineHTTPObj{2048, rootTargetMap, rootFilterOldMap, rootFilterNewMap}
// }

func coroutineHTTP(rootTargetMap map[string]string) coroutineHTTPObj {
	return coroutineHTTPObj{2048, rootTargetMap}
}

func (obj coroutineHTTPObj) acceptDeal(netListen net.Listener) {
	for {
		log.Println("rein http proxy netListen.Accept ...")
		conn, err := netListen.Accept()
		if err != nil {
			log.Println("rein http proxy netListen.Accept error!")
			// os.Exit(1)
			return
		}
		log.Println("rein http proxy netListen.Accept ok!, conn id: ", fmt.Sprintf("%0x", &conn))
		go obj.communicationDeal(conn, obj.bufferLen)
	}
}

func (obj coroutineHTTPObj) communicationDeal(orgiConn net.Conn, bufferLen int) {
	// clientConn := coroutineStream().getClientConn(destServerAddr)

	channel := make(chan chanHTTPEle)
	go obj.orgiConnReadProducter(orgiConn, bufferLen, channel, channel)

}

func (obj coroutineHTTPObj) consumerDeal(orgiConn net.Conn, clientConn *net.TCPConn, channel <-chan chanHTTPEle) {
	for {
		ce := <-channel
		strLen := strconv.Itoa(len(ce.bf))
		log.Println("flag(enter): ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen)
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
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			return
		}
	}
}

func (obj coroutineHTTPObj) clientConnReadProducter(clientConn *net.TCPConn, bufferLen int, channel chan<- chanHTTPEle) {
	for {
		var bufferStr = ""
		var err error
		var n int
		for {
			buffer := make([]byte, bufferLen)
			n, err = clientConn.Read(buffer)
			bufferStr = bufferStr + string(buffer[:n])

			if err != nil {
				break
			}

			if strings.Index(bufferStr, "\r\n") >= 0 {
				break
			}

		}
		ce := chanHTTPEle{1, len(bufferStr), bufferStr}

		// buffer := make([]byte, bufferLen)
		// n, err = clientConn.Read(buffer)
		// // utilsStream().printBytes(buffer, n)
		// log.Println(string(buffer[:n]))
		// ce := chanHTTPEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanHTTPEle{-1, 0, ""}
			channel <- ce
			return
		}

	}
}

func (obj coroutineHTTPObj) orgiConnReadProducter(orgiConn net.Conn, bufferLen int,
	channel chan<- chanHTTPEle, coChannel chan chanHTTPEle) {
	for {
		bufferStr := ""
		var err error
		var n int
		for {
			buffer := make([]byte, bufferLen)
			n, err = orgiConn.Read(buffer)
			bufferStr = bufferStr + string(buffer[:n])
			// request header end condition
			if strings.Index(bufferStr, "\r\n\r\n") > 0 || strings.Index(bufferStr, "\n\n") > 0 {
				bufferStr = bufferStr + string(buffer[:n])
				break
			}
		}

		// fmt.Println(bufferStr)
		afterBufferStr, clientConn := obj.routeByRequestURL(bufferStr)
		if nil == clientConn {
			orgiConn.Close()
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			return
		}

		go obj.clientConnReadProducter(clientConn, bufferLen, coChannel)
		go obj.consumerDeal(orgiConn, clientConn, coChannel)

		ce := chanHTTPEle{0, len(afterBufferStr), afterBufferStr}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanHTTPEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineHTTPObj) routeByRequestURL(headerContext string) (string, *net.TCPConn) {
	req := protocolParser().getHTTPRequestHeader(headerContext)
	// fmt.Println(req.RequestURI)
	urlFragments := strings.Split(req.RequestURI, "/")
	// for inx, urlFragment := range urlFragments {
	// 	fmt.Println(inx, ":", urlFragment)
	// }
	if len(urlFragments) <= 0 {
		return headerContext, nil
	}

	rootLocation := "/" + urlFragments[1]
	log.Println("rootLocation-before:", rootLocation, req.RequestURI)
	var postfixURL = ""
	if rootLocation == req.RequestURI {
		postfixURL = "/"
	} else {
		postfixURL = "/" + req.RequestURI[len(rootLocation)+1:]
	}

	log.Println("rootLocation-after:", rootLocation, "postfixURL:", postfixURL)
	if utilsConf().isExistKeyOfStrMap(rootLocation, obj.rootTargetMap) == false {
		return headerContext, nil
	}

	destServerAddr := obj.rootTargetMap[rootLocation]

	log.Println("rootLocation:", rootLocation, "postfixURL:", postfixURL)
	log.Println(obj.rootTargetMap)
	req.URL.Path = postfixURL
	req.Host = destServerAddr
	// log.Println("######################")
	log.Println("req.body: \n" + protocolParser().setHTTPRequestHeader(req))
	// log.Println("######################")

	clientConn := netopt.NetGetClientConn(destServerAddr)
	log.Println("routeByRequestURL - clientConn:", fmt.Sprintf("%0x", &clientConn))
	return protocolParser().setHTTPRequestHeader(req), clientConn
}

func (obj coroutineHTTPObj) run(sourceAddr string) {
	log.Println("rein http proxy start ...")
	log.Println("rein http proxy net.Listen ...")
	netListen, err := net.Listen("tcp", sourceAddr)
	if err != nil {
		log.Println("rein http proxy net.Listen error!")
		os.Exit(1)
	}
	log.Println("rein http proxy net.Listen ready...")
	go obj.acceptDeal(netListen)
	log.Println("rein http proxy net.Listen start...")

	for {
		time.Sleep(time.Second * 1)
		if mainRunFlag == false {
			netListen.Close()
		}
	}
}
