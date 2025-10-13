package rpc

import (
	"context"
	"fmt"
	pb "github.com/BitofferHub/proto_center/api/user/v1"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func callHTTP() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	// new dis with etcd client
	dis := etcd.New(client)
	endpoint := "discovery:///user-svr"
	connHTTP, err := http.NewClient(
		context.Background(),
		http.WithEndpoint(endpoint),
		http.WithDiscovery(dis),
		http.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connHTTP.Close()
	httpClient := pb.NewUserHTTPClient(connHTTP)
	fmt.Printf("before call\n")
	reply, err := httpClient.GetUser(context.Background(), &pb.GetUserRequest{UserID: 1})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("[http] GetUser %+v\n", reply)
	time.Sleep(10 * time.Second)

}
