package tool

import (
	"encoding/binary"
	"bytes"
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