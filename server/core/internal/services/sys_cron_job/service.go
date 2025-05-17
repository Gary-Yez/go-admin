package sys_cron_job

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/core/pkg/timer"
	"gitee.com/mxcker/go-admin/server/global"
	request2 "gitee.com/mxcker/go-admin/server/types/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) Sync() ([]timer.Job, error) {
	var list []*SysCronJob
	if err := global.DB.Model(SysCronJob{}).Find(&list).Error; err != nil {
		return nil, err
	}
	var jobs []timer.Job
	for _, task := range list {
		if task.Enable {
			jobs = append(jobs, task)
		}
	}
	return jobs, nil
}

func (s *serviceStruct) GetLogs(req *request2.ReqList) (list []*SysCronJobLog, total int64, err error) {
	db := req.WithFilter(global.DB.Model(SysCronJobLog{}), []string{"job_id"})
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, nil)).Order("id DESC").Find(&list).Error
	return
}

func (s *serviceStruct) Get(req *request2.Req) (data *SysCronJob, err error) {
	data = &SysCronJob{}
	err = req.WithQuery(global.DB.Model(SysCronJob{})).First(data).Error
	return
}

func (s *serviceStruct) List(req *request2.ReqList) (list []*SysCronJob, total int64, err error) {
	db := req.WithFilter(global.DB.Model(SysCronJob{}), nil)
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.WithPagination(req.WithSort(db, nil)).Find(&list).Error
	return
}

func (s *serviceStruct) Create(data *SysCronJob) (err error) {
	err = global.DB.
		Omit(clause.Associations).
		Create(data).Error
	return err
}

func (s *serviceStruct) Update(data *SysCronJob) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	updates := map[string]interface{}{
		"version":       gorm.Expr("version + 1"),
		"name":          data.Name,
		"handler_key":   data.HandlerKey,
		"params":        data.Params,
		"cron":          data.Cron,
		"enable":        data.Enable,
		"next_run_time": nil,
	}
	err = global.DB.Model(data).Where("id = ?", data.Id).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	if err = req.WithQuery(global.DB).Delete(&SysCronJob{}).Error; err != nil {
		return err
	}
	return nil
}
