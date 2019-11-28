package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type coroutineInpqObj struct {
	bufferLen int
}

func coroutineInpq() coroutineInpqObj {
	return coroutineInpqObj{10240}
}

func (obj coroutineInpqObj) getClientConn(destServer string) *net.TCPConn {
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
	// log.Println("getClientConn ok")
	return conn
}

func (obj coroutineInpqObj) connRecvDealOnce(conn *net.TCPConn, bufferLen int) string {

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

func (obj coroutineInpqObj) clientConnRead(clientConn *net.TCPConn, bufferLen int) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := clientConn.Read(buffer)
		if err == io.EOF {
			clientConn.Close()
			// log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			return
		}
		fmt.Print(string(buffer[:n]))
	}
}

func (obj coroutineInpqObj) run(ctrlAddr string) {
	ctrlCliConn := obj.getClientConn(ctrlAddr)
	ctrlCliConn.Write([]byte("inpq"))
	obj.clientConnRead(ctrlCliConn, obj.bufferLen)
	ctrlCliConn.Close()
	os.Exit(0)
}
