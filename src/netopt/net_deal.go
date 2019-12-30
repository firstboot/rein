package netopt

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func NetCheckClientPing(clientConn net.Conn) bool {
	clientConn.Write([]byte("hi"))

	buffer := make([]byte, 1024)
	n, err := clientConn.Read(buffer)
	if err == io.EOF {
		// clientConn.Close()
		return false
	}

	if err != nil {
		// clientConn.Close()
		return false
	}

	if string(buffer[:n]) == "hi" {
		return true
	}

	return false

}

func NetCheckClientConn(destServer string) bool {
	var status bool
	for i := 0; i < 2; i++ {
		status = NetCheckClientConnT(destServer)
		time.Sleep(time.Second * 1)
	}
	return status
}

func NetCheckClientConnT(destServer string) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", destServer)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		return false
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		return false
	}
	conn.Close()
	// log.Println("getClientConn ok")
	return true
}

func NetGetClientConn(destServer string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", destServer)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		os.Exit(1)
		// return nil
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error()))
		os.Exit(1)
		// return nil
	}
	// log.Println("getClientConn ok")
	return conn
}

func NetConnRecvDealOnce(conn net.Conn, bufferLen int) string {

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

func NetCliConnRecvPrint(clientConn *net.TCPConn, bufferLen int) {
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

func NetCliConnEcho(clientConn *net.TCPConn) {
	for {
		buffer := make([]byte, 1024)
		n, err := clientConn.Read(buffer)
		if err == io.EOF {
			clientConn.Close()
			// log.Println("clientConn:", fmt.Sprintf("%0x", &clientConn), " close.")
			return
		}
		clientConn.Write([]byte(string(buffer[:n])))
	}
}

func NetServerAccept(netListen net.Listener) net.Conn {

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

func NetServerListen(servAddr string) net.Listener {
	var netListen net.Listener
	var err error
	for i := 0; i < 3; i++ {
		netListen, err = net.Listen("tcp", servAddr)

		if err != nil {
			log.Println("net.Listen error!", err)
			// netListen.Close()
			netListen = nil
			// os.Exit(1)
		} else {
			return netListen
		}
	}
	return nil
}
