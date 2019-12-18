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
type chanInpccEle struct {
	flag int
	n    int
	bf   string
}

type coroutineInpccObj struct {
	bufferLen int
}

func coroutineInpcc() coroutineInpccObj {
	return coroutineInpccObj{10240}
}

func (obj coroutineInpccObj) getClientConn(destServer string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", destServer)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		// os.Exit(1)
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		// os.Exit(1)
		return nil
	}
	log.Println("getClientConn ok")
	return conn
}

func (obj coroutineInpccObj) connRecvDealOnce(conn *net.TCPConn, bufferLen int) string {

	buffer := make([]byte, bufferLen)
	n, err := conn.Read(buffer)

	if err == io.EOF {
		conn.Close()
		log.Println("connRecvDealOnce conn:", fmt.Sprintf("%0x", &conn), " close.")
		return ""
	}

	bufferStr := string(buffer[:n])
	log.Println("conn:", fmt.Sprintf("%0x", &conn), len(bufferStr), " recv: ", bufferStr)
	return bufferStr
}

func (obj coroutineInpccObj) clientConnRead(clientConn *net.TCPConn, bufferLen int) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		if err != nil {
			clientConn.Close()
			log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			return
		}
		log.Println("clientConnRead: ", string(buffer[:n]))
	}
}

func (obj coroutineInpccObj) serverAccept(netListen net.Listener) net.Conn {

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

func (obj coroutineInpccObj) serverListen(servAddr string) net.Listener {

	netListen, err := net.Listen("tcp", servAddr)
	if err != nil {
		log.Println("rein inps net.Listen error!")
		// os.Exit(1)
		return nil
	}
	return netListen
}

func (obj coroutineInpccObj) acceptDealEx(ctrlCliConn *net.TCPConn, userCliConn *net.TCPConn, bufferLen int) {
	obj.communicationDeal(ctrlCliConn, bufferLen, userCliConn)

}

func (obj coroutineInpccObj) acceptDeal(ctrlCliConn *net.TCPConn, userCliConn *net.TCPConn) {
	obj.acceptDealEx(ctrlCliConn, userCliConn, obj.bufferLen)
}

func (obj coroutineInpccObj) orgiConnReadProducter(orgiConn *net.TCPConn, bufferLen int, channel chan<- chanInpccEle) {
	for {

		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		ce := chanInpccEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConnReadProducter orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanInpccEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpccObj) clientConnReadProducter(clientConn *net.TCPConn, bufferLen int, channel chan<- chanInpccEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanInpccEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConnReadProducter clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanInpccEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpccObj) consumerDeal(orgiConn *net.TCPConn, clientConn *net.TCPConn, channel <-chan chanInpccEle) {
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
			log.Println("consumerDeal clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			log.Println("consumerDeal orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			return
		}
	}
}

func (obj coroutineInpccObj) communicationDeal(ctrlCliConn *net.TCPConn, bufferLen int, userCliConn *net.TCPConn) {

	channel := make(chan chanInpccEle)
	go obj.orgiConnReadProducter(ctrlCliConn, bufferLen, channel)
	go obj.clientConnReadProducter(userCliConn, bufferLen, channel)
	obj.consumerDeal(ctrlCliConn, userCliConn, channel)
}

func (obj coroutineInpccObj) run(ctrlAddr string, sourceAddr string, targetAddr string) {
	log.Println("rein inpcc start...")

	// // inpcc control signal
	// inpccCtrlConn := obj.getClientConn(ctrlAddr)
	// log.Printf("inpccCtrlConn getClientConn ok ...")
	// inpccCtrlConn.Write([]byte("inpcc-ctrl"))

	for {

		ctrlCliConn := obj.getClientConn(ctrlAddr)
		// ctrlCliConn.Write([]byte(sourceAddr))
		ctrlCliConn.Write([]byte(sourceAddr + "/" + targetAddr))
		cmd := obj.connRecvDealOnce(ctrlCliConn, obj.bufferLen)
		if cmd == "ok" {
			log.Printf("ctrlCliConn recv ok ...")
		}

		// time.Sleep(time.Second * 10)
		fmt.Println("connect proxy server target ... ok")
		userCliConn := obj.getClientConn(targetAddr) // proxy server target

		obj.acceptDeal(ctrlCliConn, userCliConn) // block
		log.Println("accept after ...")

		// for {
		// 	msg := obj.connRecvDealOnce(inpccCtrlConn, obj.bufferLen)
		// 	if msg == "new" {
		// 		break
		// 	}

		// 	if msg == "" {
		// 		// os.Exit(1)
		// 		break
		// 	}
		// 	time.Sleep(time.Second * 1)
		// }
	}

	// for {
	// 	fmt.Println("ctrlCliConn:", fmt.Sprintf("%0x", &ctrlCliConn))
	// 	fmt.Println("userCliConn:", fmt.Sprintf("%0x", &userCliConn))
	// 	time.Sleep(time.Second * 1)
	// }
}
