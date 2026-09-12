package scheduler

import (
	"context"
	"time"
)

type Scheduler interface {
	GetHandlers() map[string]*HandlerOption
	RegisterHandler(key string, option *HandlerOption) error
	StartScheduler(syncInterval time.Duration)
	StopScheduler() error
}

type Job interface {
	GetID() string
	GetHandlerKey() string
	GetCronExpr() string
	GetParams() []byte
	GetVersion() int
	GetNextRunTime() time.Time
	Claim(context.Context) (Run, error)
}

type Run interface {
	Finish(context.Context, error) error
}

type JobSyncer interface {
	Sync() ([]Job, error)
}
