package main

import (
	"net"
	//"bufio"
	"encoding/binary"
	"bytes"
	"log"
	"bufio"
)

func test2(){
	conn, err := net.Dial("tcp", "127.0.0.1:3017")
	if err != nil {
		log.Fatal("connect to server 127.0.0.1:2017 failed")
		return
	}
	var body []byte = []byte("haha")
	var total int32 = int32(len(body))
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, total)
	conn.Write(buffer.Bytes())
	conn.Write(body)
	line, err := bufio.NewReader(conn).ReadString('\n')
	log.Println("total: ", buffer.Bytes())
	log.Println("body: ", body)
	log.Println("message: "+line)
}

func main(){
	//conn, err := net.Dial("tcp", "127.0.0.1:3017")
	//if err != nil {
	//	log.Fatal("connect to server 127.0.0.1:2017 failed")
	//	return
	//}
	//line, err := bufio.NewReader(conn).ReadString('\n')
	//log.Println("message: "+line)
	const str = "\x0D\x0A\x0D\x0A\x00\x0D"
	log.Println(len(str))
}