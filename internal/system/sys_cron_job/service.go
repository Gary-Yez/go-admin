package sys_cron_job

import (
	"errors"
	"github.com/Gary-Yez/go-admin/internal/state"
	"time"

	"github.com/Gary-Yez/go-admin/internal/scheduler"
	request2 "github.com/Gary-Yez/go-admin/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) Sync() ([]scheduler.Job, error) {
	var list []*SysCronJob
	if err := state.DB().Model(SysCronJob{}).Where("enable = ?", true).Find(&list).Error; err != nil {
		return nil, err
	}
	var jobs []scheduler.Job
	for _, task := range list {
		jobs = append(jobs, task)
	}
	return jobs, nil
}

func (s *serviceStruct) GetLogs(req *request2.ReqList) (list []*SysCronJobLog, total int64, err error) {
	db := req.WithFilter(state.DB().Model(SysCronJobLog{}), []string{"job_id"})
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, nil)).Order("id DESC").Find(&list).Error
	return
}

func (s *serviceStruct) List(req *request2.ReqList) (list []*SysCronJob, total int64, err error) {
	db := req.WithFilter(state.DB().Model(SysCronJob{}), nil)
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, []string{"id"})).Find(&list).Error
	return
}

func (s *serviceStruct) Create(data *SysCronJob) (err error) {
	next, err := nextRun(data.Cron, time.Now())
	if err != nil {
		return err
	}
	data.NextRunTime, data.LastRunTime, data.Version = &next, nil, 1
	err = state.DB().
		Omit(clause.Associations).
		Create(data).Error
	return err
}

func (s *serviceStruct) Update(data *SysCronJob) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	next, err := nextRun(data.Cron, time.Now())
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"version":       gorm.Expr("version + 1"),
		"name":          data.Name,
		"handler_key":   data.HandlerKey,
		"params":        data.Params,
		"cron":          data.Cron,
		"enable":        data.Enable,
		"next_run_time": next,
	}
	err = state.DB().Model(data).Where("id = ?", data.Id).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	if err = req.WithQuery(state.DB()).Delete(&SysCronJob{}).Error; err != nil {
		return err
	}
	return nil
}
