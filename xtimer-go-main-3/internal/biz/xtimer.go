package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/BitofferHub/pkg/middlewares/lock"
	"github.com/BitofferHub/pkg/middlewares/log"
	"github.com/BitofferHub/xtimer/internal/conf"
	"github.com/BitofferHub/xtimer/internal/constant"
	"github.com/BitofferHub/xtimer/internal/utils"
	context2 "golang.org/x/net/context"
)

const defaultEnableGapSeconds = 3

// 表定义也需要放在biz层, 还是为了解耦biz与data层
type Timer struct {
	TimerId         int64 `gorm:"column:id"`
	App             string
	Name            string
	Status          int
	Cron            string
	NotifyHTTPParam string     `gorm:"column:notify_http_param;NOT NULL" json:"notify_http_param,omitempty"` // Http 回调参数
	CreateTime      *time.Time `gorm:"column:create_time;default:null"`
	ModifyTime      *time.Time `gorm:"column:modify_time;default:null"`
}

// TableName 表名
func (p *Timer) TableName() string {
	return "xtimer"
}

// 在一段未来的时间，将符合cron表达式的调度时间，批量导入到TimerTask对象数组中
func (t *Timer) BatchTasksFromTimer(executeTimes []time.Time) []*TimerTask {
	tasks := make([]*TimerTask, 0, len(executeTimes))
	for _, executeTime := range executeTimes {
		tasks = append(tasks, &TimerTask{
			App:      t.App,
			TimerID:  t.TimerId,
			Status:   constant.NotRunned.ToInt(),
			RunTimer: executeTime.UnixMilli(),
		})
	}
	return tasks
}

// TimerRepo is a Greater timerRepo.
type TimerRepo interface {
	Save(context.Context, *Timer) (*Timer, error)
	Update(context.Context, *Timer) (*Timer, error)
	FindByID(context.Context, int64) (*Timer, error)
	FindByAPP(context.Context, string) (*Timer, error) //app
	FindByStatus(context.Context, int) ([]*Timer, error)
	Delete(context.Context, string) error
	FindAllTimers(context.Context) ([]*Timer, error)
}

// xtimerUseCase is a User usecase.
type XtimerUseCase struct {
	confData  *conf.Data
	timerRepo TimerRepo
	taskRepo  TimerTaskRepo
	taskCache TaskCache
	tm        Transaction
	muc       *MigratorUseCase
}

// NewUserUseCase new a User usecase.
func NewXtimerUseCase(confData *conf.Data, timerRepo TimerRepo, taskRepo TimerTaskRepo, taskCache TaskCache, tm Transaction, muc *MigratorUseCase) *XtimerUseCase {
	return &XtimerUseCase{confData: confData, timerRepo: timerRepo, taskRepo: taskRepo, taskCache: taskCache, tm: tm, muc: muc}
}

func (uc *XtimerUseCase) CreateTimer(ctx context.Context, g *Timer) (*Timer, error) {
	return uc.timerRepo.Save(ctx, g)
}

// 激活 xtimer
func (uc *XtimerUseCase) EnableTimer(ctx context.Context, app string, timerId int64) error {
	// 限制激活和去激活频次(获取分布式锁)
	locker := lock.NewRedisLock(utils.GetEnableLockKey(app),
		lock.WithExpireSeconds(defaultEnableGapSeconds),
		lock.WithWatchDogMode())
	defer func(locker *lock.RedisLock, ctx context2.Context) {
		err := locker.Unlock(ctx)
		if err != nil {
			log.ErrorContextf(ctx, "EnableTimer 自动解锁失败", err.Error())
		}
	}(locker, ctx)

	err := locker.Lock(ctx)
	if err != nil {
		log.InfoContextf(ctx, "激活/去激活操作过于频繁，请稍后再试！", err.Error())
		// 抢锁失败, 直接跳过执行, 下一轮
		return nil
	}

	// 开启事务
	_ = uc.tm.InTx(ctx, func(ctx context.Context) error {
		// 1. 数据库获取Timer
		timer, err := uc.timerRepo.FindByID(ctx, timerId)
		if err != nil {
			log.ErrorContextf(ctx, "激活失败，timer不存在：timerId, err: %v", err)
			return err
		}
		//// TODO: 做一个小修改，就是使用app来查询这个xtimer，因为他给的protobuf中没有这个字段信息
		//timer, err := uc.timerRepo.FindByAPP(ctx, app)
		//if err != nil {
		//	log.ErrorContextf(ctx, "激活失败，timer不存在：timerId, err: %v", err)
		//	return err
		//}

		// 2. 校验状态
		// 如果是激活装填，就直接返回这个定时器的timerId，不需要去激活
		if timer.Status != constant.Unabled.ToInt() {
			// log.Error("测试这个是否正确")
			// time.Sleep(10 * time.Minute)
			return fmt.Errorf("Timer非Unable状态，激活失败，timerId:: %d", timerId)
		}

		// 修改timer为激活状态
		timer.Status = constant.Enabled.ToInt()
		_, err = uc.timerRepo.Update(ctx, timer)
		if err != nil {
			log.ErrorContextf(ctx, "激活失败，timer不存在：timerId, err: %v", err)
			return err
		}

		// 迁移数据
		// 如果是未激活状态，那么就需要将这个定时器激活，然后去迁移这个定时的数据
		// 在我看来这不是真正意义上的数据迁移，只是通过这个定时器的cron表达式，然后去批量的生成接下来一段时间内可能会触发任务的时间节点
		if err := uc.muc.MigratorTimer(ctx, timer); err != nil {
			log.ErrorContextf(ctx, "迁移timer失败: %v", err)
			return err
		}
		return nil
	})
	return nil
}

// 取消定时任务
func (uc *XtimerUseCase) UnableTimer(ctx context.Context, app string, timerId int64) error {
	_ = uc.tm.InTx(ctx, func(ctx context.Context) error {
		// 1. 数据库获取 Timer
		timer, err := uc.timerRepo.FindByAPP(ctx, app)
		if err != nil {
			log.ErrorContextf(ctx, "取消失败，timer不存在: timerId=%d, err=%v", timerId, err)
			return err
		}

		// 2. 校验状态
		if timer.Status != constant.Enabled.ToInt() {
			return fmt.Errorf("Timer非 Enabled 状态，无法取消，timerId=%d", timerId)
		}

		// 3. 修改 status 为 Unabled
		timer.Status = constant.Unabled.ToInt()
		_, err = uc.timerRepo.Update(ctx, timer)
		if err != nil {
			log.ErrorContextf(ctx, "取消失败，更新数据库错误: timerId=%d, err=%v", timerId, err)
			return err
		}
		return nil
	})
	return nil
}

// 删除定时任务
func (uc *XtimerUseCase) DelTimer(ctx context.Context, app string, timerId int64) error {
	return uc.timerRepo.Delete(ctx, app)
}

// 获取所有xtimer
func (uc *XtimerUseCase) GetTimers(ctx context.Context) ([]*Timer, error) {
	return uc.timerRepo.FindAllTimers(ctx)
}
