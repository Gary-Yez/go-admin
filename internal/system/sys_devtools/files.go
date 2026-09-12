package sys_devtools

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func contentHash(content []byte) string { return fmt.Sprintf("%x", sha256.Sum256(content)) }

// Restrict outputs to the generator's directories and reject symlink traversal.
func checkedPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	allowed := false
	for _, root := range []string{filepath.Join(ServerPath, "modules"), filepath.Join(WebPath, "apis"), filepath.Join(WebPath, "views")} {
		rootAbs, err := filepath.Abs(root)
		if err != nil {
			return "", err
		}
		rel, err := filepath.Rel(rootAbs, abs)
		if err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel) {
			allowed = true
		}
	}
	if !allowed {
		return "", fmt.Errorf("生成路径超出允许目录：%s", abs)
	}
	for p := abs; ; p = filepath.Dir(p) {
		info, err := os.Lstat(p)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("生成路径不能经过符号链接：%s", p)
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return abs, nil
}

func inspectFiles(items []WriteItem) ([]WriteItem, error) {
	seen := make(map[string]bool)
	for i := range items {
		item := &items[i]
		path, err := checkedPath(item.Path)
		if err != nil {
			return nil, err
		}
		item.Path = path
		key := strings.ToLower(path)
		if seen[key] {
			return nil, fmt.Errorf("生成目标重复：%s", path)
		}
		seen[key] = true
		content, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			item.Action = "create"
			continue
		}
		if err != nil {
			return nil, err
		}
		item.ExistingContent = string(content)
		item.ExistingHash = contentHash(content)
		if item.Content == string(content) {
			item.Action = "unchanged"
		} else if item.registration {
			item.Action = "modify"
		} else {
			item.Action = "conflict"
			item.RequiresOverwrite = true
		}
	}
	return items, nil
}

type stagedFile struct {
	item    WriteItem
	stage   string
	backup  string
	applied bool
}

func checkSnapshot(item WriteItem) error {
	if _, err := checkedPath(item.Path); err != nil {
		return err
	}
	content, err := os.ReadFile(item.Path)
	if item.Action == "create" || item.Action == "missing" {
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("目标文件已出现，请重新预览：%s", item.Path)
	}
	if err != nil {
		return err
	}
	if contentHash(content) != item.ExistingHash {
		return fmt.Errorf("文件已变化，请重新预览：%s", item.Path)
	}
	return nil
}

func stageContent(dir, pattern string, content []byte, mode os.FileMode) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	name := f.Name()
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(content)
	}
	if err == nil {
		err = f.Sync()
	}
	err = errors.Join(err, f.Close())
	if err != nil {
		return "", errors.Join(err, os.Remove(name))
	}
	return name, nil
}

// Stage every output and backup before replacing anything. A failed history save
// is part of the same operation and restores the files as well.
func commitFiles(items []WriteItem, approvals map[string]string, saveHistory func() error) (result error) {
	for _, item := range items {
		if item.RequiresOverwrite && approvals[item.Path] != item.ExistingHash {
			return fmt.Errorf("文件已存在，默认禁止覆盖：%s；请预览后在 overwrite_files 中提交该路径及 existing_hash", item.Path)
		}
		if err := checkSnapshot(item); err != nil {
			return err
		}
	}
	var staged []*stagedFile
	var createdDirs []string
	defer func() {
		failed := result != nil
		for i := len(staged) - 1; i >= 0; i-- {
			f := staged[i]
			if failed && f.applied {
				var err error
				if f.backup != "" {
					err = os.Rename(f.backup, f.item.Path)
					if err == nil {
						f.backup = ""
					}
				} else {
					err = os.Remove(f.item.Path)
				}
				if err != nil {
					result = errors.Join(result, fmt.Errorf("恢复 %s 失败（备份保留在 %s）：%w", f.item.Path, f.backup, err))
					// Keep the backup for manual recovery.
					f.backup = ""
				}
			}
			for _, temp := range []string{f.stage, f.backup} {
				if temp != "" {
					if err := os.Remove(temp); err != nil && !os.IsNotExist(err) {
						result = errors.Join(result, fmt.Errorf("清理临时文件 %s 失败：%w", temp, err))
					}
				}
			}
		}
		if failed {
			for i := len(createdDirs) - 1; i >= 0; i-- {
				if err := os.Remove(createdDirs[i]); err != nil && !os.IsNotExist(err) {
					result = errors.Join(result, fmt.Errorf("清理新建目录 %s 失败：%w", createdDirs[i], err))
				}
			}
		}
	}()
	for _, item := range items {
		if item.Action == "unchanged" || item.Action == "missing" {
			continue
		}
		dir := filepath.Dir(item.Path)
		var missing []string
		for p := dir; ; p = filepath.Dir(p) {
			_, err := os.Stat(p)
			if err == nil {
				break
			}
			if !os.IsNotExist(err) {
				return err
			}
			missing = append(missing, p)
			if filepath.Dir(p) == p {
				return fmt.Errorf("找不到输出目录根路径")
			}
		}
		for i := len(missing) - 1; i >= 0; i-- {
			if err := os.Mkdir(missing[i], 0755); err != nil {
				return err
			}
			createdDirs = append(createdDirs, missing[i])
		}
		f := &stagedFile{item: item}
		staged = append(staged, f)
		mode := os.FileMode(0644)
		if item.Action != "create" {
			info, err := os.Stat(item.Path)
			if err != nil {
				return err
			}
			mode = info.Mode().Perm()
			f.backup, err = stageContent(dir, ".go-admin-backup-*", []byte(item.ExistingContent), mode)
			if err != nil {
				return err
			}
		}
		var err error
		if item.Action != "delete" {
			f.stage, err = stageContent(dir, ".go-admin-stage-*", []byte(item.Content), mode)
			if err != nil {
				return err
			}
		}
	}
	// Recheck all snapshots after staging, before applying the first change.
	for _, item := range items {
		if err := checkSnapshot(item); err != nil {
			return err
		}
	}
	for _, f := range staged {
		if err := checkSnapshot(f.item); err != nil {
			return err
		}
		if f.item.Action == "delete" {
			if err := os.Remove(f.item.Path); err != nil {
				return err
			}
			f.applied = true
			continue
		}
		if f.item.Action == "create" {
			// Reserve a new path exclusively so a concurrently created file is not replaced.
			file, err := os.OpenFile(f.item.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			f.applied = true
			if err := file.Close(); err != nil {
				return err
			}
		}
		if err := os.Rename(f.stage, f.item.Path); err != nil {
			return err
		}
		f.stage = ""
		f.applied = true
	}
	return saveHistory()
}
