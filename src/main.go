package main

/*
rein
date: 2019-6-17
author: lz
*/

import (
	"os"
	"time"
)

func main() {
	commandDeal().parse(os.Args)

	for {
		time.Sleep(time.Second * 1)
	}
}
