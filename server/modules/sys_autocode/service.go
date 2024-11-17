package sys_autocode

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/mxcker/go-admin/server/global"
	"gitee.com/mxcker/go-admin/server/models"
	"gitee.com/mxcker/go-admin/server/models/request"
	"gorm.io/gorm"
	"path/filepath"
)

type Service struct {
}

func (s *Service) GetTemplates(data *GenerateBody) ([]WriteItem, error) {
	var templateList []WriteItem
	for _, templateItem := range ServerTemplatesPath {
		content, err := getTemplateContent(templateItem.Path, data)
		if err != nil {
			return nil, err
		}
		var filePosition string
		if templateItem.Name == "model" {
			filePosition = filepath.Join(ServerPath, "./models", data.ModuleName+".go")
		} else {
			filePosition = filepath.Join(ServerPath, "./modules", data.ModuleName, templateItem.Name+".go")
		}
		templateList = append(templateList, WriteItem{
			Path:    filePosition,
			Content: content,
		})
	}
	routerContent, routerPath, err := getRouterContent(data.ModuleName)
	if err != nil {
		return nil, err
	}
	templateList = append(templateList, WriteItem{
		Path:    routerPath,
		Content: routerContent,
	})
	initializationContent, initPath, err := getInitializationContent(data.ModelName)
	if err != nil {
		return nil, err
	}
	templateList = append(templateList, WriteItem{
		Path:    initPath,
		Content: initializationContent,
	})
	for _, templateItem := range WebTemplatesPath {
		content, err := getTemplateContent(templateItem.Path, data)
		if err != nil {
			return nil, err
		}
		var filePosition string
		if templateItem.Name == "api" {
			filePosition = filepath.Join(WebPath, "./apis", data.ModuleName+".ts")
		} else {
			filePosition = filepath.Join(WebPath, "./views", data.ModuleName, "index.vue")
		}
		templateList = append(templateList, WriteItem{
			Path:    filePosition,
			Content: content,
		})
		//templates[filePosition] = content
	}
	return templateList, nil
}

func (s *Service) Generate(data *GenerateBody) error {
	templateMap, err := s.GetTemplates(data)
	if err != nil {
		return err
	}
	for _, templateItem := range templateMap {
		err = writeFile(filepath.Join(RootPath, templateItem.Path), templateItem.Content)
		if err != nil {
			fmt.Println(templateItem.Path, err.Error())
			return err
		}
	}
	return err
}

func (s *Service) History(req *request.ReqList) (list []*models.SysAutoCode, total int64, err error) {
	db := global.DB.Model(models.SysAutoCode{})
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	req.SetDB(db)
	err = db.Find(&list).Error
	return
}

func (s *Service) Get(id uint) (data *models.SysAutoCode, err error) {
	data = &models.SysAutoCode{}
	err = global.DB.Model(models.SysAutoCode{}).Where("id = ?", id).Find(data).Error
	return
}

func (s *Service) SaveHistory(data *GenerateBody) error {
	history := models.SysAutoCode{
		ModuleName: data.ModuleName,
	}
	if err := global.DB.Where(history).First(&history).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	history.ModelName = data.ModelName
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	history.Form = string(marshal)
	return global.DB.Save(&history).Error
}

func (s *Service) DeleteByIds(ids []uint) (err error) {
	err = global.DB.Delete(&models.SysAutoCode{}, ids).Error
	return
}
