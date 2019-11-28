package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// channel element struct
type chanInpsEle struct {
	flag int
	n    int
	bf   string
}

type coroutineInpsObj struct {
	bufferLen int
}

func coroutineInps() coroutineInpsObj {
	return coroutineInpsObj{10240}
}

func (obj coroutineInpsObj) serverAccept(netListen net.Listener) net.Conn {

	log.Println("serverAccept accept ...")
	conn, err := netListen.Accept()
	if err != nil {
		log.Println("serverAccept error!")
		// os.Exit(1)
		return nil
	}
	log.Println("netListen.Accept ok!, conn id: ", fmt.Sprintf("%0x", &conn))
	return conn
}

func (obj coroutineInpsObj) serverListen(servAddr string) net.Listener {
	var netListen net.Listener
	var err error
	for i := 0; i < 3; i++ {
		netListen, err = net.Listen("tcp", servAddr)
		if err != nil {
			log.Println("rein inps net.Listen error!", err)
			netListen.Close()
			netListen = nil
			// os.Exit(1)
		} else {
			return netListen
		}
	}
	return nil
}

func (obj coroutineInpsObj) connRecvDealOnce(conn net.Conn, bufferLen int) string {

	buffer := make([]byte, bufferLen)
	n, err := conn.Read(buffer)

	if err == io.EOF {
		conn.Close()
		log.Println("conn:", fmt.Sprintf("%0x", &conn), " close.")
		return ""
	}

	bufferStr := string(buffer[:n])
	log.Println("conn:", fmt.Sprintf("%0x", &conn), len(bufferStr), " recv: ", bufferStr)
	return bufferStr
}

func (obj coroutineInpsObj) acceptDealEx(userServLis net.Listener, ctrlServConn net.Conn, bufferLen int) {
	var firstRunFlag = false
	// for {
	log.Println("userServLis.Accept ...")
	conn, err := userServLis.Accept()
	if err != nil {
		log.Println("userServLis.Accept error!", err)
		// os.Exit(1)
		return
	}
	if firstRunFlag == false {
		ctrlServConn.Write([]byte(string("ok")))
		firstRunFlag = true
	}
	log.Println("userServLis.Accept ok!, conn id: ", fmt.Sprintf("%0x", &conn))
	obj.communicationDeal(conn, bufferLen, ctrlServConn)
	// }
}

func (obj coroutineInpsObj) acceptDeal(userServLis net.Listener, ctrlServConn net.Conn) {
	obj.acceptDealEx(userServLis, ctrlServConn, obj.bufferLen)
}

func (obj coroutineInpsObj) getClientConn(destServer string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", destServer)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		os.Exit(1)
	}
	log.Println("getClientConn ok")
	return conn
}

func (obj coroutineInpsObj) orgiConnReadProducter(orgiConn net.Conn, bufferLen int, channel chan<- chanInpsEle) {
	for {

		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		ce := chanInpsEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanInpsEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpsObj) clientConnReadProducter(clientConn net.Conn, bufferLen int, channel chan<- chanInpsEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanInpsEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanInpsEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpsObj) consumerDeal(orgiConn net.Conn, clientConn net.Conn, channel <-chan chanInpsEle) {
	for {
		ce := <-channel
		strLen := strconv.Itoa(strings.Count(ce.bf, ""))
		log.Println("flag: ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen)
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

func (obj coroutineInpsObj) communicationDeal(userServConn net.Conn, bufferLen int, ctrlServConn net.Conn) {

	channel := make(chan chanInpsEle)
	go obj.orgiConnReadProducter(userServConn, bufferLen, channel)
	go obj.clientConnReadProducter(ctrlServConn, bufferLen, channel)
	obj.consumerDeal(userServConn, ctrlServConn, channel)
}

func (obj coroutineInpsObj) run(ctrlAddr string) {

	ctrlServLis := obj.serverListen(ctrlAddr)
	inpccConnPairs := make(map[string]int)

	for {
		log.Println("wait for ctrlServConn link in ...")
		ctrlServConn := obj.serverAccept(ctrlServLis) // block
		log.Println("ctrlServConn link in ok ...")
		msg := obj.connRecvDealOnce(ctrlServConn, obj.bufferLen)
		log.Println("connRecvDealOnce link in ok ...")

		// supoort inpq
		inpccConnPairs[msg] = 0
		if msg == "inpq" {
			for key, ele := range inpccConnPairs {
				log.Println("inpq: ", key, ", ", ele)
				ctrlServConn.Write([]byte(key + "\n"))
			}
			ctrlServConn.Close()
			continue
		}

		pos := strings.Index(msg, "/")
		sourceAddr := msg[:pos]
		userServLis := obj.serverListen(sourceAddr) // proxy server source

		go func() {

			// userServLis := obj.serverListen("0.0.0.0:9800") // proxy server source
			obj.acceptDeal(userServLis, ctrlServConn) // block

			// fmt.Println("bottom")
			userServLis.Close()

		}()
	}

}
