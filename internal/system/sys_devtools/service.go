package sys_devtools

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Gary-Yez/go-admin/internal/state"
	request2 "github.com/Gary-Yez/go-admin/request"
	"gorm.io/gorm"
)

type serviceStruct struct {
}

var generationMu sync.Mutex

func (s *serviceStruct) GetTemplates(data *GenerateBody) ([]WriteItem, error) {
	if err := validateGenerateBody(data); err != nil {
		return nil, err
	}
	var templateList []WriteItem
	if err := prepareMenus(data); err != nil {
		return nil, err
	}
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
	routerPath, err := checkedPath(filepath.Join(ServerPath, "modules", "enter.go"))
	if err != nil {
		return nil, err
	}
	originalRouter, err := os.ReadFile(routerPath)
	if err != nil {
		return nil, err
	}
	routerContent, routerPath, err := getModuleEnterContent(data.ModuleName, originalRouter)
	if err != nil {
		return nil, err
	}
	templateList = append(templateList, WriteItem{
		Path:         routerPath,
		Content:      routerContent,
		registration: true,
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
	items, err := inspectFiles(templateList)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.registration && item.ExistingHash != contentHash(originalRouter) {
			return nil, fmt.Errorf("模块注册文件已变化，请重新预览")
		}
	}
	return items, nil
}

func (s *serviceStruct) Generate(data *GenerateBody) error {
	generationMu.Lock()
	defer generationMu.Unlock()
	templateMap, err := s.GetTemplates(data)
	if err != nil {
		return err
	}
	return commitFiles(templateMap, data.OverwriteFiles, func() error {
		return s.SaveHistory(data)
	})
}

func (s *serviceStruct) History(req *HistoryQuery) (list []HistoryItem, total int64, err error) {
	db := state.DB().Model(&SysAutoCode{})
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		// Treat search input literally rather than as SQL LIKE wildcards.
		keyword = strings.NewReplacer("!", "!!", "%", "!%", "_", "!_").Replace(keyword)
		pattern := "%" + keyword + "%"
		db = db.Where("(module_name LIKE ? ESCAPE '!' OR model_name LIKE ? ESCAPE '!')", pattern, pattern)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	var histories []SysAutoCode
	err = db.Order("updated_at DESC").Order("id DESC").Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Find(&histories).Error
	if err != nil {
		return nil, total, err
	}
	list = make([]HistoryItem, 0, len(histories))
	for _, history := range histories {
		item := HistoryItem{Id: history.Id, CreatedAt: history.CreatedAt, UpdatedAt: history.UpdatedAt, ModuleName: history.ModuleName, ModelName: history.ModelName}
		var config GenerateBody
		if json.Unmarshal([]byte(history.Form), &config) == nil && config.ModuleName == history.ModuleName && config.ModelName == history.ModelName && (!config.CreateCURD || config.Fields != nil) {
			item.ChineseModuleName = config.ChineseModuleName
			item.CreateCURD, item.UseSoftDelete = config.CreateCURD, config.UseSoftDelete
			if config.CreateCURD {
				item.FieldCount = len(config.Fields)
			}
			item.ConfigValid = true
			if config.CreateCURD {
				for _, field := range config.Fields {
					if field.Name == "" || field.Key == "" || field.Type == "" {
						item.ConfigValid = false
						break
					}
				}
			}
		}
		list = append(list, item)
	}
	return
}

func (s *serviceStruct) GetHistory(req *request2.Req) (*SysAutoCode, error) {
	history := new(SysAutoCode)
	err := req.WithQuery(state.DB()).First(history).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("生成记录不存在或已被删除，请刷新列表")
	}
	return history, err
}

func (s *serviceStruct) SaveHistory(data *GenerateBody) error {
	history := SysAutoCode{
		ModuleName: data.ModuleName,
	}
	if err := state.DB().Where(history).First(&history).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	history.ModelName = data.ModelName
	saved := *data
	saved.OverwriteFiles = nil
	marshal, err := json.Marshal(&saved)
	if err != nil {
		return err
	}
	history.Form = string(marshal)
	return state.DB().Save(&history).Error
}

func (s *serviceStruct) DeleteByIds(req *request2.ReqIds) (err error) {
	if len(req.Ids) == 0 {
		return errors.New("请选择要删除的生成记录")
	}
	err = req.WithQuery(state.DB()).Delete(&SysAutoCode{}).Error
	return
}
