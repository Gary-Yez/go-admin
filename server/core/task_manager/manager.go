package task_manager

import (
	"errors"
	"sync"
)

type defaultManage struct {
	lock      sync.Mutex
	handlers  []*handler
	functions map[string]handlerFun
}

func (t *defaultManage) GetHandlers() []*handler {
	return t.handlers
}

func (t *defaultManage) RegisterHandler(key string, option *HandleOption, fn handlerFun) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok := t.functions[key]
	if ok {
		return errors.New("TaskKey already registered")
	}
	t.handlers = append(t.handlers, &handler{
		Key:          key,
		HandleOption: option,
	})
	t.functions[key] = fn
	return nil
}

//func (t *defaultManage) Run(id uint, key string, params interface{}) error {
//	return t.functions[key](params)
//	//go func() {
//	//	err := t.functions[key](params)
//	//	if err != nil {
//	//		return
//	//	}
//	//	t.result <- err
//	//}()
//}
