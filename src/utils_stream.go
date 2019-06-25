package main

import "fmt"

type utilsStreamObj struct {
}

func utilsStream() utilsStreamObj {
	obj := utilsStreamObj{}
	return obj
}

// is exist 0x00 in stream
func (obj utilsStreamObj) isTail(buffer []byte, len int) bool {
	for inx := 0; inx < len; inx++ {
		// fmt.Printf("%0x ", buffer[inx])
		if buffer[inx] == 0 {
			return true
		}
	}
	return false
}

// is exist '\n' in stream
func (obj utilsStreamObj) isEnter(buffer []byte, len int) bool {
	for inx := 0; inx < len; inx++ {
		// fmt.Printf("%0x ", buffer[inx])
		if buffer[inx] == '\n' {
			return true
		}
	}
	return false
}

func (obj utilsStreamObj) printBytes(buffer []byte, len int) bool {
	for inx := 0; inx < len; inx++ {
		fmt.Printf("%0x ", buffer[inx])
	}
	fmt.Println("")
	return false
}
