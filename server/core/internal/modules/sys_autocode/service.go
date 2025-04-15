package sys_autocode

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/mxcker/go-admin/server/core/global"
	request2 "gitee.com/mxcker/go-admin/server/core/types/request"
	"gorm.io/gorm"
	"path/filepath"
)

type serviceStruct struct {
}

func (s *serviceStruct) GetTemplates(data *GenerateBody) ([]WriteItem, error) {
	var templateList []WriteItem
	for _, templateItem := range ServerTemplatesPath {
		if !data.CreateCURD && templateItem.Name == "model" {
			continue
		}
		content, err := getTemplateContent(templateItem.Path, data)
		if err != nil {
			return nil, err
		}
		var filePosition string
		filePosition = filepath.Join(ServerPath, "./modules", data.ModuleName, templateItem.Name+".go")
		templateList = append(templateList, WriteItem{
			Path:    filePosition,
			Content: content,
		})
	}
	routerContent, routerPath, err := getModuleEnterContent(data.ModuleName)
	if err != nil {
		return nil, err
	}
	templateList = append(templateList, WriteItem{
		Path:    routerPath,
		Content: routerContent,
	})
	if data.CreateCURD {
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
		}
	}
	return templateList, nil
}

func (s *serviceStruct) Generate(data *GenerateBody) error {
	templateMap, err := s.GetTemplates(data)
	if err != nil {
		return err
	}
	for _, templateItem := range templateMap {
		err = writeFile(filepath.Join(templateItem.Path), templateItem.Content)
		if err != nil {
			fmt.Println(templateItem.Path, err.Error())
			return err
		}
	}
	return err
}

func (s *serviceStruct) History(req *request2.ReqList) (list []*SysAutoCode, total int64, err error) {
	db := global.DB.Model(SysAutoCode{})
	err = req.BuildWhere(db).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = req.BuildQuery(db).Find(&list).Error
	return
}

func (s *serviceStruct) Get(req *request2.Req) (data *SysAutoCode, err error) {
	data = &SysAutoCode{}
	err = req.BuildQuery(global.DB.Model(SysAutoCode{})).First(data).Error
	return
}

func (s *serviceStruct) SaveHistory(data *GenerateBody) error {
	history := SysAutoCode{
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

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	err = req.BuildQuery(global.DB).Delete(&SysAutoCode{}).Error
	return
}
