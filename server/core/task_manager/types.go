package task_manager

import "context"

type handlerFun func(ctx context.Context, params map[string]interface{}) error
type ParamsTye string

const StringParams = ParamsTye("string")
const IntParams = ParamsTye("int")
const BoolParams = ParamsTye("bool")

type Manager interface {
	GetHandlers() map[string]*HandlerOption
	RegisterHandler(key string, option *HandlerOption) error
	RegisterJob(job *ScheduledJob) error
	RemoveJobs(ids []uint)
	StartScheduler()
}

type ScheduledJob struct {
	ctx              context.Context
	cancel           context.CancelFunc
	Id               uint                   // 任务ID
	HandlerKey       string                 // 注册的任务处理函数的Key
	CronExpr         string                 // Cron表达式 * * * * * *
	Params           map[string]interface{} // 传递给任务处理函数的参数
	AllowConcurrency bool                   // 是否允许并发
}

type HandlerParams struct {
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Type        ParamsTye `json:"type"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
}
type HandlerOption struct {
	Name    string           `json:"name"`
	Params  []*HandlerParams `json:"params"`
	Handler handlerFun       `json:"-"`
}
