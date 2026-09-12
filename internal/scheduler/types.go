package scheduler

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"time"
)

const StringParams = ParamsType("string")
const IntParams = ParamsType("int")
const BoolParams = ParamsType("bool")

// DefaultLocation 统一使用北京时间，避免执行周期依赖服务器的本地时区。
var DefaultLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

type handlerFun func(ctx context.Context, params []byte) error

type SchedulerOption struct {
	JobSyncer         JobSyncer
	DistributedLocker gocron.Locker
	Location          *time.Location
}

type ParamsType string

type HandlerParam struct {
	Name        string     `json:"name"`
	Key         string     `json:"key"`
	Type        ParamsType `json:"type"`
	Description string     `json:"description"`
	Required    bool       `json:"required"`
}

type HandlerParams []*HandlerParam

type HandlerOption struct {
	Name    string        `json:"name"`
	Params  HandlerParams `json:"params"`
	Handler handlerFun    `json:"-"`
}
