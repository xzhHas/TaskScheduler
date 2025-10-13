//go:build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/BitofferHub/xtimer/internal/biz"
	"github.com/BitofferHub/xtimer/internal/conf"
	"github.com/BitofferHub/xtimer/internal/data"
	"github.com/BitofferHub/xtimer/internal/interfaces"
	"github.com/BitofferHub/xtimer/internal/server"
	"github.com/BitofferHub/xtimer/internal/service"
	"github.com/BitofferHub/xtimer/internal/task"
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,       // 最底层，负责数据库、缓存、事务等
		biz.ProviderSet,        // 业务层，核心业务逻辑
		service.ProviderSet,    // 服务层，对外提供api接口，调用biz层
		interfaces.ProviderSet, // 接口层，负责api请求解析，调用service层
		task.ProviderSet,
		newApp))
}
