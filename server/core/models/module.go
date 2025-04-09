package models

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type ModuleRegisterFunc func(path string, adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error

type ModuleMap struct {
	sequence []string
	mp       map[string]ModuleRegisterFunc
	lock     sync.Mutex
}

func (m *ModuleMap) Add(moduleName string, registerFunc ModuleRegisterFunc) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mp == nil {
		m.mp = make(map[string]ModuleRegisterFunc)
	}
	if _, ok := m.mp[moduleName]; ok {
		panic("module: " + moduleName + " has been existed")
	}
	m.sequence = append(m.sequence, moduleName)
	m.mp[moduleName] = registerFunc
}

func (m *ModuleMap) RegisterAll(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, moduleName := range m.sequence {
		err := m.mp[moduleName](moduleName, adminAuthGroup, publicGroup)
		if err != nil {
			return err
		}
	}
	return nil
}
