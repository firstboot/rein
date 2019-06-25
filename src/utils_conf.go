package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type utilsConfObj struct {
}

func utilsConf() utilsConfObj {
	obj := utilsConfObj{}
	return obj
}

func (obj utilsConfObj) getConf(confPath string) map[string]interface{} {
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	// log.Println(b)
	confStr := string(b)
	log.Println(confStr)

	var c map[string]interface{}
	jsonErr := json.Unmarshal([]byte(confStr), &c)
	if jsonErr != nil {
		log.Println(jsonErr)
		os.Exit(1)
	}
	return c
}

func (obj utilsConfObj) isExistKeyOfMap(key string, confMap map[string]interface{}) bool {
	for confKey := range confMap {
		if confKey == key {
			return true
		}
	}
	return false
}

func (obj utilsConfObj) isExistKeyOfStrMap(key string, confMap map[string]string) bool {
	for confKey := range confMap {
		if confKey == key {
			return true
		}
	}
	return false
}
