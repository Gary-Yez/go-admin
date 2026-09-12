package sys_devtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/request"
)

func (s *serviceStruct) PreviewDeleteHistory(ids []uint) (*DeleteHistoryPlan, error) {
	generationMu.Lock()
	defer generationMu.Unlock()
	return buildDeletePlan(ids)
}

func (s *serviceStruct) DeleteHistory(req *DeleteHistoryBody) error {
	generationMu.Lock()
	defer generationMu.Unlock()
	if !req.DeleteFiles {
		return s.DeleteByIds(&request.ReqIds{Ids: req.Ids})
	}
	plan, err := buildDeletePlan(req.Ids)
	if err != nil {
		return err
	}
	if req.PreviewToken == "" || req.PreviewToken != plan.Token {
		return errors.New("记录或本地文件已变化，请重新查看删除清单并确认")
	}
	if err := commitFiles(plan.Files, nil, func() error { return s.DeleteByIds(&request.ReqIds{Ids: req.Ids}) }); err != nil {
		return err
	}
	// Backups are gone only after commitFiles returns; failed operations must keep
	// their directories available for rollback.
	if err := cleanupEmptyModuleDirs(plan.Files); err != nil {
		return fmt.Errorf("文件和生成记录已删除，但部分空目录清理失败：%w", err)
	}
	return nil
}

func cleanupEmptyModuleDirs(files []WriteItem) error {
	roots := []string{filepath.Join(ServerPath, "modules"), filepath.Join(WebPath, "views")}
	for i := range roots {
		root, err := filepath.Abs(roots[i])
		if err != nil {
			return err
		}
		roots[i] = root
	}
	seen := make(map[string]bool)
	var result error
	for _, file := range files {
		if file.Action != "delete" && file.Action != "missing" {
			continue
		}
		dir := filepath.Dir(file.Path)
		// Only a module's direct directory is eligible, never shared roots/apis.
		eligible := false
		for _, root := range roots {
			rel, err := filepath.Rel(root, dir)
			if err == nil && rel != "." && rel != ".." && !filepath.IsAbs(rel) && !strings.ContainsAny(rel, `/\`) {
				eligible = true
			}
		}
		key := strings.ToLower(dir)
		if !eligible || seen[key] {
			continue
		}
		seen[key] = true
		if _, err := checkedPath(dir); err != nil {
			result = errors.Join(result, err)
			continue
		}
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			result = errors.Join(result, fmt.Errorf("%s：%w", dir, err))
			continue
		}
		if len(entries) != 0 {
			continue
		}
		// Non-recursive removal cannot delete files added concurrently.
		if err := os.Remove(dir); err != nil && !os.IsNotExist(err) {
			if current, readErr := os.ReadDir(dir); readErr == nil && len(current) > 0 {
				continue
			}
			result = errors.Join(result, fmt.Errorf("%s：%w", dir, err))
		}
	}
	return result
}

func buildDeletePlan(ids []uint) (*DeleteHistoryPlan, error) {
	unique := make(map[uint]bool)
	for _, id := range ids {
		if id == 0 {
			return nil, errors.New("记录编号不正确")
		}
		unique[id] = true
	}
	if len(unique) == 0 || len(unique) > 100 {
		return nil, errors.New("请选择1至100条生成记录")
	}
	var histories []SysAutoCode
	if err := state.DB().Where("id IN ?", ids).Order("id ASC").Find(&histories).Error; err != nil {
		return nil, err
	}
	if len(histories) != len(unique) {
		return nil, errors.New("部分生成记录已被删除，请刷新列表")
	}
	modules := make(map[string]bool)
	var targets []WriteItem
	for _, history := range histories {
		// Only validated module names can determine local paths; never accept paths from the client.
		if err := validateGenerateBody(&GenerateBody{ModuleName: history.ModuleName, ModelName: history.ModelName, ChineseModuleName: "删除"}); err != nil {
			return nil, fmt.Errorf("模块 %s 的记录无法安全定位文件：%w", history.ModuleName, err)
		}
		modules[history.ModuleName] = true
		paths := []string{
			filepath.Join(ServerPath, "modules", history.ModuleName, "model.go"),
			filepath.Join(ServerPath, "modules", history.ModuleName, "controller.go"),
			filepath.Join(ServerPath, "modules", history.ModuleName, "service.go"),
			filepath.Join(ServerPath, "modules", history.ModuleName, "enter.go"),
			filepath.Join(WebPath, "apis", history.ModuleName+".ts"),
			filepath.Join(WebPath, "views", history.ModuleName, "index.vue"),
		}
		for _, path := range paths {
			targets = append(targets, WriteItem{Path: path})
		}
	}
	files, err := inspectFiles(targets)
	if err != nil {
		return nil, err
	}
	// Missing files stay in the snapshot to detect files appearing after confirmation.
	for i := range files {
		files[i].RequiresOverwrite = false
		if files[i].Action != "create" {
			files[i].Action = "delete"
		} else {
			files[i].Action = "missing"
		}
	}
	entryPath, err := checkedPath(filepath.Join(ServerPath, "modules", "enter.go"))
	if err != nil {
		return nil, err
	}
	source, err := os.ReadFile(entryPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	entry := WriteItem{Path: entryPath, Action: "missing"}
	if err == nil {
		content, err := removeModuleRegistrations(source, modules)
		if err != nil {
			return nil, err
		}
		entry = WriteItem{Path: entryPath, Content: string(content), ExistingContent: string(source), ExistingHash: contentHash(source), Action: "modify"}
		if bytes.Equal(source, content) {
			entry.Action = "unchanged"
		}
	}
	files = append(files, entry)
	snapshot, err := json.Marshal(struct {
		Histories []SysAutoCode
		Files     []WriteItem
	}{histories, files})
	if err != nil {
		return nil, err
	}
	return &DeleteHistoryPlan{Files: files, Token: contentHash(snapshot)}, nil
}

func removeModuleRegistrations(source []byte, modules map[string]bool) ([]byte, error) {
	host, err := readModulePath(filepath.Join(ServerPath, "go.mod"))
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "enter.go", source, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("模块入口无法解析：%w", err)
	}
	imports := make(map[string]string)
	adminAlias := "admin"
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		if path == "github.com/Gary-Yez/go-admin" && spec.Name != nil {
			adminAlias = spec.Name.Name
		}
		for module := range modules {
			if path == strings.TrimSuffix(host, "/")+"/modules/"+module {
				alias := module
				if spec.Name != nil {
					alias = spec.Name.Name
				}
				if alias == "." || alias == "_" {
					return nil, fmt.Errorf("模块 %s 使用特殊导入，请先手动移除注册", module)
				}
				imports[module] = alias
			}
		}
	}
	changed := false
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "init" || fn.Recv != nil || fn.Body == nil {
			continue
		}
		kept := make([]ast.Stmt, 0, len(fn.Body.List))
		for _, stmt := range fn.Body.List {
			remove := false
			if expr, ok := stmt.(*ast.ExprStmt); ok {
				if call, ok := expr.X.(*ast.CallExpr); ok && len(call.Args) == 2 {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "MustRegister" {
						owner, ok := sel.X.(*ast.Ident)
						name, nameOK := call.Args[0].(*ast.BasicLit)
						if ok && owner.Name == adminAlias && nameOK && name.Kind == token.STRING {
							module, _ := strconv.Unquote(name.Value)
							if modules[module] {
								ctor, ok := call.Args[1].(*ast.CallExpr)
								if !ok || len(ctor.Args) != 1 {
									return nil, fmt.Errorf("模块 %s 的注册代码经过定制，请先手动移除", module)
								}
								newIdent, ok := ctor.Fun.(*ast.Ident)
								if !ok || newIdent.Name != "new" {
									return nil, fmt.Errorf("模块 %s 的注册代码经过定制", module)
								}
								typ, ok := ctor.Args[0].(*ast.SelectorExpr)
								if !ok || typ.Sel.Name != "Mounter" {
									return nil, fmt.Errorf("模块 %s 的注册类型经过定制", module)
								}
								pkg, ok := typ.X.(*ast.Ident)
								if !ok || pkg.Name != imports[module] {
									return nil, fmt.Errorf("模块 %s 的注册导入不匹配", module)
								}
								remove = true
								changed = true
							}
						}
					}
				}
			}
			if !remove {
				kept = append(kept, stmt)
			}
		}
		fn.Body.List = kept
	}
	used := make(map[string]bool)
	// Refuse to remove imports still referenced by custom code.
	for _, decl := range file.Decls {
		if g, ok := decl.(*ast.GenDecl); ok && g.Tok == token.IMPORT {
			continue
		}
		ast.Inspect(decl, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				used[id.Name] = true
			}
			return true
		})
	}
	removeImports := make(map[string]bool)
	for module, alias := range imports {
		if used[alias] {
			return nil, fmt.Errorf("模块入口仍有对 %s 的自定义引用，请先手动处理", module)
		}
		removeImports[strings.TrimSuffix(host, "/")+"/modules/"+module] = true
	}
	if changed && !used[adminAlias] && adminAlias != "_" && adminAlias != "." {
		removeImports["github.com/Gary-Yez/go-admin"] = true
	}
	decls := make([]ast.Decl, 0, len(file.Decls))
	for _, decl := range file.Decls {
		if g, ok := decl.(*ast.GenDecl); ok && g.Tok == token.IMPORT {
			specs := make([]ast.Spec, 0, len(g.Specs))
			for _, spec := range g.Specs {
				path, _ := strconv.Unquote(spec.(*ast.ImportSpec).Path.Value)
				if removeImports[path] {
					changed = true
				} else {
					specs = append(specs, spec)
				}
			}
			g.Specs = specs
			if len(specs) == 0 {
				continue
			}
		}
		decls = append(decls, decl)
	}
	if !changed {
		return source, nil
	}
	file.Decls = decls
	var result bytes.Buffer
	if err := format.Node(&result, fset, file); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
