package module_loader

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type defaultLoader struct {
	sequence []string
	mp       map[string]Mounter
	lock     sync.Mutex
}

func (m *defaultLoader) Add(moduleName string, mounter Mounter) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.mp[moduleName]; ok {
		panic("module: " + moduleName + " has been existed")
	}
	m.sequence = append(m.sequence, moduleName)
	m.mp[moduleName] = mounter
}

func (m *defaultLoader) InitAll() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, moduleName := range m.sequence {
		err := m.mp[moduleName].Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *defaultLoader) RegisterAll(adminAuthGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, moduleName := range m.sequence {
		err := m.mp[moduleName].Register(moduleName, adminAuthGroup, publicGroup)
		if err != nil {
			return err
		}
	}
	return nil
}
