package sys_config

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/bizconfig"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/Gary-Yez/go-admin/internal/system/sys_devtools"
)

type serviceStruct struct{}

var writeMu sync.Mutex

var baseConfigSchema *bizconfig.Schema

// SetBaseSchema 在启动时提供内置字段，供生成器检查冲突。
func SetBaseSchema(schema *bizconfig.Schema) { baseConfigSchema = schema }

func hash(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(content))
}

func configPath() (string, error) {
	path, err := filepath.Abs(filepath.Join(sys_devtools.ServerPath, "settings", "config.go"))
	if err != nil {
		return "", err
	}
	for current := path; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil && !(current == path && os.IsNotExist(err)) {
			return "", err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("配置路径不能经过符号链接")
		}
		if current == filepath.Dir(current) {
			break
		}
	}
	return path, nil
}

func load() (DefinitionFile, string, []byte, error) {
	definition := DefinitionFile{Fields: []Field{}, Groups: []string{}}
	path, err := configPath()
	if err != nil {
		return definition, "", nil, err
	}
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return definition, "", nil, err
	}
	definition.Hash = hash(content)
	definition.Fields, definition.typeName, err = parse(content)
	return definition, path, content, err
}

func (*serviceStruct) List() (*DefinitionFile, error) {
	definition, _, _, err := load()
	if err != nil {
		return nil, err
	}
	stored, err := bizconfig.Definitions()
	if err != nil {
		return nil, err
	}
	definitions := make(map[string]bizconfig.Definition, len(stored))
	groups := map[string]bool{}
	for _, row := range stored {
		definitions[row.Key] = row
		if row.Group != "" {
			groups[row.Group] = true
		}
	}
	definition.BuiltinGroups = []string{}
	for _, field := range baseConfigSchema.Fields {
		group := field.Group
		if stored, ok := definitions[field.Key]; ok {
			group = stored.Group
		}
		if !slices.Contains(definition.BuiltinGroups, group) {
			definition.BuiltinGroups = append(definition.BuiltinGroups, group)
		}
	}
	for i := range definition.Fields {
		field := &definition.Fields[i]
		if row, exists := definitions[field.Key]; exists {
			if row.Type != field.Type {
				return nil, fmt.Errorf("配置 %s 的代码和数据库类型不一致", field.Key)
			}
			field.Definition = row
		}
		if field.Group != "" {
			groups[field.Group] = true
		}
	}
	for group := range groups {
		definition.Groups = append(definition.Groups, group)
	}
	slices.Sort(definition.Groups)
	return &definition, nil
}

func (*serviceStruct) Preview(body *SaveBody) (*Preview, error) {
	definition, path, before, err := load()
	if err != nil {
		return nil, err
	}
	if body.Hash != definition.Hash {
		return nil, errors.New("配置文件已变化，请刷新后重新操作")
	}
	if err := prepareFields(body.Fields, definition.Fields); err != nil {
		return nil, err
	}
	after, err := render(body.Fields, definition.typeName)
	if err != nil {
		return nil, err
	}
	action := "modify"
	if len(before) == 0 {
		action = "create"
	} else if string(before) == string(after) {
		action = "unchanged"
	}
	return &Preview{Path: path, Before: string(before), After: string(after), Action: action}, nil
}
func (s *serviceStruct) Apply(body *SaveBody) error {
	writeMu.Lock()
	defer writeMu.Unlock()
	preview, err := s.Preview(body)
	if err != nil {
		return err
	}
	definitions := make([]bizconfig.Definition, 0, len(body.Fields))
	for i, field := range body.Fields {
		definition := field.Definition
		definition.Sort = len(baseConfigSchema.Fields) + i
		definitions = append(definitions, definition)
	}
	fileApplied := false
	err = bizconfig.SaveDefinitions(definitions, func() error {
		if err := applyDefinitionFile(preview); err != nil {
			return err
		}
		fileApplied = preview.Action != "unchanged"
		return nil
	})
	if err != nil && fileApplied {
		current, readErr := os.ReadFile(preview.Path)
		if os.IsNotExist(readErr) {
			readErr = nil
		}
		if readErr != nil || string(current) != preview.After {
			return errors.Join(err, errors.New("数据库提交失败且文件发生变化，请手动核对配置文件"))
		}
		var restoreErr error
		if preview.Before == "" {
			restoreErr = os.Remove(preview.Path)
		} else {
			restoreErr = os.WriteFile(preview.Path, []byte(preview.Before), 0644)
		}
		return errors.Join(err, restoreErr)
	}
	return err
}

func applyDefinitionFile(preview *Preview) error {
	current, err := os.ReadFile(preview.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if string(current) != preview.Before {
		return errors.New("配置文件已变化，请重新预览")
	}
	if preview.Action == "unchanged" {
		return nil
	}
	file, err := os.CreateTemp(filepath.Dir(preview.Path), ".config-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err = file.WriteString(preview.After); err != nil {
		file.Close()
		return err
	}
	if err = file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	if err = os.Chmod(file.Name(), 0644); err != nil {
		return err
	}
	// 写入前再次核对，避免覆盖预览后被其他工具修改的文件。
	current, err = os.ReadFile(preview.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if string(current) != preview.Before {
		return errors.New("配置文件已变化，请重新预览")
	}
	return os.Rename(file.Name(), preview.Path)
}

func (*serviceStruct) Site() (*SiteInfo, error) {
	info := new(SiteInfo)
	// 公开字段固定列出，新增配置或调整分组不会扩大公开范围。
	for _, field := range []struct {
		key   string
		value *string
	}{
		{"site.name", &info.Name},
		{"site.logo", &info.Logo},
		{"site.favicon", &info.Favicon},
		{"site.title", &info.Title},
		{"site.copyright", &info.Copyright},
	} {
		value, err := bizconfig.ReadValue(field.key)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(value, field.value); err != nil {
			return nil, err
		}
	}
	return info, nil
}
