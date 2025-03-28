package global

import (
	"os"
	"path/filepath"
)

func Init() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	RootPath = filepath.Join(cwd, "../")
	err = InitConfig()
	if err != nil {
		return err
	}
	err = InitDB()
	if err != nil {
		return err
	}
	return err
}
