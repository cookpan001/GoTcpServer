package main

import (
	"net"
	//"bufio"
	"encoding/binary"
	"bytes"
	"log"
	"bufio"
	"./codec"
	"fmt"
	//"log"
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

type TestStruct struct {
	A byte
	B [2]byte
	C string
}

func bytes2struct() {
	tmp := []byte("haha")
	A := TestStruct{}
	buffer := bytes.NewReader(tmp)
	log.Println(buffer)
	binary.Read(buffer, binary.BigEndian, &A)
	log.Println(A)
}

func testRedisEncode() {
	//log.Println("int: ", codec.Serialize(1))
	//log.Println("str: ", codec.Serialize("abc"))
	//log.Println("bool: ", codec.Serialize(true))
	//log.Println("float: ", codec.Serialize(3.1415))
	//log.Println("array: ", codec.Serialize([]int{1,2,3,4}))
	log.Println("array: ", codec.Serialize([]string{"abc", "def"}))
	//fmt.Printf("%q\n", bytes.SplitN(bytes.TrimSpace([]byte("a\r\nb\r\n")), []byte("\r\n"), 2))
	tmp:= make([]byte, 4)
	binary.PutVarint(tmp, 10)
	fmt.Println(tmp)
	fmt.Println(binary.Varint(tmp))
}

func testRedisDecode() {
	str := codec.Serialize([]string{"abc", "def"})
	val, _, err := codec.Unserialize([]byte(str))
	log.Println(str)
	log.Println(val, err)
}

func testInterface() {
	t := []interface{} {
		uint8(1),
		int64(1343422),
		"abcd",
	}
	log.Println(t)
}

func main(){
	//conn, err := net.Dial("tcp", "127.0.0.1:3017")
	//if err != nil {
	//	log.Fatal("connect to server 127.0.0.1:2017 failed")
	//	return
	//}
	//line, err := bufio.NewReader(conn).ReadString('\n')
	//log.Println("message: "+line)
	//testRedisEncode()
	testRedisDecode()
}