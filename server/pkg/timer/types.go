package timer

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"time"
)

const StringParams = ParamsTye("string")
const IntParams = ParamsTye("int")
const BoolParams = ParamsTye("bool")

type handlerFun func(ctx context.Context, params []byte) error

type SchedulerOption struct {
	JobSyncer         JobSyncer
	DistributedLocker gocron.Locker
	Location          *time.Location
}

type ParamsTye string

type handlerParam struct {
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Type        ParamsTye `json:"type"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
}

type HandlerParams []*handlerParam

type HandlerOption struct {
	Name    string        `json:"name"`
	Params  HandlerParams `json:"params"`
	Handler handlerFun    `json:"-"`
}
