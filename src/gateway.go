package main

import (
	"log"
	"net"
	"time"
	"./tool"
)

type GatewayServer struct {
	listener net.Listener
}

func (server *GatewayServer) Start() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.Handle(connection)
	}
}

func (server *GatewayServer) Handle(conn net.Conn) {
	conn.SetDeadline(time.Minute)
	for{
		buffer, err := tool.ReadConn(conn)
		if err != nil {
			break
		}
		log.Println("receive from client: ", len(buffer))
		switch buffer[0] {
		case 1:
			HandleRedis(buffer[1:2], buffer[2:])
		case 2:
			HandleMySQL(buffer[1:2], buffer[2:])
		}
	}
}

func HandleRedis(name byte, request []byte){

}

func HandleMySQL(name byte, request []byte){

}

func NewGateway(port string) (*HamalServer){
	serv, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &GatewayServer{
		listener: serv,
	}
}

func init(){

}

func main(){
	server := NewGateway("2017")
	server.Start()
}