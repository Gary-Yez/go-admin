package sys_devtools

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	adminTemplates "github.com/Gary-Yez/go-admin/internal/templates"
)

var ServerPath string
var WebPath string
var ServerTemplatesPath []TemplateItem
var WebTemplatesPath []TemplateItem

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ServerPath = filepath.Join(cwd, "../", "server")
	WebPath = filepath.Join(cwd, "../", "web/src")
	// 服务端模板
	serverTemplateNames := []string{
		"model",
		"controller",
		"service",
		"enter",
	}
	for _, name := range serverTemplateNames {
		ServerTemplatesPath = append(ServerTemplatesPath, TemplateItem{
			Name: name,
			Path: filepath.ToSlash(filepath.Join("server", name+".go.tmpl")),
		})
	}
	// Web端模板
	WebTemplatesPath = append(WebTemplatesPath, TemplateItem{
		Name: "api",
		Path: "web/api.ts.tmpl",
	})
	WebTemplatesPath = append(WebTemplatesPath, TemplateItem{
		Name: "view",
		Path: "web/view.vue.tmpl",
	})
}

func getTemplateContent(templatePath string, data *GenerateBody) (string, error) {
	// 打开模板文件
	var buffer bytes.Buffer
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"requiredPointer": requiredPointer,
		"columnName":      columnName,
	}).ParseFS(adminTemplates.FS, templatePath)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buffer, *data)
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(templatePath, ".go.tmpl") {
		content, err := format.Source(buffer.Bytes())
		if err != nil {
			return "", fmt.Errorf("格式化模板 %s 失败: %w", templatePath, err)
		}
		return string(content), nil
	}
	return buffer.String(), nil
}

func getModuleEnterContent(moduleName string, source []byte) (string, string, error) {
	// 配置文件和位置
	fileSet := token.NewFileSet()
	filePath := filepath.Join(ServerPath, "./modules/enter.go")
	hostModule, err := readModulePath(filepath.Join(ServerPath, "go.mod"))
	if err != nil {
		return "", "", err
	}
	targetImportPath := strings.TrimSuffix(hostModule, "/") + "/modules/" + moduleName
	// 解析源代码文件
	f, err := parser.ParseFile(fileSet, filePath, source, parser.ParseComments)
	if err != nil {
		return "", "", err
	}
	adminAlias, err := ensureModuleImport(f, "github.com/Gary-Yez/go-admin", "admin")
	if err != nil {
		return "", "", err
	}
	moduleAlias, err := ensureModuleImport(f, targetImportPath, moduleName)
	if err != nil {
		return "", "", err
	}
	// 查找模块自动注册的 init 函数。
	foundInit := false
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "init" || fn.Recv != nil || fn.Body == nil {
			return true
		}
		foundInit = true
		moduleLit := &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + moduleName + `"`,
		}
		// 判断是否已经存在 admin.MustRegister("模块名", ...)
		for _, stmt := range fn.Body.List {
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			sel, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if sel.Sel.Name != "MustRegister" {
				continue
			}
			xIdent, ok := sel.X.(*ast.Ident)
			if !ok || xIdent.Name != adminAlias {
				continue
			}
			if len(callExpr.Args) > 0 {
				arg0, ok := callExpr.Args[0].(*ast.BasicLit)
				if ok && arg0.Kind == token.STRING && arg0.Value == moduleLit.Value {
					// 已经添加过，不重复添加
					return false
				}
			}
		}

		// 构造 admin.MustRegister("模块名", new(模块名.Mounter))
		newCall := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent(adminAlias),
				Sel: ast.NewIdent("MustRegister"),
			},
			Args: []ast.Expr{
				moduleLit,
				&ast.CallExpr{
					Fun: ast.NewIdent("new"),
					Args: []ast.Expr{&ast.SelectorExpr{
						X:   ast.NewIdent(moduleAlias),
						Sel: ast.NewIdent("Mounter"),
					}},
				},
			},
		}
		newStmt := &ast.ExprStmt{X: newCall}

		fn.Body.List = append(fn.Body.List, newStmt)

		return false
	})
	if !foundInit {
		return "", "", fmt.Errorf("模块入口缺少 init 函数，无法自动注册")
	}
	// 将修改后的AST写回文件
	var buf bytes.Buffer
	if err = format.Node(&buf, fileSet, f); err != nil {
		return "", "", err
	}
	return buf.String(), filePath, nil
}

// 复用已有别名，补回删除最后一个模块时移除的导入。
func ensureModuleImport(file *ast.File, path, preferred string) (string, error) {
	used := make(map[string]bool)
	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return "", err
		}
		alias := filepath.Base(importPath)
		if importPath == "github.com/Gary-Yez/go-admin" {
			alias = "admin"
		}
		if spec.Name != nil {
			alias = spec.Name.Name
		}
		if importPath == path {
			if alias == "." || alias == "_" {
				return "", fmt.Errorf("导入 %s 使用特殊别名，请先手动调整", path)
			}
			return alias, nil
		}
		used[alias] = true
	}
	alias := preferred
	for suffix := 2; used[alias]; suffix++ {
		alias = fmt.Sprintf("%s%d", preferred, suffix)
	}
	spec := &ast.ImportSpec{Name: ast.NewIdent(alias), Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)}}
	for _, decl := range file.Decls {
		if group, ok := decl.(*ast.GenDecl); ok && group.Tok == token.IMPORT {
			group.Specs = append(group.Specs, spec)
			file.Imports = append(file.Imports, spec)
			return alias, nil
		}
	}
	file.Decls = append([]ast.Decl{&ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{spec}}}, file.Decls...)
	file.Imports = append(file.Imports, spec)
	return alias, nil
}

func readModulePath(goModPath string) (string, error) {
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("module directive not found in %s", goModPath)
}
