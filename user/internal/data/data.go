package data

import (
	"github.com/BitofferHub/pkg/middlewares/cache"
	"github.com/BitofferHub/pkg/middlewares/gormcli"
	"github.com/BitofferHub/user/internal/conf"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *cache.Client
}

var data *Data

func GetData() *Data {
	return data
}
func (p *Data) GetDB() *gorm.DB {
	return p.db
}

func (p *Data) GetCache() *cache.Client {
	return p.rdb
}

// NewData
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description:
//	@param dt
//	@return *Data
//	@return error
func NewData(dt *conf.Data) (*Data, error) {
	gormcli.Init(
		gormcli.WithAddr(dt.GetDatabase().GetAddr()),
		gormcli.WithUser(dt.GetDatabase().GetUser()),
		gormcli.WithPassword(dt.GetDatabase().GetPassword()),
		gormcli.WithDataBase(dt.GetDatabase().GetDataBase()),
		gormcli.WithMaxIdleConn(int(dt.GetDatabase().GetMaxIdleConn())),
		gormcli.WithMaxOpenConn(int(dt.GetDatabase().GetMaxOpenConn())),
		gormcli.WithMaxIdleTime(int64(dt.GetDatabase().GetMaxIdleTime())),
	)
	cache.Init(
		cache.WithAddr(dt.GetRedis().Addr),
		cache.WithPassWord(dt.GetRedis().PassWord),
		cache.WithDB(int(dt.GetRedis().Db)),
		cache.WithPoolSize(int(dt.GetRedis().PoolSize)))
	dta := &Data{db: gormcli.GetDB(), rdb: cache.GetRedisCli()}
	data = dta
	return dta, nil
}
