package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// channel element struct
type chanInpcsEle struct {
	flag int
	n    int
	bf   string
}

type coroutineInpcsObj struct {
	bufferLen int
}

func coroutineInpcs() coroutineInpcsObj {
	return coroutineInpcsObj{10240}
}

func (obj coroutineInpcsObj) serverAccept(netListen net.Listener) net.Conn {

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

func (obj coroutineInpcsObj) serverListen(servAddr string) net.Listener {
	var netListen net.Listener
	var err error
	for i := 0; i < 3; i++ {
		netListen, err = net.Listen("tcp", servAddr)
		if err != nil {
			log.Println("rein Inpcs net.Listen error!", err)
			netListen = nil
			// os.Exit(1)
		} else {
			return netListen
		}
	}
	return nil
}

func (obj coroutineInpcsObj) connRecvDealOnce(conn net.Conn, bufferLen int) string {

	buffer := make([]byte, bufferLen)
	n, err := conn.Read(buffer)

	if err == io.EOF {
		conn.Close()
		log.Println("connRecvDealOnce conn:", fmt.Sprintf("%0x", &conn), " close.")
		return ""
	}

	bufferStr := string(buffer[:n])
	log.Println("connRecvDealOnce conn:", fmt.Sprintf("%0x", &conn), len(bufferStr), " recv: ", bufferStr)
	return bufferStr
}

func (obj coroutineInpcsObj) acceptDealEx(userServLis net.Listener, ctrlServConn net.Conn, bufferLen int) {
	var firstRunFlag = false

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
}

func (obj coroutineInpcsObj) acceptDeal(userServLis net.Listener, ctrlServConn net.Conn) {
	obj.acceptDealEx(userServLis, ctrlServConn, obj.bufferLen)
}

func (obj coroutineInpcsObj) getClientConn(destServer string) *net.TCPConn {
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

func (obj coroutineInpcsObj) orgiConnReadProducter(orgiConn net.Conn, bufferLen int, channel chan<- chanInpcsEle) {
	for {

		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		ce := chanInpcsEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanInpcsEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpcsObj) clientConnReadProducter(clientConn net.Conn, bufferLen int, channel chan<- chanInpcsEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanInpcsEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanInpcsEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpcsObj) consumerDeal(orgiConn net.Conn, clientConn net.Conn, channel <-chan chanInpcsEle) {
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

func (obj coroutineInpcsObj) communicationDeal(userServConn net.Conn, bufferLen int, ctrlServConn net.Conn) {

	channel := make(chan chanInpcsEle)
	go obj.orgiConnReadProducter(userServConn, bufferLen, channel)
	go obj.clientConnReadProducter(ctrlServConn, bufferLen, channel)
	obj.consumerDeal(userServConn, ctrlServConn, channel)
}

func (obj coroutineInpcsObj) execute(ctrlAddr string) {

	for {
		// var userServLis net.Listener
		// runFlag := true

		fmt.Println("new cycle ...")
		ctrlClientConn := obj.getClientConn(ctrlAddr)
		ctrlClientConn.Write([]byte("inpcs"))
		fmt.Println("connRecvDealOnce wait ...")
		sourceAddr := obj.connRecvDealOnce(ctrlClientConn, obj.bufferLen)
		fmt.Println("connRecvDealOnce ok ...")
		// error pair
		if sourceAddr == "exit" {
			continue
		}

		userServLis := obj.serverListen(sourceAddr) // proxy server source
		if userServLis == nil {
			break
		}
		fmt.Println("1")

		// go func() {
		// 	obj.acceptDeal(userServLis, ctrlClientConn) // block
		// 	userServLis.Close()
		// 	ctrlClientConn.Close()
		// 	runFlag = false
		// 	fmt.Println("go func botton")
		// }()

		// fmt.Println("botton")

		///////////////////////////

		obj.acceptDeal(userServLis, ctrlClientConn) // block
		userServLis.Close()
		ctrlClientConn.Close()

		// runFlag = false

		////////////////////////

		// for {
		// 	// fmt.Println("runFlag :", runFlag)
		// 	time.Sleep(time.Second * 1)
		// 	if runFlag == false {
		// 		break
		// 	}
		// }

		// for {
		// 	msg := obj.connRecvDealOnce(inpcsCtrlConn, obj.bufferLen)
		// 	if msg == "new" {
		// 		// userServLis.Close()
		// 		break
		// 	}

		// 	if msg == "" {
		// 		// os.Exit(1)
		// 		// userServLis.Close()
		// 		break
		// 	}
		// 	time.Sleep(time.Second * 1)
		// }
	}
}

func (obj coroutineInpcsObj) run(ctrlAddr string) {
	log.Println("rein inpcs start...")

	for {
		// inpcs control signal
		inpcsCtrlConn := obj.getClientConn(ctrlAddr)
		inpcsCtrlConn.Write([]byte("inpcs-ctrl"))

		for {
			msg := obj.connRecvDealOnce(inpcsCtrlConn, obj.bufferLen)
			if len(msg) >= 10 && msg[:11] == "leftConnNum" {
				leftConnNum := msg[12:]
				num, _ := strconv.Atoi(leftConnNum)
				log.Println("rein inpcs leftConnNum = ", num)
				for i := 0; i < num; i++ {
					go obj.execute(ctrlAddr)
				}
			}
			if msg == "" {
				inpcsCtrlConn.Close()
				time.Sleep(time.Second * 3)
				break
			}

			if msg == "reboot" {
				inpcsCtrlConn.Close()
				time.Sleep(time.Second * 3)
				break
			}
		}
	}
}
