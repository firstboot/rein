package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"./netopt"
)

type coroutineInpsObj struct {
	bufferLen int
}

func coroutineInps() coroutineInpsObj {
	return coroutineInpsObj{10240}
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

func (obj coroutineInpsObj) communicationDeal(userServConn net.Conn, bufferLen int, ctrlServConn net.Conn) {
	channel := make(chan netopt.ChanEle)
	go netopt.SrvConnReadLeft(userServConn, bufferLen, channel)
	go netopt.SrvConnReadRight(ctrlServConn, bufferLen, channel)
	netopt.SrvSrvConsumer(userServConn, ctrlServConn, channel)
}

func (obj coroutineInpsObj) run(ctrlAddr string) {

	ctrlServLis := netopt.NetServerListen(ctrlAddr)
	inpccConnPairs := make(map[string]int)

	lisPairs := make(map[string]net.Listener)

	for {
		log.Println("wait for ctrlServConn link in ...")
		// ctrlServConn := obj.serverAccept(ctrlServLis) // block
		ctrlServConn := netopt.NetServerAccept(ctrlServLis)
		log.Println("ctrlServConn link in ok ...")
		// msg := obj.connRecvDealOnce(ctrlServConn, obj.bufferLen)
		msg := netopt.NetConnRecvDealOnce(ctrlServConn, obj.bufferLen)

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

		userBackupLis, userLisErr := lisPairs[sourceAddr]
		var userServLis net.Listener
		if false == userLisErr {
			log.Println("new userServLis is " + sourceAddr)
			userServLis = netopt.NetServerListen(sourceAddr) // proxy server source
			lisPairs[sourceAddr] = userServLis
		} else {
			log.Println("rebuild userServLis is " + sourceAddr)
			userBackupLis.Close()
			userBackupLis = nil
			userServLis = netopt.NetServerListen(sourceAddr) // proxy server source
			lisPairs[sourceAddr] = userServLis
		}

		go func() {

			// userServLis := obj.serverListen("0.0.0.0:9800") // proxy server source
			obj.acceptDeal(userServLis, ctrlServConn) // block

			// fmt.Println("bottom")
			userServLis.Close()

		}()
	}

}
