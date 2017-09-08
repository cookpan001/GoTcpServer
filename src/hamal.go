package main

import (
	"net"
	"log"
	"time"
	"./tool"
)

type HamalServer struct {
	count uint64
	proxyListener net.Listener
	appListener net.Listener
	clients map[uint64]uint32
	quit chan byte
}

func (server *HamalServer) Start() {
	select {
	case quit := <- server.quit:
		log.Println("quit: ", quit)
	}

}
// receive message from daemon
func (server *HamalServer) HandleApp() {
	defer server.appListener.Close()
	for {
		conn, err := server.appListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		server.count += 1
		go func(conn net.Conn) {
			for {
				buffer, err := tool.ReadConn(conn)
				if err != nil {
					break
				}
				log.Println("receive from app: " + len(buffer))
			}
			server.count -= 1
		}(conn)
	}
}
// receive message from proxy
func (server *HamalServer) HandleProxy() {
	defer server.proxyListener.Close()
	for {
		conn, err := server.proxyListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		server.count += 1
		go func(conn net.Conn) {
			for {
				buffer, err := tool.ReadConn(conn)
				if err != nil {
					break
				}
				log.Println("receive from proxy: " + len(buffer))
			}
			server.count -= 1
		}(conn)
	}
}

func (server *HamalServer) ProxyLisener(port string){
	serv, err := net.Listen("tcp", ":" + port)
	if err != nil {
		time.Sleep(time.Second)
		go server.ProxyLisener(port)
		return
	}
	server.proxyListener = serv
	go server.HandleProxy()
}

func (server *HamalServer) AppLisener(port string){
	serv, err := net.Listen("tcp", ":" + port)
	if err != nil {
		time.Sleep(time.Second)
		go server.AppLisener(port)
		return
	}
	server.appListener = serv
	go server.HandleApp()
}

func NewHamal() (*HamalServer){
	return &HamalServer{
		count: 0,
		clients: make(map[uint64]uint32),
	}
}

func main(){
	server := NewHamal()
	go server.proxyListener("2017")
	go server.appListener("3017")
	server.Start()
}
