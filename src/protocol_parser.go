package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/textproto"
	"strings"
)

type protocolParserObj struct {
}

func protocolParser() protocolParserObj {
	return protocolParserObj{}
}

func (obj protocolParserObj) getHTTPRequestHeader(requestContent string) *http.Request {
	reader := bufio.NewReader(strings.NewReader(requestContent))
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Fatal("fatal: ", err)
	}
	return req
}

func (obj protocolParserObj) setHTTPRequestHeader(req *http.Request) string {
	buf := new(bytes.Buffer)
	err := req.Write(buf)
	if err != nil {
		log.Println("setHTTPRequestHeader err")
	}
	return buf.String()
}

func (obj protocolParserObj) getHTTPResponseHeader(responseContent string) http.Header {
	reader := bufio.NewReader(strings.NewReader(responseContent))
	tp := textproto.NewReader(reader)
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		log.Fatal(err)
	}
	// http.Header and textproto.MIMEHeader are both just a map[string][]string
	httpHeader := http.Header(mimeHeader)
	// log.Println(httpHeader)
	return httpHeader
}
