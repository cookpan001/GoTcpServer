package main

import (
	"net"
	"log"
	"time"
	"./tool"
)

type MQServer struct {
	acceptor net.Listener
	exchanger net.Listener
	waiter map[string]net.Listener
	quit chan byte
}

func Save(){

}

func Add(b []byte) {

}

func Notify(){

}

func (server *MQServer) acceptorLisener(port string){
	serv, err := net.Listen("tcp", ":" + port)
	if err != nil {
		time.Sleep(time.Second)
		go server.acceptorLisener(port)
		return
	}
	server.acceptor = serv
	go server.HandleAcceptor()
}

func (server *MQServer) exchangerLisener(port string){
	serv, err := net.Listen("tcp", ":" + port)
	if err != nil {
		time.Sleep(time.Second)
		go server.exchangerLisener(port)
		return
	}
	server.exchanger = serv
	go server.HandleExchanger()
}

// receive message from acceptor
func (server *MQServer) HandleAcceptor() {
	defer server.acceptor.Close()
	for {
		conn, err := server.acceptor.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn) {
			for {
				buffer, err := tool.ReadConn(conn)
				if err != nil {
					break
				}
				log.Println("receive from acceptor: " + len(buffer))
			}
		}(conn)
	}
}
// receive message from exchanger
func (server *MQServer) HandleExchanger() {
	defer server.exchanger.Close()
	for {
		conn, err := server.exchanger.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn) {
			for {
				buffer, err := tool.ReadConn(conn)
				if err != nil {
					break
				}
				log.Println("receive from exchanger: " + len(buffer))
			}
		}(conn)
	}
}

func (server *MQServer) ConnectWaiter(addr string) (net.Listener, error){
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect to waiter " + addr + " failed")
		return conn, err
	}
	server.waiter[addr] = conn
	return conn, nil
}

func (server *MQServer) HandleWaiter(address []string) {
	for _, addr := range address {
		_, ok := server.waiter[addr]
		if ok {
			continue
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println("connect to waiter " + addr + " failed")
			continue
		}
		server.waiter[addr] = conn
		go func(conn net.Listener) {
			for{
				_, ok := server.waiter[addr]
				if !ok {
					break
				}
				buffer, err := tool.ReadConn(conn)
				if err != nil {
					delete(server.waiter, addr)
					break
				}
				log.Println("from waiter: ", buffer)
			}
		}(conn)
	}
	time.Sleep(time.Second)
	go server.HandleWaiter(address)
}

func (server *MQServer) Start() {
	select {
	case quit := <- server.quit:
		log.Println("quit: ", quit)
	}

}

func NewMQServer() (*MQServer){
	return &MQServer{
		waiter: make(map[string]net.Listener),
	}
}

func init(){
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main(){
	server := NewMQServer()
	go server.acceptorLisener("2017")
	go server.exchangerLisener("3017")
	server.Start()
}