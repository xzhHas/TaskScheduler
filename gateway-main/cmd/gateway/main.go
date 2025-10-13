package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/BitofferHub/gateway/limiter"
	"github.com/BitofferHub/pkg/middlewares/discovery"
	"github.com/redis/go-redis/v9"

	"github.com/BitofferHub/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs2 *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs2,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	initClient(bc.Data)
	var endpoints = []string{bc.Micro.GetLb().GetAddr()}
	//discovery.InitServiceDiscovery(endpoints, []string{"user-svr", "sec_kill-svr"})
	discovery.InitServiceDiscovery(endpoints, bc.Micro.GetLb().GetDisSvrList())
	err := limiter.InitLimiter("../../configs/router.json", rdb, 3, 10, 100)
	if err != nil {
		fmt.Println("panic : ", err)
		panic(err)
	}
	fmt.Println("come into wireapp")
	app, cleanup, err := wireApp(bc.Server, bc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

var (
	rdb *redis.Client
)

// initClient
// @Description: 初始化Redis，分布式锁，限流，消息队列，服务注册，使用Redis集合做一个业务方的权限鉴权
// @param cfData
// @return err
// @Router /api/initClient [${http_method}]
func initClient(cfData *conf.Data) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfData.GetRedis().GetAddr(),
		Password: cfData.GetRedis().GetPassWord(), // no password set
		DB:       0,                               // use default DB
		PoolSize: 100,                             // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}
