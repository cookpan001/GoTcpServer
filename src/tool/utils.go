package tool

import (
	"encoding/binary"
	"bytes"
	"net"
	"log"
)

const (
	PREFIX_BYTE = 4
	PROXY_VERSION1 = "PROXY "
	PROXY_VERSION2_START = "\x0D\x0A\x0D\x0A\x00\x0D"
	PROXY_VERSION2_END = "\x0D\x0A\x0D\x0A\x00\x0D"
)

func ToInt32(buf []byte) (int32){
	var total int32
	err := binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &total)
	if err != nil {
		return 0
	}
	return total
}

func FromInt32(i int32) ([]byte){
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, int32(i))
	return buffer.Bytes()
}

func ReadConn(conn net.Conn) ([]byte, error){
	tmp := make([]byte, PREFIX_BYTE)
	_, err := conn.Read(tmp)
	if err != nil {
		log.Printf("receive prefix error: %v\n", err)
		return []byte{}, err
	}
	var total int32
	err = binary.Read(bytes.NewBuffer(tmp), binary.BigEndian, &total)
	if err != nil {
		log.Printf("parse prefix error: %v\n", err)
		return []byte{}, err
	}
	log.Printf("receive prefix size: %v\n", total)
	buffer := make([]byte, total)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Printf("receive message error: %v\n", err)
		return []byte{}, err
	}
	log.Printf("receive from client: %v\n", string(buffer))
	return buffer, nil
}