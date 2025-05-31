package modular

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
)

func NewLoader() Loader {
	return &defaultLoader{
		sequence: []string{},
		mp:       make(map[string]Mounter),
		lock:     sync.Mutex{},
	}
}

type defaultLoader struct {
	sequence []string
	mp       map[string]Mounter
	lock     sync.Mutex
}

func (m *defaultLoader) Mount(routerPrefix string, mounter Mounter) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.mp[routerPrefix]; ok {
		panic("module routerPrefix: " + routerPrefix + " has been existed")
	}
	m.sequence = append(m.sequence, routerPrefix)
	m.mp[routerPrefix] = mounter
}

func (m *defaultLoader) InitializeAll() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, moduleName := range m.sequence {
		err := m.mp[moduleName].Initialize()
		if err != nil {
			return errors.New(moduleName + " initialize fail: " + err.Error())
		}
	}
	return nil
}

func (m *defaultLoader) RegisterRouter(adminGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, moduleName := range m.sequence {
		m.mp[moduleName].AdminRouter(adminGroup.Group(moduleName))
		m.mp[moduleName].PublicRouter(publicGroup.Group(moduleName))
	}
	return nil
}
