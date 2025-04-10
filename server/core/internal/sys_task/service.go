package sys_task

import (
	"encoding/json"
	"errors"
	"gitee.com/mxcker/go-admin/server/core/global"
	"gitee.com/mxcker/go-admin/server/core/task_manager"
	"gitee.com/mxcker/go-admin/server/core/types/request"
	"gorm.io/gorm/clause"
)

type serviceStruct struct {
}

func (s *serviceStruct) Sync() error {
	var list []*SysTask
	if err := global.DB.Model(SysTask{}).Find(&list).Error; err != nil {
		return err
	}
	for _, task := range list {
		if task.Enable {
			_ = s.RegisterJob(task)
		}
	}
	return nil
}

func (s *serviceStruct) RegisterJob(task *SysTask) error {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(task.Params), &params)
	if err != nil {
		return err
	}
	return global.TaskManager.RegisterJob(&task_manager.ScheduledJob{
		Id:               task.Id,
		HandlerKey:       task.HandlerKey,
		CronExpr:         task.Cron,
		Params:           params,
		AllowConcurrency: task.AllowConcurrency,
	})
}

func (s *serviceStruct) Get(req *request.Req) (data *SysTask, err error) {
	data = &SysTask{}
	err = req.BuildQuery(global.DB.Model(SysTask{})).First(data).Error
	return
}

func (s *serviceStruct) List(req *request.ReqList) (list []*SysTask, total int64, err error) {
	db := global.DB.Model(SysTask{})
	err = req.BuildWhere(db).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.BuildQuery(db).Find(&list).Error
	return
}

func (s *serviceStruct) Create(data *SysTask) (err error) {
	err = global.DB.
		Omit(clause.Associations).
		Create(data).Error
	if err != nil {
		return err
	}
	if data.Enable {
		return s.RegisterJob(data)
	}
	return
}

func (s *serviceStruct) Update(data *SysTask) (err error) {
	if data.Id == 0 {
		return errors.New("id不能为空")
	}
	err = global.DB.Select("*").
		Omit(clause.Associations).
		Omit("Id", "CreatedAt", "UpdatedAt").
		Where("id = ?", data.Id).Updates(data).Error
	if err != nil {
		return err
	}
	global.TaskManager.RemoveJobs([]uint{data.Id})
	if data.Enable {
		return s.RegisterJob(data)
	}
	return nil
}

func (s *serviceStruct) DeleteByIds(req *request.ReqIds) (err error) {
	if err = req.BuildQuery(global.DB).Delete(&SysTask{}).Error; err != nil {
		return err
	}
	global.TaskManager.RemoveJobs(req.Ids)
	return
}
