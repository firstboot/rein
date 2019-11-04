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
type chanInpcEle struct {
	flag int
	n    int
	bf   string
}

type coroutineInpcObj struct {
	bufferLen int
}

func coroutineInpc() coroutineInpcObj {
	return coroutineInpcObj{10240}
}

func (obj coroutineInpcObj) getClientConn(destServer string) *net.TCPConn {
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

func (obj coroutineInpcObj) connRecvDealOnce(conn *net.TCPConn, bufferLen int) string {

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

func (obj coroutineInpcObj) clientConnRead(clientConn *net.TCPConn, bufferLen int) {
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

func (obj coroutineInpcObj) serverAccept(netListen net.Listener) net.Conn {

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

func (obj coroutineInpcObj) serverListen(servAddr string) net.Listener {

	netListen, err := net.Listen("tcp", servAddr)
	if err != nil {
		log.Println("rein inps net.Listen error!")
		os.Exit(1)
	}
	return netListen
}

func (obj coroutineInpcObj) acceptDealEx(ctrlCliConn *net.TCPConn, userCliConn *net.TCPConn, bufferLen int) {
	obj.communicationDeal(ctrlCliConn, bufferLen, userCliConn)

}

func (obj coroutineInpcObj) acceptDeal(ctrlCliConn *net.TCPConn, userCliConn *net.TCPConn) {
	obj.acceptDealEx(ctrlCliConn, userCliConn, obj.bufferLen)
}

func (obj coroutineInpcObj) orgiConnReadProducter(orgiConn *net.TCPConn, bufferLen int, channel chan<- chanInpcEle) {
	for {

		buffer := make([]byte, bufferLen)
		n, err := orgiConn.Read(buffer)
		ce := chanInpcEle{0, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			orgiConn.Close()
			log.Println("orgiConnReadProducter orgiConn:", fmt.Sprintf("%0x", &orgiConn), " close.")
			ce := chanInpcEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpcObj) clientConnReadProducter(clientConn *net.TCPConn, bufferLen int, channel chan<- chanInpcEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		ce := chanInpcEle{1, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			clientConn.Close()
			log.Println("clientConnReadProducter clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			ce := chanInpcEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func (obj coroutineInpcObj) consumerDeal(orgiConn *net.TCPConn, clientConn *net.TCPConn, channel <-chan chanInpcEle) {
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

func (obj coroutineInpcObj) communicationDeal(ctrlCliConn *net.TCPConn, bufferLen int, userCliConn *net.TCPConn) {

	channel := make(chan chanInpcEle)
	go obj.orgiConnReadProducter(ctrlCliConn, bufferLen, channel)
	go obj.clientConnReadProducter(userCliConn, bufferLen, channel)
	obj.consumerDeal(ctrlCliConn, userCliConn, channel)
}

func (obj coroutineInpcObj) run(ctrlAddr string, sourceAddr string, targetAddr string) {

	for {
		ctrlCliConn := obj.getClientConn(ctrlAddr)
		ctrlCliConn.Write([]byte(sourceAddr))
		cmd := obj.connRecvDealOnce(ctrlCliConn, obj.bufferLen)
		if cmd == "ok" {
			log.Printf("ctrlCliConn recv ok ...")
		}

		// time.Sleep(time.Second * 10)
		// fmt.Printf("connect proxy server target ... ok")
		userCliConn := obj.getClientConn(targetAddr) // proxy server target

		obj.acceptDeal(ctrlCliConn, userCliConn) // go
		log.Println("accept after ...")
	}

	// for {
	// 	fmt.Println("ctrlCliConn:", fmt.Sprintf("%0x", &ctrlCliConn))
	// 	fmt.Println("userCliConn:", fmt.Sprintf("%0x", &userCliConn))
	// 	time.Sleep(time.Second * 1)
	// }
}
