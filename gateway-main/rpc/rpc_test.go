package rpc

import (
	"log"
	"net/url"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints)
	defer ser.Close()
	ser.WatchService("/microservices/user-svr")
	//ser.WatchService("/microservices")
	//	ser.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(2 * time.Second):
			//log.Println(ser.GetServices())
			addr := ser.GetHttpEndPoint()
			u, err := url.Parse(addr)
			if err != nil {
				panic(err)
			}
			log.Println(u.Host)
		}
	}
}
