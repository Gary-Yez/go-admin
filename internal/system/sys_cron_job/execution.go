package sys_cron_job

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/scheduler"
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"time"
)

func nextRun(expr string, now time.Time) (time.Time, error) {
	schedule, err := cron.ParseStandard(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("Cron 表达式无效：%w", err)
	}
	next := schedule.Next(now.In(scheduler.DefaultLocation))
	if next.IsZero() {
		return time.Time{}, errors.New("Cron 表达式没有可用的下次执行时间")
	}
	return next, nil
}

type execution struct{ id uint }

// 数据库行锁串行化领取、编辑和删除，唯一索引防止同一计划批次重复执行。
func (j *SysCronJob) Claim(ctx context.Context) (scheduler.Run, error) {
	var claimed *execution
	err := state.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current SysCronJob
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, j.Id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		now := time.Now().UTC()
		if !current.Enable || current.Version != j.Version || current.NextRunTime == nil || current.NextRunTime.After(now) {
			return nil
		}
		next, err := nextRun(current.Cron, now)
		if err != nil {
			return err
		}
		host, err := os.Hostname()
		if err != nil {
			return err
		}
		record := &SysCronJobLog{JobId: current.Id, ScheduledAt: *current.NextRunTime, StartTime: now,
			Status: "running", Instance: fmt.Sprintf("%s:%d", host, os.Getpid())}
		result := tx.Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).Create(record)
		if result.Error != nil {
			return result.Error
		}
		// 无论新领取还是已存在的批次，都推进计划；不重跑失败或中断的批次。
		updates := map[string]interface{}{"next_run_time": next}
		if result.RowsAffected == 1 {
			updates["last_run_time"] = now
		}
		if err := tx.Model(&current).UpdateColumns(updates).Error; err != nil {
			return err
		}
		if result.RowsAffected == 1 {
			claimed = &execution{id: record.Id}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if claimed == nil {
		return nil, nil
	}
	return claimed, nil
}

func (run *execution) Finish(ctx context.Context, taskErr error) error {
	status, message := "success", ""
	if taskErr != nil {
		status, message = "failed", taskErr.Error()
	}
	if errors.Is(taskErr, context.Canceled) {
		status = "interrupted"
	}
	result := state.DB().WithContext(ctx).Model(&SysCronJobLog{}).Where("id = ? AND status = ?", run.id, "running").
		Updates(map[string]interface{}{"status": status, "error": message, "end_time": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("执行记录不存在或已结束")
	}
	return nil
}
