package netopt

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type ChanEle struct {
	flag int
	n    int
	bf   string
}

func SrvConnReadT(flag int, srvConn net.Conn, bufferLen int, channel chan<- ChanEle) {
	for {

		buffer := make([]byte, bufferLen)
		n, err := srvConn.Read(buffer)
		ce := ChanEle{flag, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			srvConn.Close()
			log.Println("srvConn:", fmt.Sprintf("%0x", &srvConn), " close.")
			ce := ChanEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func CliConnReadT(flag int, cliConn *net.TCPConn, bufferLen int, channel chan<- ChanEle) {
	for {
		buffer := make([]byte, bufferLen)
		n, err := cliConn.Read(buffer)
		ce := ChanEle{flag, n, string(buffer[:n])}
		channel <- ce

		if err != nil {
			cliConn.Close()
			log.Println("cliConn:", fmt.Sprintf("%0x", &cliConn), " close.")
			ce := ChanEle{-1, 0, ""}
			channel <- ce
			return
		}
	}
}

func SrvConnReadLeft(srvConn net.Conn, bufferLen int, channel chan<- ChanEle) {
	SrvConnReadT(0, srvConn, bufferLen, channel)
}

func SrvConnReadRight(srvConn net.Conn, bufferLen int, channel chan<- ChanEle) {
	SrvConnReadT(1, srvConn, bufferLen, channel)
}

func CliConnReadLeft(cliConn *net.TCPConn, bufferLen int, channel chan<- ChanEle) {
	CliConnReadT(0, cliConn, bufferLen, channel)
}

func CliConnReadRight(cliConn *net.TCPConn, bufferLen int, channel chan<- ChanEle) {
	CliConnReadT(1, cliConn, bufferLen, channel)
}

func SrvCliConsumer(srvConn net.Conn, cliConn *net.TCPConn, channel <-chan ChanEle) {
	for {
		ce := <-channel
		strLen := strconv.Itoa(strings.Count(ce.bf, ""))
		//log.Println("flag: ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen, " text: ", ce.bf)
		log.Println("flag: ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen)
		if ce.flag == 0 {
			cliConn.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == 1 {
			srvConn.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == -1 {
			cliConn.Close()
			srvConn.Close()
			log.Println("cliConn:", fmt.Sprintf("%0x", &cliConn), " close.")
			log.Println("srvConn:", fmt.Sprintf("%0x", &srvConn), " close.")
			return
		}
	}
}

func SrvSrvConsumer(srvConnLeft net.Conn, srvConnRight net.Conn, channel <-chan ChanEle) {
	for {
		ce := <-channel
		strLen := strconv.Itoa(strings.Count(ce.bf, ""))
		log.Println("flag: ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen)
		if ce.flag == 0 {
			srvConnRight.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == 1 {
			srvConnLeft.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == -1 {
			srvConnRight.Close()
			srvConnLeft.Close()
			log.Println("srvConnRight:", fmt.Sprintf("%0x", &srvConnRight), " close.")
			log.Println("srvConnLeft:", fmt.Sprintf("%0x", &srvConnLeft), " close.")
			return
		}
	}
}

func CliCliConsumer(cliConnLeft *net.TCPConn, cliConnRight *net.TCPConn, channel <-chan ChanEle) {
	for {
		ce := <-channel
		strLen := strconv.Itoa(strings.Count(ce.bf, ""))
		log.Println("flag: ", strconv.Itoa(ce.flag), " n: ", strconv.Itoa(ce.n), " buffers: ", strLen)
		if ce.flag == 0 {
			cliConnRight.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == 1 {
			cliConnLeft.Write([]byte(ce.bf))
			continue
		}
		if ce.flag == -1 {
			cliConnRight.Close()
			cliConnLeft.Close()
			log.Println("cliConnRight:", fmt.Sprintf("%0x", &cliConnRight), " close.")
			log.Println("cliConnLeft:", fmt.Sprintf("%0x", &cliConnLeft), " close.")
			return
		}
	}
}
