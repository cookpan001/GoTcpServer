package main

import (
	"net"
	//"bufio"
	"encoding/binary"
	"bytes"
	"log"
	"bufio"
)

func main(){
	conn, err := net.Dial("tcp", "127.0.0.1:3017")
	if err != nil {
		log.Fatal("connect to server 127.0.0.1:3017 failed")
		return
	}
	var body []byte = []byte("haha")
	var total int32 = int32(len(body))
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, total)
	conn.Write(buffer.Bytes())
	conn.Write(body)
	log.Println("total: ", buffer.Bytes())
	log.Println("body: ", string(body))
	line, err := bufio.NewReader(conn).ReadString('\n')
	log.Println("message: "+line)
}