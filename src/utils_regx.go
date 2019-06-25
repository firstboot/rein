package main

import (
	"regexp"
	"strings"
)

type utilsRegxObj struct {
}

func utilsRegx() utilsRegxObj {
	obj := utilsRegxObj{}
	return obj
}

// replace by regx
func (p utilsRegxObj) replace(data string, regx string, newStr string) string {
	reg := regexp.MustCompile(regx)
	str := reg.ReplaceAllString(data, newStr)
	return str
}

// is exist by regx in data
func (p utilsRegxObj) isExist(data string, regx string) bool {
	reg := regexp.MustCompile(regx)
	return reg.MatchString(data)
}

// find str by regx in data
func (p utilsRegxObj) getFindStrs(data string, regx string) []string {
	reg := regexp.MustCompile(regx)
	return reg.FindAllString(data, -1)
}

func (p utilsRegxObj) getRootLocation(data string) string {
	var beg = strings.Index(data, " ")
	var end = strings.LastIndex(data, " ")
	// fmt.Println("##", beg, end)
	subData := data[beg:end]
	beg = strings.Index(subData, "/")
	end = strings.LastIndex(subData, "/")
	if beg != end {
		end = strings.Index(subData[beg+1:], "/")
		// fmt.Println("beg!=end", beg, end+beg+1, subData)
		return subData[beg : end+beg+1]
	}
	// fmt.Println("beg==end", beg, end)
	return strings.ReplaceAll(subData[beg:], " ", "")

}
