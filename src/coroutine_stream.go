package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"./netopt"
)

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

func (obj coroutineStreamObj) communicationDeal(srvConn net.Conn, bufferLen int, destServerAddr string) {
	cliConn := netopt.NetGetClientConn(destServerAddr)
	channel := make(chan netopt.ChanEle)
	go netopt.SrvConnReadLeft(srvConn, bufferLen, channel)
	go netopt.CliConnReadRight(cliConn, bufferLen, channel)
	go netopt.SrvCliConsumer(srvConn, cliConn, channel)
}

func (obj coroutineStreamObj) run(sourceAddr string, targetAddr string) {
	destServerAddr := targetAddr
	log.Println("rein stream server start ...")
	netListen := netopt.NetServerListen(sourceAddr)
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
