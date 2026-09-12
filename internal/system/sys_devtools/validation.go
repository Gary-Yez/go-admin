package sys_devtools

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"gorm.io/gorm/schema"
)

var modulePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
var exportedPattern = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
var jsonNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
var devicePattern = regexp.MustCompile(`^(COM|LPT)[0-9]$`)

var supportedTypes = map[string]bool{
	"string": true, "bool": true, "int": true, "int8": true,
	"int16": true, "int32": true, "int64": true, "uint": true,
	"uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"float32": true, "float64": true, "time.Time": true,
}

func requiredPointer(f Filed) bool {
	return f.Required && f.Type != "string" && f.Type != "time.Time"
}

func columnName(field Filed) string {
	return (schema.NamingStrategy{}).ColumnName("", field.Key)
}

func safeLabel(s string) bool {
	return !strings.ContainsAny(s, "\"'`<>\\&${}") && !strings.ContainsFunc(s, unicode.IsControl)
}

func validateGenerateBody(data *GenerateBody) error {
	if data == nil {
		return fmt.Errorf("生成配置不能为空")
	}
	name := data.ModuleName
	reservedModules := " admin main init internal vendor any comparable bool byte complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr true false iota nil append cap clear close complex copy delete imag len make max min new panic print println real recover "
	if len(name) > 64 || !modulePattern.MatchString(name) || token.Lookup(name).IsKeyword() || strings.Contains(reservedModules, " "+name+" ") {
		return fmt.Errorf("模块名必须是非关键字的小写字母、数字或下划线标识符，长度不超过64")
	}
	upper := strings.ToUpper(name)
	if upper == "CON" || upper == "PRN" || upper == "AUX" || upper == "NUL" || devicePattern.MatchString(upper) {
		return fmt.Errorf("模块名不能使用系统保留文件名")
	}
	if len(data.ModelName) > 64 || !exportedPattern.MatchString(data.ModelName) {
		return fmt.Errorf("模型名必须以大写英文字母开头，仅包含字母和数字，长度不超过64")
	}
	switch data.ModelName {
	case "Mounter", "Controller", "Service":
		return fmt.Errorf("模型名 %s 与模块内置名称冲突", data.ModelName)
	}
	if strings.TrimSpace(data.ChineseModuleName) == "" || !safeLabel(data.ChineseModuleName) {
		return fmt.Errorf("模块显示名称不能为空，且不能包含引号、模板符号或控制字符")
	}
	if len(data.Fields) > 200 {
		return fmt.Errorf("字段数量不能超过200")
	}
	keys := map[string]bool{"id": true, "createdat": true, "updatedat": true, "deletedat": true, "tomodel": true}
	names := map[string]bool{"id": true, "created_at": true, "updated_at": true, "deleted_at": true}
	columns := map[string]bool{"id": true, "created_at": true, "updated_at": true, "deleted_at": true}
	jsReserved := " break case catch class const continue debugger default delete do else enum export extends false finally for function if import in instanceof let new null return super switch this throw true try typeof var void while with yield await implements interface package private protected public static prototype constructor __proto__ "
	for i, f := range data.Fields {
		fail := func(message string) error { return fmt.Errorf("第%d个字段：%s", i+1, message) }
		if len(f.Key) > 64 || !exportedPattern.MatchString(f.Key) {
			return fail("Go字段名必须以大写英文字母开头，仅包含字母和数字，长度不超过64")
		}
		key := strings.ToLower(f.Key)
		column := (schema.NamingStrategy{}).ColumnName("", f.Key)
		if keys[key] || columns[column] {
			return fail("Go字段名重复或与系统字段冲突")
		}
		keys[key], columns[column] = true, true
		if len(f.Name) > 64 || !jsonNamePattern.MatchString(f.Name) || strings.Contains(jsReserved, " "+f.Name+" ") {
			return fail("JSON字段名必须是非保留的字母、数字或下划线标识符，长度不超过64")
		}
		jsonName := strings.ToLower(f.Name)
		if names[jsonName] {
			return fail("JSON字段名重复或与系统字段冲突")
		}
		names[jsonName] = true
		if !supportedTypes[f.Type] {
			return fail("不支持的字段类型：" + f.Type)
		}
		switch f.IndexType {
		case "", "index", "unique", "uniqueIndex":
		default:
			return fail("索引类型只支持空值、index、unique、uniqueIndex")
		}
		switch f.QueryType {
		case "", "like", "=", "!=", ">", ">=", "<", "<=", "between":
		default:
			return fail("不支持的筛选方式")
		}
		if f.QueryType == "like" && f.Type != "string" {
			return fail("模糊搜索仅支持字符串字段")
		}
		if f.QueryType == "between" && f.Type != "time.Time" {
			return fail("时间范围筛选仅支持日期时间字段")
		}
		if f.Type == "bool" && f.QueryType != "" && f.QueryType != "=" && f.QueryType != "!=" {
			return fail("布尔字段仅支持等于或不等于筛选")
		}
		if strings.Contains(f.ChineseName, ";") {
			return fail("显示名称不能包含分号，以免影响数据库字段注释")
		}
		if !safeLabel(f.ChineseName) {
			return fail("显示名称不能包含引号、模板符号或控制字符")
		}
		if f.Required && !f.Editable {
			return fail("请求必填字段必须可编辑")
		}
	}
	return nil
}
