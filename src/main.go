package main

/*
rein
date: 2019-6-17
author: lz
*/

import (
	"log"
	"os"
	"sync"
	"time"
)

var mainWg sync.WaitGroup
var mainRunFlag bool

func main() {
	for {
		mainRunFlag = true
		startAll()
		log.Println("reload rein...")
	}

}

func startAll() {
	mainWg.Add(1)
	commandDeal().parse(os.Args)
	log.Println("start rein...")
	mainWg.Wait()
}

func shutdownAll() {
	mainRunFlag = false
	time.Sleep(time.Second * 3)
	defer mainWg.Done()
	log.Println("shutdown all rein...")
}
