package service

import (
	"context"
)

// 定时迁移模块
func (s *XTimerService) MigratorTimer(ctx context.Context) error {
	return s.migratorUC.BatchMigratorTimer(ctx)
}
