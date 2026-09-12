package sys_cron_job

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	_ "time/tzdata" // 保证未安装时区数据的部署环境也能解析默认任务的上海时区。
)

type initializationRecord struct {
	Key       string `gorm:"primaryKey;size:100"`
	CreatedAt time.Time
}

func (initializationRecord) TableName() string { return "sys_initializations" }

func InitData() error {
	return state.DB().Transaction(func(tx *gorm.DB) error {
		// 主键冲突时不更新记录；并发启动只有成功插入的实例负责填充。
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&initializationRecord{Key: "default_cron_jobs"})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		jobs := []SysCronJob{
			{Name: "清理计划任务日志", HandlerKey: "sys_clear_cron_logs", Params: `{"day":30}`, Cron: "0 3 * * *", Enable: true},
			{Name: "清理登录日志", HandlerKey: "sys_clear_login_logs", Params: `{"days":30}`, Cron: "10 3 * * *", Enable: true},
		}
		now := time.Now()
		for i := range jobs {
			job := &jobs[i]
			var count int64
			if err := tx.Model(&SysCronJob{}).Where("handler_key = ?", job.HandlerKey).Count(&count).Error; err != nil {
				return err
			}
			// 已有任务（包括禁用任务）保留管理员设置。
			if count > 0 {
				continue
			}
			next, err := nextRun(job.Cron, now)
			if err != nil {
				return err
			}
			job.NextRunTime, job.Version = &next, 1
			if err := tx.Create(job).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
