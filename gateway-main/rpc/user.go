package rpc

import (
	"context"
	"fmt"
	"github.com/BitofferHub/pkg/constant"
	"github.com/BitofferHub/pkg/middlewares/discovery"
	pb "github.com/BitofferHub/proto_center/api/user/v1"
	"github.com/go-kratos/kratos/v2/metadata"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"log"
)

const (
	UserSvrName = "user-svr"
)

func GetUserInfo(userID int64) string {
	// new etcd client
	dscry := discovery.GetServiceDiscovery(UserSvrName)
	endpoint := dscry.GetGrpcEndPoint()

	conn, err := transgrpc.DialInsecure(context.Background(),
		transgrpc.WithEndpoint(endpoint),
		transgrpc.WithMiddleware(mmd.Client()),
	)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	cli := pb.NewUserClient(conn)
	ctx := metadata.AppendToClientContext(context.Background(), constant.TraceID, "2233666")
	reply, err := cli.GetUser(ctx, &pb.GetUserRequest{UserID: userID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[grpc] GetUser reply %+v\n", reply)
	return reply.Data.UserName

}

func GetUserInfoByName(userName string) *pb.GetUserReplyData {
	dscry := discovery.GetServiceDiscovery(UserSvrName)
	endpoint := dscry.GetGrpcEndPoint()

	conn, err := transgrpc.DialInsecure(context.Background(),
		transgrpc.WithEndpoint(endpoint),
		//	transgrpc.WithDiscovery(dis),
		transgrpc.WithMiddleware(mmd.Client()),
	)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	cli := pb.NewUserClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, fmt.Sprintf("x-md-global-%s", constant.TraceID), "2233")
	fmt.Printf("rpc traceID %v\n", ctx)
	reply, err := cli.GetUserByName(ctx, &pb.GetUserByNameRequest{UserName: userName})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("[grpc] GetUser reply %+v\n", reply)
	return reply.Data
}

//func GetUserInfoByName(userName string) string {
//	// new etcd client
//	client, err := clientv3.New(clientv3.Config{
//		Endpoints: []string{"127.0.0.1:2379"},
//	})
//	if err != nil {
//		panic(err)
//	}
//	// new dis with etcd client
//	dis := etcd.New(client)
//
//	endpoint := "discovery:///user-svr"
//	conn, err := transgrpc.DialInsecure(context.Background(), transgrpc.WithEndpoint(endpoint), transgrpc.WithDiscovery(dis))
//	if err != nil {
//		panic(err)
//	}
//
//	defer conn.Close()
//	cli := pb.NewUserClient(conn)
//	reply, err := cli.GetUserByName(context.Background(), &pb.GetUserByNameRequest{UserName: userName})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("[grpc] GetUser reply %+v\n", reply)
//	return reply.Data.UserName
//
//}
//
