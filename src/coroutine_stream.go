package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// channel element struct
type chanEle struct {
	flag int
	n    int
	bf   string
}

type coroutineStreamObj struct {
	bufferLen int
}

func coroutineStream() coroutineStreamObj {
	return coroutineStreamObj{2048}
}

func (obj coroutineStreamObj) acceptDealEx(netListen net.Listener, destServerAddr string, bufferLen int) {
	for {
		log.Println("netListen.Accept ...")
		conn, err := netListen.Accept()
		if err != nil {
			log.Println("netListen.Accept error!")
			// os.Exit(1)
			return
		}
		log.Println("netListen.Accept ok!, conn id: ", fmt.Sprintf("%0x", &conn))
		go obj.communicationDeal(conn, bufferLen, destServerAddr)
	}
}

func (obj coroutineStreamObj) acceptDeal(netListen net.Listener, destServerAddr string) {
	obj.acceptDealEx(netListen, destServerAddr, obj.bufferLen)
}

func (obj coroutineStreamObj) getClientConn(destServer string) *net.TCPConn {
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

func (obj coroutineStreamObj) orgiConnReadProducter(orgiConn net.Conn, bufferLen int, channel chan<- chanEle) {
	for {

		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		ce := chanEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineStreamObj) clientConnReadProducter(clientConn *net.TCPConn, bufferLen int, channel chan<- chanEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineStreamObj) consumerDeal(orgiConn net.Conn, clientConn *net.TCPConn, channel <-chan chanEle) {
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

func (obj coroutineStreamObj) communicationDeal(orgiConn net.Conn, bufferLen int, destServerAddr string) {
	clientConn := obj.getClientConn(destServerAddr)

	channel := make(chan chanEle)
	go obj.orgiConnReadProducter(orgiConn, bufferLen, channel)
	go obj.clientConnReadProducter(clientConn, bufferLen, channel)
	go obj.consumerDeal(orgiConn, clientConn, channel)
}

func (obj coroutineStreamObj) run(sourceAddr string, targetAddr string) {
	destServerAddr := targetAddr
	log.Println("rein stream server start ...")

	netListen, err := net.Listen("tcp", sourceAddr)
	if err != nil {
		log.Println("rein stream net.Listen error!")
		os.Exit(1)
	}

	log.Println("rein stream net.Listen ready...")
	go obj.acceptDeal(netListen, destServerAddr)
	log.Println("rein stream net.Listen start...")

	for {
		time.Sleep(time.Second * 1)
		if mainRunFlag == false {
			netListen.Close()
		}
	}
}
