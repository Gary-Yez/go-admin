package sys_autocode

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var RootPath string
var ServerPath string
var WebPath string
var ServerTemplatesPath []TemplateItem
var WebTemplatesPath []TemplateItem

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	RootPath = filepath.Join(cwd, "../")
	fmt.Println(RootPath)
	ServerPath = "admin-server"
	WebPath = "admin-web/src"
	// 服务端模板
	serverTemplateNames := []string{
		"model",
		"controller",
		"service",
		"route",
	}
	for _, name := range serverTemplateNames {
		ServerTemplatesPath = append(ServerTemplatesPath, TemplateItem{
			Name: name,
			Path: filepath.Join(RootPath, ServerPath, "./templates/server", name+".go.tmpl"),
		})
	}
	// Web端模板
	WebTemplatesPath = append(WebTemplatesPath, TemplateItem{
		Name: "api",
		Path: filepath.Join(RootPath, ServerPath, "./templates/web/api.ts.tmpl"),
	})
	WebTemplatesPath = append(WebTemplatesPath, TemplateItem{
		Name: "view",
		Path: filepath.Join(RootPath, ServerPath, "./templates/web/view.vue.tmpl"),
	})
}

func convertToPascalCase(input string) string {
	// 将下划线替换为空格，然后分割成单词
	re := regexp.MustCompile(`_+`)
	input = re.ReplaceAllString(input, " ")

	// 转换成大写驼峰命名格式
	var output string
	for _, word := range strings.Fields(input) {
		// 第一个单词首字母大写，其他单词首字母也大写
		output += strings.Title(word)
	}

	return output
}

func getTemplateContent(templatePath string, data *GenerateBody) (string, error) {
	// 打开模板文件
	var buffer bytes.Buffer
	tmpl, err := template.ParseFiles(templatePath)
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func getRouterContent(moduleName string) (string, string, error) {
	routerPath := filepath.Join(ServerPath, "./router/enter.go")
	routerEnterContent, err := os.ReadFile(filepath.Join(RootPath, routerPath))
	if err != nil {
		return "", "", err
	}
	// 定义正则表达式，匹配 `import` 块
	regImport := regexp.MustCompile(`(?s)import\s+\(([^)]+)\)`)
	matchesImport := regImport.FindStringSubmatch(string(routerEnterContent))
	if len(matchesImport) < 2 {
		fmt.Println("路由文件格式不正确")
		return "", "", err
	}
	// 构造要插入的 `import` 行
	newImport := "\t\"gitee.com/mxcker/go-admin/server/modules/" + moduleName + "\"\n"
	// 检查 `import` 块中是否已存在
	if strings.Contains(matchesImport[1], newImport) {
		fmt.Println("模块已在 import 中存在")
	} else {
		// 追加新模块到 `import` 块
		updatedImport := matchesImport[1] + newImport
		routerEnterContent = []byte(regImport.ReplaceAllString(string(routerEnterContent), `import (`+updatedImport+`)`))
	}
	// 定义正则表达式，匹配 `var routes = []Route{` 后的内容
	regRoutes := regexp.MustCompile(`(?s)var routes = \[]Route{([^}]*)}`)
	matchesRoutes := regRoutes.FindStringSubmatch(string(routerEnterContent))
	if len(matchesRoutes) < 2 {
		fmt.Println("路由文件格式不正确")
		return "", "", err
	}
	// 构造要插入的 `routes` 行
	newRoute := "\tnew(" + moduleName + ".Route),\n"
	// 检查 `routes` 块中是否已存在
	if strings.Contains(matchesRoutes[1], newRoute) {
		fmt.Println("模块已在 routes 中存在")
	} else {
		// 追加新路由到 `routes` 块
		updatedRoutes := matchesRoutes[1] + newRoute
		routerEnterContent = []byte(regRoutes.ReplaceAllString(string(routerEnterContent), `var routes = []Route{`+updatedRoutes+`}`))
	}
	return string(routerEnterContent), routerPath, nil
}

func getInitializationContent(modelName string) (string, string, error) {
	initializationPath := filepath.Join(ServerPath, "./initialization/enter.go")
	routerEnterContent, err := os.ReadFile(filepath.Join(RootPath, initializationPath))
	if err != nil {
		return "", "", err
	}
	// 正则表达式匹配 AutoMigrate 调用的内容
	re := regexp.MustCompile(`(?s)(global\.DB\.AutoMigrate\(\s*.*?\))`)
	matchesImport := re.FindStringSubmatch(string(routerEnterContent))
	if len(matchesImport) < 2 {
		fmt.Println("初始化文件格式不正确")
		return "", "", err
	}
	fmt.Println(matchesImport[1])
	// 替换匹配到的内容，插入新行
	newLine := "&models." + modelName + "{},"
	updatedContent := re.ReplaceAllStringFunc(string(routerEnterContent), func(match string) string {
		// 判断是否已包含新行
		if strings.Contains(match, newLine) {
			return match // 已包含则直接返回原内容
		}
		// 在匹配到的 AutoMigrate 调用中，最后一个 `)` 前插入新行
		return regexp.MustCompile(`\)`).ReplaceAllString(match, fmt.Sprintf("\t%s\n\t)", newLine))
	})
	//// 使用正则替换原始的 routes 部分
	//updatedContent := re.ReplaceAllString(string(routerEnterContent), updatedImport)
	return updatedContent, initializationPath, nil
}

func writeFile(filePath, fileContent string) error {
	// 获取文件的目录部分
	dir := filePath[:len(filePath)-len("/"+getFileName(filePath))]
	// 使用 os.MkdirAll 创建目录，如果目录已经存在不会报错
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// 写入文件内容
	err = os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("File written successfully:", filePath)
	return nil
}

func getFileName(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[i+1:]
		}
	}
	return filePath
}
