package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

// channel element struct
type chanInpssEle struct {
	flag int
	n    int
	bf   string
}

type coroutineInpssObj struct {
	bufferLen int
}

func coroutineInpss() coroutineInpssObj {
	return coroutineInpssObj{10240}
}

func (obj coroutineInpssObj) serverAccept(netListen net.Listener) net.Conn {

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

func (obj coroutineInpssObj) serverListen(servAddr string) net.Listener {
	var netListen net.Listener
	var err error
	for i := 0; i < 3; i++ {
		netListen, err = net.Listen("tcp", servAddr)
		if err != nil {
			log.Println("rein Inpss net.Listen error!", err)
			netListen.Close()
			netListen = nil
			// os.Exit(1)
		} else {
			return netListen
		}
	}
	return nil
}

func (obj coroutineInpssObj) connRecvDealOnce(conn net.Conn, bufferLen int) string {

	buffer := make([]byte, bufferLen)
	n, err := conn.Read(buffer)

	if err == io.EOF {
		conn.Close()
		log.Println("conn:", fmt.Sprintf("%0x", &conn), " close.")
		return ""
	}

	bufferStr := string(buffer[:n])
	log.Println("** conn:", fmt.Sprintf("%0x", &conn), len(bufferStr), " recv: ", bufferStr)
	return bufferStr
}

func (obj coroutineInpssObj) acceptDealEx(userServLis net.Conn, ctrlServConn net.Conn, bufferLen int) {
	obj.communicationDeal(userServLis, bufferLen, ctrlServConn)

}

func (obj coroutineInpssObj) acceptDeal(userServLis net.Conn, ctrlServConn net.Conn) {
	go obj.acceptDealEx(userServLis, ctrlServConn, obj.bufferLen)
}

func (obj coroutineInpssObj) getClientConn(destServer string) *net.TCPConn {
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

func (obj coroutineInpssObj) orgiConnReadProducter(orgiConn net.Conn, bufferLen int, channel chan<- chanInpssEle) {
	// var firstRunFlag = false
	for {

		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		// if firstRunFlag == false {
		// 	firstCe := chanInpssEle{0, n, string("ok")}
		// 	channel <- firstCe
		// 	firstRunFlag = true
		// }
		ce := chanInpssEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanInpssEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpssObj) clientConnReadProducter(clientConn net.Conn, bufferLen int, channel chan<- chanInpssEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanInpssEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanInpssEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpssObj) consumerDeal(orgiConn net.Conn, clientConn net.Conn, channel <-chan chanInpssEle) {
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

func (obj coroutineInpssObj) communicationDeal(userServConn net.Conn, bufferLen int, ctrlServConn net.Conn) {

	channel := make(chan chanInpssEle)
	go obj.orgiConnReadProducter(userServConn, bufferLen, channel)
	go obj.clientConnReadProducter(ctrlServConn, bufferLen, channel)
	obj.consumerDeal(userServConn, ctrlServConn, channel)
}

func (obj coroutineInpssObj) run(ctrlAddr string) {
	ctrlServLis := obj.serverListen(ctrlAddr)
	leftConns := []net.Conn{}
	leftConnTexts := []string{}
	rightConns := []net.Conn{}
	inpccConnPairs := make(map[string]int)

	var inpccCtrlConn net.Conn
	// var inpcsCtrlConn net.Conn

	for {
		log.Println("wait for ctrlServConn link in ...")
		ctrlServConn := obj.serverAccept(ctrlServLis) // block
		log.Println("ctrlServConn link in ok ...")
		msg := obj.connRecvDealOnce(ctrlServConn, obj.bufferLen)
		log.Println("connRecvDealOnce link in ok ...")

		if msg == "inpq" {
			for key, ele := range inpccConnPairs {
				log.Println("inpq: ", key, ", ", ele)
				ctrlServConn.Write([]byte(key + "\n"))
			}
			ctrlServConn.Close()
			continue
		}

		if msg == "inpcc-ctrl" {
			inpccCtrlConn = ctrlServConn
			log.Println("inpccCtrlConn: ", fmt.Sprintf("%0x", &inpccCtrlConn))
			continue
		}

		// if msg == "inpcs-ctrl" {
		// 	inpcsCtrlConn = ctrlServConn
		// 	log.Println("inpcsCtrlConn: ", fmt.Sprintf("%0x", &inpcsCtrlConn))
		// 	continue
		// }

		if len(msg) >= 6 && msg[:7] == "0.0.0.0" {
			leftConns = append(leftConns, ctrlServConn)
			log.Println("leftConn: ", fmt.Sprintf("%0x", &ctrlServConn))
			inpccConnPairs[msg] = 0
			pos := strings.Index(msg, "/")
			sourceAddr := msg[:pos]
			leftConnTexts = append(leftConnTexts, sourceAddr)
		}

		if msg == "inpcs" {
			rightConns = append(rightConns, ctrlServConn)
			log.Println("rightConn: ", fmt.Sprintf("%0x", &ctrlServConn))
			//rightConn.Write([]byte(sourceAddr))
		}

		fmt.Println("before len(leftConns): ", len(leftConns), ", len(rightConns): ", len(rightConns))

		// if len(leftConns) > 0 && len(rightConns) > 0 && len(leftConns) < len(rightConns) {
		// 	rightConns[len(rightConns)-1].Write([]byte("exit"))
		// 	rightConns[len(rightConns)-1].Close()
		// 	rightConns = rightConns[:len(rightConns)-1]
		// }

		if len(leftConns) > 0 && len(rightConns) > 0 {

			leftLen := len(leftConns)
			rightLen := len(rightConns)
			arrSize := int(math.Min(float64(leftLen), float64(rightLen)))

			for i := 0; i < arrSize; i++ {
				text := leftConnTexts[i]
				rightConns[i].Write([]byte(text))
				obj.acceptDeal(rightConns[i], leftConns[i]) // block
			}

			for i := 0; i < arrSize; i++ {
				leftConns = leftConns[i+1:]
				leftConnTexts = leftConnTexts[i+1:]
				rightConns = rightConns[i+1:]
			}

		}

		if len(leftConns) == 0 && len(rightConns) > 0 {
			inpccCtrlConn.Write([]byte("new"))
		}

		fmt.Println("after len(leftConns): ", len(leftConns), ", len(rightConns): ", len(rightConns))
	}

	// for {

	// 	var leftConn net.Conn
	// 	var rightConn net.Conn
	// 	var sourceAddr string
	// 	// leftConn
	// 	log.Println("leftConn wait for ctrlServConn link in ...")
	// 	ctrlServConn := obj.serverAccept(ctrlServLis) // block
	// 	log.Println("leftConn ctrlServConn link in ok ...")
	// 	msg := obj.connRecvDealOnce(ctrlServConn, obj.bufferLen)
	// 	log.Println("leftConn connRecvDealOnce link in ok ...")
	// 	if len(msg) >= 6 && msg[:7] == "0.0.0.0" {
	// 		leftConn = ctrlServConn
	// 		log.Println("leftConn: ", fmt.Sprintf("%0x", &leftConn))
	// 		sourceAddr = msg

	// 		// rightConn
	// 		log.Println("rightConn wait for ctrlServConn link in ...")
	// 		ctrlServConn = obj.serverAccept(ctrlServLis) // block
	// 		log.Println("rightConn ctrlServConn link in ok ...")
	// 		msg = obj.connRecvDealOnce(ctrlServConn, obj.bufferLen)
	// 		log.Println("rightConn connRecvDealOnce link in ok ...")
	// 		if msg == "inpcs" {
	// 			rightConn = ctrlServConn
	// 			log.Println("rightConn: ", fmt.Sprintf("%0x", &rightConn))
	// 			rightConn.Write([]byte(sourceAddr))
	// 		}

	// 		if leftConn != nil && rightConn != nil {
	// 			// go func() {
	// 			//              cs         cc
	// 			obj.acceptDeal(rightConn, leftConn) // block
	// 			// rightConn.Close()
	// 			// leftConn.Close()

	// 			// }()
	// 		}

	// 	} else {
	// 		ctrlServConn.Write([]byte("exit"))
	// 		ctrlServConn.Close()
	// 	}

	// 	fmt.Println("bottom")

	// }

}
