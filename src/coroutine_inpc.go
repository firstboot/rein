package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"./netopt"
)

type coroutineInpcObj struct {
	bufferLen int
}

func coroutineInpc() coroutineInpcObj {
	return coroutineInpcObj{2048}
}

func (obj coroutineInpcObj) acceptDealEx(ctrlCliConn *net.TCPConn, userCliConn *net.TCPConn, bufferLen int) {
	obj.communicationDeal(ctrlCliConn, bufferLen, userCliConn)

}

func (obj coroutineInpcObj) acceptDeal(ctrlCliConn *net.TCPConn, userCliConn *net.TCPConn) {
	obj.acceptDealEx(ctrlCliConn, userCliConn, obj.bufferLen)
}

func (obj coroutineInpcObj) communicationDeal(ctrlCliConn *net.TCPConn, bufferLen int, userCliConn *net.TCPConn) {
	channel := make(chan netopt.ChanEle)
	go netopt.CliConnReadLeft(ctrlCliConn, bufferLen, channel)
	go netopt.CliConnReadRight(userCliConn, bufferLen, channel)
	go netopt.CliCliConsumer(ctrlCliConn, userCliConn, channel)
}

func (obj coroutineInpcObj) run(ctrlAddr string, sourceAddr string, targetAddr string) {

	ctrlConnLst := []*net.TCPConn{}
	userConnLst := []*net.TCPConn{}

	go func() {
		// status checker
		ctrlCliConn := netopt.NetGetClientConn(ctrlAddr)
		ctrlConnLst = append(ctrlConnLst, ctrlCliConn)
		ctrlCliConn.Write([]byte("inpc_status:" + sourceAddr + "/" + targetAddr))
		go netopt.NetCliConnEcho(ctrlCliConn)

		// signal checker
		c := make(chan os.Signal)
		signal.Notify(c)
		fmt.Println("start signal listen ...")
		s := <-c
		fmt.Println("check signal is ", s)
		for _, conn := range ctrlConnLst {
			log.Println("inpc: ", fmt.Sprintf("%0x", &conn), " close.")
			conn.Close()
		}

		for _, conn := range userConnLst {
			log.Println("inpc: ", fmt.Sprintf("%0x", &conn), " close.")
			conn.Close()
		}
		os.Exit(1)
	}()

	for {
		ctrlCliConn := netopt.NetGetClientConn(ctrlAddr)
		// ctrlCliConn.Write([]byte(sourceAddr))
		ctrlCliConn.Write([]byte(sourceAddr + "/" + targetAddr))
		cmd := netopt.NetConnRecvDealOnce(ctrlCliConn, obj.bufferLen)
		if cmd == "ok" {
			log.Printf("ctrlCliConn recv ok ...")
		}

		// time.Sleep(time.Second * 10)
		// fmt.Printf("connect proxy server target ... ok")
		userCliConn := netopt.NetGetClientConn(targetAddr) // proxy server target

		ctrlConnLst = append(ctrlConnLst, ctrlCliConn)
		userConnLst = append(userConnLst, userCliConn)

		obj.acceptDeal(ctrlCliConn, userCliConn) // go
		// log.Println("accept after ...")
	}

}
