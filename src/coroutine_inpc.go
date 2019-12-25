package main

import (
	"log"
	"net"

	"./netopt"
)

type coroutineInpcObj struct {
	bufferLen int
}

func coroutineInpc() coroutineInpcObj {
	return coroutineInpcObj{10240}
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
	netopt.CliCliConsumer(ctrlCliConn, userCliConn, channel)
}

func (obj coroutineInpcObj) run(ctrlAddr string, sourceAddr string, targetAddr string) {

	for {
		for {
			ctrlCliConn := netopt.NetGetClientConn(ctrlAddr)
			// ctrlCliConn.Write([]byte(sourceAddr))
			ctrlCliConn.Write([]byte(sourceAddr + "/" + targetAddr))
			cmd := netopt.NetConnRecvDealOnce(ctrlCliConn, obj.bufferLen)
			if cmd == "ok" {
				log.Printf("ctrlCliConn recv ok ...")
			}
			if cmd == "reboot" {
				log.Printf("ctrlCliConn recv reboot ...")
				break
			}

			// time.Sleep(time.Second * 10)
			// fmt.Printf("connect proxy server target ... ok")
			userCliConn := netopt.NetGetClientConn(targetAddr) // proxy server target

			obj.acceptDeal(ctrlCliConn, userCliConn) // go
			log.Println("accept after ...")
		}
	}

	// for {
	// 	fmt.Println("ctrlCliConn:", fmt.Sprintf("%0x", &ctrlCliConn))
	// 	fmt.Println("userCliConn:", fmt.Sprintf("%0x", &userCliConn))
	// 	time.Sleep(time.Second * 1)
	// }
}
