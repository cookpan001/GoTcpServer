package main

import (
	"net"
	"log"
	"time"
	"fmt"
	"encoding/binary"
	"bytes"
	"bufio"
	"strings"
	"./tool"
)

type ProxyProtocol2 struct {

}

type Client struct{
	conn net.Conn
	server Server
	uid uint64
	id uint64
	ip string
}

type Server struct{
	network, address string
	count uint64
	listener net.Listener
	clients map[uint64]Client
	toChan chan []byte
	fromChan chan []byte
}

type Subscriber struct {
	network, address string
	conn net.Conn
}

func ReadExtraBytes(conn net.Conn, head []byte, body int32) ([]byte, error) {
	if len(head) > 0 {
		tmp := make([]byte, tool.PREFIX_BYTE - len(head))
		_, err := conn.Read(tmp)
		if err != nil {
			log.Printf("fail to receive extra prefix: %v\n", err)
			return []byte{}, err
		}
		prefix := make([]byte, tool.PREFIX_BYTE)
		for k, v := range head {
			prefix[k] = v
		}
		for k, v := range tmp {
			prefix[k] = v
		}
		total := tool.ToInt32(prefix)
		bodyBuffer := make([]byte, total)
		_, err = conn.Read(bodyBuffer)
		if err != nil {
			log.Printf("fail to receive body: %v\n", err)
			return []byte{}, err
		}
		return bodyBuffer, nil
	}
	if body > 0 {
		tmp := make([]byte, body)
		_, err := conn.Read(tmp)
		if err != nil {
			log.Printf("receive prefix error: %v\n", err)
			return []byte{}, err
		}
		return tmp, nil
	}
	return []byte{}, nil
}

func HandleRemainingBytes(buffer []byte, conn net.Conn, c chan []byte){
	var cur, total int32
	count := int32(len(buffer))
	for cur=0; cur < count; {
		if len(buffer[cur:]) < tool.PREFIX_BYTE {
			extra, err := ReadExtraBytes(conn, buffer[cur:], 0)
			if err != nil {
				return
			}
			c <- extra
			return
		}
		err := binary.Read(bytes.NewBuffer(buffer[cur:cur+tool.PREFIX_BYTE]), binary.BigEndian, &total)
		if err != nil {
			return
		}
		cur += tool.PREFIX_BYTE
		if count < cur {
			extra, err := ReadExtraBytes(conn, []byte{}, cur - count)
			if err != nil {
				return
			}
			c <- bytes.Join([][]byte{
				buffer[cur:],
				extra,
			}, []byte{})
			return
		}
		c <- buffer[cur:cur+total]
		cur += total
		if cur == count {
			break
		}
	}
}

func (client *Client) Receive() ([]byte, error){
	return tool.ReadConn(client.conn)
}

func (client *Client) Send(buf interface{}) (bool){
	fmt.Printf("send: %v, %v", buf.(string), len(buf.(string)))
	switch v := buf.(type) {
	case string:
		_, err := client.conn.Write([]byte(buf.(string)))
		if err != nil {
			fmt.Printf("send error: %v", err)
		}
	case []uint8:
		_, err := client.conn.Write(buf.([]byte))
		if err != nil {
			fmt.Printf("send error: %v", err)
		}
	default:
		fmt.Printf("unsupport type: %v", v)
		return false
	}
	return true
}

func (client *Client) OnConnect(server *Server) {
	head := make([]byte, 6)
	_, err := client.conn.Read(head)
	if err != nil {
		return
	}
	if string(head) == tool.PROXY_VERSION1 {
		buffer :=  bufio.NewReader(client.conn)
		str, e := buffer.ReadString('\n')
		if e != nil {
			log.Println("fail to decode PROXY protocol version 1, err: ", err)
			return
		}
		arr := strings.Fields(str)
		client.ip = arr[1]
		log.Println(client.ip)
		log.Println(buffer.Buffered())
		left := buffer.Buffered()
		if left > 0 {
			leftBytes := make([]byte, left)
			_, err = buffer.Read(leftBytes)
			if err != nil {
				return
			}
			HandleRemainingBytes(leftBytes, client.conn, server.toChan)
		}
		return
	}else if string(head) == tool.PROXY_VERSION2_START{
		headExtra := make([]byte, 6)
		_, err := client.conn.Read(headExtra)
		if err == nil {
			if string(headExtra) == tool.PROXY_VERSION2_END{

			}
			return
		}
	}else{
		client.ip = client.conn.RemoteAddr().String()
	}
	log.Println("no PROXY protocol")
	var total int32
	err = binary.Read(bytes.NewBuffer(head[:5]), binary.BigEndian, &total)
	if err == nil {
		buffer := make([]byte, total - 2)
		_, err = client.conn.Read(buffer)
		if err != nil {
			fmt.Printf("receive message error: %v\n", err)
			return
		}
		s := [][]byte{
			head[5:],
			buffer,
		}
		server.toChan <- bytes.Join(s, []byte(""))
	}
}

func (server *Server) Handle(client Client){
	defer client.conn.Close()
	client.conn.SetDeadline(time.Now().Add(time.Minute))
	client.OnConnect(server)
	for{
		buffer, err := client.Receive()
		if err != nil {
			break
		}
		server.toChan <- buffer
	}
	_, ok := server.clients[client.id]
	if ok {
		delete(server.clients, client.id)
	}
}

func SaveMessage(c chan []byte){
	for {
		select {
		case buf := <-c:
			if len(buf) > 0 {
				log.Println("save message: ", buf)
			}
		}
	}
}

func PushReply(c chan []byte){
	for {
		select {
		case buf := <-c:
			if len(buf) > 0 {
				log.Println("push message: ", buf)
			}
		}
	}
}

func (server *Server) Start(){
	defer server.listener.Close()
	defer close(server.toChan)
	defer close(server.fromChan)
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		c := Client{
			conn: connection,
			server: *server,
			uid: 0,
			id: server.count,
			ip: "",
		}
		server.count += 1
		server.clients[c.id] = c
		go server.Handle(c)
	}
}

func NewTcpServer(host, port string) (*Server, error){
	address := net.JoinHostPort(host, port)
	serv, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Server{
		network: "tcp",
		address: address,
		listener: serv,
		count: 0,
		clients: make(map[uint64]Client),
		toChan: make(chan []byte, 10),
		fromChan:  make(chan []byte, 10),
	}, nil
}

func NewSubscriber(host, port, regMsg string, c chan []byte){
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal("connect to server 127.0.0.1:2017 failed")
		return
	}
	var body []byte = []byte(regMsg)
	var total int32 = int32(len(body))
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, total)
	conn.Write(buffer.Bytes())
	conn.Write(body)
	log.Println("total: ", buffer.Bytes())
	log.Println("body: ", body)
	go func(conn net.Conn){
		for{
			recvBytes, err := tool.ReadConn(conn)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("subscribe recieved: ", recvBytes)
			c <- recvBytes
		}
	}(conn)
	go  NewSubscriber(host, port, regMsg, c)
}

func init(){
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main(){
	server, err := NewTcpServer("127.0.0.1", "2017")
	log.Println("listen to " + server.address)
	if err == nil {
		go SaveMessage(server.toChan)
		go PushReply(server.fromChan)
		server.Start()
		log.Println("start to accept connection")
	}
	log.Println("server ended")
}