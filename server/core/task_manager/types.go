package task_manager

import "sync"

type handlerFun func(params interface{}) error

type Manager interface {
	GetHandlers() []*handler
	RegisterHandler(key string, option *HandleOption, fn handlerFun) error
}

type HandlerParams struct {
	Key         string      `json:"key"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default"`
}
type HandleOption struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Params      []*HandlerParams `json:"params"`
}

type handler struct {
	Key string `json:"key"`
	*HandleOption
}

func NewManager() Manager {
	return &defaultManage{
		lock:     sync.Mutex{},
		handlers: []*handler{},
		//functions: make(map[string]taskFunc),
	}
}
