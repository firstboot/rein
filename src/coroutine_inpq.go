package main

import (
	"os"

	"./netopt"
)

type coroutineInpqObj struct {
	bufferLen int
}

func coroutineInpq() coroutineInpqObj {
	return coroutineInpqObj{10240}
}

func (obj coroutineInpqObj) run(ctrlAddr string) {
	ctrlCliConn := netopt.NetGetClientConn(ctrlAddr)
	ctrlCliConn.Write([]byte("inpq"))
	netopt.NetCliConnRecvPrint(ctrlCliConn, obj.bufferLen)
	ctrlCliConn.Close()
	os.Exit(0)
}
