package tool

type Codec interface{
	encode(data interface{}) ([]byte)
	decode() (interface{})
}

type MsgPackCodec struct {
	data []byte
}

func (message *MsgPackCodec) encode(arr interface{}) ([]byte){

	return []byte{}
}

func (message *MsgPackCodec) decode()  (interface{}){

	return nil
}

type JsonCodec struct {
	data []byte
}

func (message *JsonCodec) encode(arr interface{}) ([]byte){

	return []byte{}
}

func (message *JsonCodec) decode()  (interface{}){

	return nil
}

type ProtobufCodec struct {
	data []byte
}