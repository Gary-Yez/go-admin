package sys_autocode

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
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
			Path: filepath.Join(ServerPath, "./core/templates/server", name+".go.tmpl"),
		})
	}
	// Web端模板
	WebTemplatesPath = append(WebTemplatesPath, TemplateItem{
		Name: "api",
		Path: filepath.Join(ServerPath, "./core/templates/web/api.ts.tmpl"),
	})
	WebTemplatesPath = append(WebTemplatesPath, TemplateItem{
		Name: "view",
		Path: filepath.Join(ServerPath, "./core/templates/web/view.vue.tmpl"),
	})
}

func createAutoMigrateCall(args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("global"),
				Sel: ast.NewIdent("DB"),
			},
			Sel: ast.NewIdent("AutoMigrate"),
		},
		Args: args,
	}
}
func isGlobalDBAutoMigrate(expr ast.Expr) bool {
	// 验证调用链结构
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "AutoMigrate" {
		return false
	}

	if subSel, ok := sel.X.(*ast.SelectorExpr); ok {
		return subSel.X.(*ast.Ident).Name == "global" &&
			subSel.Sel.Name == "DB"
	}
	return false
}

func formatAutoMigrateCalls(code []byte) []byte {
	// 匹配AutoMigrate调用模式
	re := regexp.MustCompile(`(global\.DB\.AutoMigrate\()([\s\S]*?)(\))`)

	return re.ReplaceAllFunc(code, func(match []byte) []byte {
		parts := re.FindSubmatch(match)
		opening := parts[1]
		content := parts[2]
		closing := parts[3]
		// 拆分参数并添加换行
		params := strings.Split(strings.TrimSuffix(strings.TrimSpace(string(content)), ","), ", ")
		//formatted := "(\n\t" + strings.Join(params, ",\n\t") + ",\n)"
		indent := "\t\t" // 二级缩进
		formatted := "\n" + indent + strings.Join(params, ",\n"+indent) + ",\n\t"
		return []byte(string(opening) + formatted + string(closing))
	})
}

func getTemplateContent(templatePath string, data *GenerateBody) (string, error) {
	// 打开模板文件
	var buffer bytes.Buffer
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buffer, *data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func getRouterContent(moduleName string) (string, string, error) {
	// 配置文件和位置
	fileSet := token.NewFileSet()
	filePath := filepath.Join(ServerPath, "./modules/enter.go")
	targetImportPath := "gitee.com/mxcker/go-admin/server/modules/" + moduleName
	// 解析源代码文件
	f, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", "", err
	}
	// 1. 先处理导入
	importAdded := false
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		// 检查是否已经存在该导入
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			if importSpec.Path.Value == strconv.Quote(targetImportPath) {
				importAdded = true
				break
			}
		}

		// 添加新导入
		if !importAdded {
			newImport := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(targetImportPath),
				},
			}
			genDecl.Specs = append(genDecl.Specs, newImport)
			importAdded = true
		}
	}

	// 如果没有找到import声明，创建新的
	if !importAdded {
		newImportDecl := &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(targetImportPath),
					},
				},
			},
		}

		// 将新的import声明插入到文件最前面
		newDecls := make([]ast.Decl, 0, len(f.Decls)+1)
		newDecls = append(newDecls, newImportDecl)
		f.Decls = append(newDecls, f.Decls...)
	}
	// 查找目标函数
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Register" {
			return true
		}
		// 找到函数体中的return语句
		for i, stmt := range fn.Body.List {
			ret, ok := stmt.(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				continue
			}
			// 构造要插入的表达式：otherModule.Register("/other", adminAuthGroup, publicGroup)
			newCall := &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(moduleName),
					Sel: ast.NewIdent("Register"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.STRING, Value: `"` + moduleName + `"`},
					ast.NewIdent("adminAuthGroup"),
					ast.NewIdent("publicGroup"),
				},
			}
			// 在return前插入新语句
			newStmt := &ast.ExprStmt{X: newCall}
			fn.Body.List = append(
				fn.Body.List[:i],
				append([]ast.Stmt{newStmt}, fn.Body.List[i:]...)...,
			)
			break
		}
		return false
	})
	// 将修改后的AST写回文件
	var buf bytes.Buffer
	if err = printer.Fprint(&buf, fileSet, f); err != nil {
		return "", "", err
	}
	return buf.String(), filePath, nil
}

func getInitializationContent(moduleName string, modelName string) (string, string, error) {
	fileSet := token.NewFileSet()
	initializationPath := filepath.Join(ServerPath, "./modules/init.go")
	targetImportPath := "gitee.com/mxcker/go-admin/server/modules/" + moduleName
	// 解析文件
	f, err := parser.ParseFile(fileSet, initializationPath, nil, parser.ParseComments)
	if err != nil {
		return "", "", err
	}
	// 1. 处理导入声明
	importAdded := false
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		// 检查是否已存在该导入
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			if importSpec.Path.Value == strconv.Quote(targetImportPath) {
				importAdded = true
				break
			}
		}
		// 添加新导入
		if !importAdded {
			newImport := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(targetImportPath),
				},
			}
			genDecl.Specs = append(genDecl.Specs, newImport)
			importAdded = true
		}
	}
	// 如果没有找到import声明，创建新的
	if !importAdded {
		newImportDecl := &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(targetImportPath),
					},
				},
			},
		}
		f.Decls = append([]ast.Decl{newImportDecl}, f.Decls...)
	}
	// 2. 修改函数
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Init" {
			return true
		}

		// 查找现有的AutoMigrate调用
		var targetCall *ast.CallExpr
		for _, stmt := range fn.Body.List {
			if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
				if call, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
					if isGlobalDBAutoMigrate(call.Fun) {
						targetCall = call
						break
					}
				}
			}
		}

		// 创建新参数
		newArg := &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent(moduleName),
					Sel: ast.NewIdent(modelName),
				},
			},
		}

		if targetCall != nil {
			// 追加到现有参数列表
			targetCall.Args = append(targetCall.Args, newArg)
		} else {
			// 创建新调用语句
			newCall := createAutoMigrateCall(newArg)
			newAssign := &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("err")},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{newCall},
			}
			// 插入到return语句前
			for i, stmt := range fn.Body.List {
				if _, ok := stmt.(*ast.ReturnStmt); ok {
					fn.Body.List = append(
						fn.Body.List[:i],
						append([]ast.Stmt{newAssign}, fn.Body.List[i:]...)...,
					)
					break
				}
			}
		}

		return false
	})
	// 写回文件
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fileSet, f); err != nil {
		return "", "", err
	}
	// 使用正则表达式格式化AutoMigrate调用
	formatted := formatAutoMigrateCalls(buf.Bytes())
	return string(formatted), initializationPath, nil
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
