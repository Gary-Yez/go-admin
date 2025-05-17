package timer

import (
	"time"
)

type Scheduler interface {
	GetHandlers() map[string]*HandlerOption
	RegisterHandler(key string, option *HandlerOption) error
	StartScheduler(syncInterval time.Duration)
}

type Job interface {
	GetID() string
	GetHandlerKey() string
	GetCronExpr() string
	GetParams() []byte
	GetVersion() int
	GetLastRunTime() time.Time
	GetNextRunTime() time.Time
	AfterRun(startTime time.Time, endTime time.Time, nextRunTime time.Time, err error)
}

type JobSyncer interface {
	Sync() ([]Job, error)
}
