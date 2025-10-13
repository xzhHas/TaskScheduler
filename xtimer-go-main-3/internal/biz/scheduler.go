package biz

import (
	"context"
	"time"

	"github.com/BitofferHub/pkg/middlewares/lock"
	"github.com/BitofferHub/pkg/middlewares/log"
	"github.com/BitofferHub/xtimer/internal/conf"
	"github.com/BitofferHub/xtimer/internal/utils"
)

// xtimerUseCase is a User usecase.
type SchedulerUseCase struct {
	confData  *conf.Data
	timerRepo TimerRepo
	taskRepo  TimerTaskRepo
	taskCache TaskCache
	tm        Transaction
	pool      WorkerPool
	trigger   *TriggerUseCase
}

// NewUserUseCase new a User usecase.
func NewSchedulerUseCase(confData *conf.Data, timerRepo TimerRepo, taskRepo TimerTaskRepo, taskCache TaskCache, trigger *TriggerUseCase) *SchedulerUseCase {
	return &SchedulerUseCase{
		confData:  confData,
		timerRepo: timerRepo,
		taskRepo:  taskRepo,
		taskCache: taskCache,
		pool:      NewGoWorkerPool(int(confData.Scheduler.WorkersNum)),
		trigger:   trigger,
	}
}

func (w *SchedulerUseCase) Work(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(w.confData.Scheduler.TryLockGapMilliSeconds) * time.Millisecond)
	defer ticker.Stop()
	// 每隔1s去执行这个定时任务
	for range ticker.C {
		select {
		case <-ctx.Done():
			log.WarnContextf(ctx, "stopped")
			return nil
		default:
		}

		w.handleSlices(ctx)
	}
	return nil
}

func (w *SchedulerUseCase) handleSlices(ctx context.Context) {
	for i := 0; i < w.getValidBucket(ctx); i++ {
		w.handleSlice(ctx, i)
	}
}

// 禁用动态分桶能力
func (w *SchedulerUseCase) getValidBucket(ctx context.Context) int {
	return int(w.confData.Scheduler.BucketsNum) // 一开始桶就只有 5 个
}

func (w *SchedulerUseCase) handleSlice(ctx context.Context, bucketID int) {
	defer func() {
		// 捕获异常，不再上报
		if r := recover(); r != nil {
			log.ErrorContextf(ctx, "handleSlice %v run err. Recovered from panic:%v", bucketID, r)
		}
	}()
	log.InfoContextf(ctx, "scheduler_%v start: %v", bucketID, time.Now())
	// 重复执行上一次的分片任务
	/*
		1.当前分钟执行的时候，可以重试，去执行前一分钟可能未执行的任务
		2.如果前一分钟执行失败，不会延迟分布式锁，那么这次重试就可以获取，并执行
		3.如果执行成功，就会延迟到130s，那么这次就无法获取到锁，去执行上一个分片任务了
	*/
	now := time.Now()
	if err := w.pool.Submit(func() {
		w.asyncHandleSlice(ctx, now.Add(-time.Minute), bucketID)
	}); err != nil {
		log.ErrorContextf(ctx, "[handle slice] submit task failed, err: %v", err)
	}

	// 执行当前时间分钟级的切片任务
	if err := w.pool.Submit(func() {
		w.asyncHandleSlice(ctx, now, bucketID)
	}); err != nil {
		log.ErrorContextf(ctx, "[handle slice] submit task failed, err: %v", err)
	}
	log.InfoContextf(ctx, "scheduler_%v end: %v", bucketID, time.Now())
}

func (w *SchedulerUseCase) asyncHandleSlice(ctx context.Context, t time.Time, bucketID int) {
	// 限制激活和去激活频次
	locker := lock.NewRedisLock(utils.GetTimeBucketLockKey(t, bucketID),
		lock.WithExpireSeconds(w.confData.Scheduler.TryLockSeconds))
	err := locker.Lock(ctx)
	if err != nil {
		//log.InfoContextf(ctx, "asyncHandleSlice 获取分布式锁失败: %v", err.Error())
		// 抢锁失败, 直接跳过执行, 下一轮
		return
	}

	log.InfoContextf(ctx, "get scheduler lock success, key: %s", utils.GetTimeBucketLockKey(t, bucketID))
	// 如果当前分钟级定时任务执行成功了，就会执行ack，延迟锁的过期时间
	ack := func() {
		if err := locker.DelayExpire(ctx, w.confData.Scheduler.SuccessExpireSeconds); err != nil {
			log.ErrorContextf(ctx, "expire lock failed, lock key: %s, err: %v", utils.GetTimeBucketLockKey(t, bucketID), err)
		}
	}
	// 触发器
	if err := w.trigger.Work(ctx, utils.GetSliceMsgKey(t, bucketID), ack); err != nil {
		log.ErrorContextf(ctx, "trigger work failed, SliceMsgKey[%v] err: %v", utils.GetSliceMsgKey(t, bucketID), err)
	}
}
