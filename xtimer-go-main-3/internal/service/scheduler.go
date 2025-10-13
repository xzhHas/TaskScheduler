package service

import (
	"context"
)

// 任务tast调度执行流程
func (s *XTimerService) ScheduleTask(ctx context.Context) error {
	return s.schedulerUC.Work(ctx)
}
