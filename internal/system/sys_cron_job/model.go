package sys_cron_job

import (
	"gitee.com/mxcker/go-admin/internal/state"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type SysCronJob struct {
	Id          uint       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	Name        string     `json:"name" binding:"required"`
	HandlerKey  string     `json:"handler_key" binding:"required"`
	Params      string     `json:"params"`
	Cron        string     `json:"cron" binding:"required"`
	Enable      bool       `json:"enable"`
	LastRunTime *time.Time `json:"last_run_time"`
	NextRunTime *time.Time `json:"next_run_time"`
	Version     int        `json:"-" gorm:"default:1"`
}

func (j *SysCronJob) GetHandlerKey() string {
	return j.HandlerKey
}

func (j *SysCronJob) GetCronExpr() string {
	return j.Cron
}

func (j *SysCronJob) GetParams() []byte {
	return []byte(j.Params)
}

func (j *SysCronJob) GetVersion() int {
	return j.Version
}

func (j *SysCronJob) GetID() string {
	return strconv.Itoa(int(j.Id))
}

func (j *SysCronJob) GetLastRunTime() time.Time {
	if j.LastRunTime != nil {
		return *j.LastRunTime
	}
	return time.Time{}
}

func (j *SysCronJob) GetNextRunTime() time.Time {
	if j.NextRunTime != nil {
		return *j.NextRunTime
	}
	return time.Time{}
}

func (j *SysCronJob) AfterRun(startTime time.Time, endTime time.Time, nextRunTime time.Time, err error) {
	state.DB().Model(j).UpdateColumns(map[string]interface{}{
		"last_run_time": startTime,
		"next_run_time": nextRunTime,
	})
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	log := &SysCronJobLog{
		StartTime: startTime,
		EndTime:   endTime,
		Error:     errMsg,
		JobId:     j.Id,
	}
	state.DB().Omit(clause.Associations).Create(log)
}

type SysCronJobLog struct {
	Id        uint        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time   `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"comment:更新时间"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Error     string      `json:"error"`
	JobId     uint        `json:"job_id"`
	Job       *SysCronJob `json:"job" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
