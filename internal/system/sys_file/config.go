package sys_file

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/Gary-Yez/go-admin/internal/bizconfig"
)

const (
	thresholdKey         = "storage.multipart_threshold_mib"
	partSizeKey          = "storage.part_size_mib"
	sessionDaysKey       = "storage.session_days"
	linkExpireKey        = "storage.link_expire_seconds"
	allowedExtensionsKey = "storage.allowed_extensions"
	maxUploadParts       = 10000
)

type uploadPolicy struct {
	AllowedExtensions []string `json:"allowed_extensions"`
	OrdinaryLimit     int64    `json:"ordinary_limit"`
	PartSize          int64    `json:"part_size"`
	MaxSize           int64    `json:"max_size"`
	SessionDays       int64    `json:"session_days"`
}

func decodeUploadValue(key string, raw json.RawMessage) (int64, error) {
	var size int64
	minimum, maximum, multiplier := int64(1), int64(5120), int64(1<<20)
	label, unit := "分片上传阈值", "MiB"
	switch key {
	case thresholdKey:
		minimum = 100
	case partSizeKey:
		minimum, maximum, label = 20, 1024, "分片大小"
	case linkExpireKey:
		minimum, maximum, multiplier, label, unit = 60, 86400, 1, "访问链接有效期", "秒"
	case sessionDaysKey:
		maximum, multiplier, label, unit = 365, 1, "上传会话有效期", "天"
	}
	if err := json.Unmarshal(raw, &size); err != nil || size < minimum || size > maximum {
		return 0, fmt.Errorf("%s必须为 %d～%d %s的整数", label, minimum, maximum, unit)
	}
	return size * multiplier, nil
}

// ValidateConfigValue 在配置保存时复用上传模块的业务范围检查。
func ValidateConfigValue(key string, value json.RawMessage) error {
	if key == allowedExtensionsKey {
		_, err := decodeAllowedExtensions(value)
		return err
	}
	if key != thresholdKey && key != partSizeKey && key != sessionDaysKey && key != linkExpireKey {
		return nil
	}
	_, err := decodeUploadValue(key, value)
	return err
}

func readUploadPolicy() (uploadPolicy, error) {
	var policy uploadPolicy
	for _, item := range []struct {
		key    string
		target *int64
	}{
		{thresholdKey, &policy.OrdinaryLimit},
		{partSizeKey, &policy.PartSize},
		{sessionDaysKey, &policy.SessionDays},
	} {
		raw, err := bizconfig.ReadValue(item.key)
		if err != nil {
			return policy, err
		}
		size, err := decodeUploadValue(item.key, raw)
		if err != nil {
			return policy, err
		}
		*item.target = size
	}
	extensions, err := readAllowedExtensions()
	if err != nil {
		return policy, err
	}
	policy.AllowedExtensions = extensions
	policy.MaxSize = policy.PartSize * maxUploadParts
	return policy, nil
}

func readLinkTTL() (time.Duration, error) {
	raw, err := bizconfig.ReadValue(linkExpireKey)
	if err != nil {
		return 0, err
	}
	seconds, err := decodeUploadValue(linkExpireKey, raw)
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

// 配置只接受单个扩展名，不接受 MIME、通配符或复合后缀。
var extensionPattern = regexp.MustCompile("^[a-z0-9]+$")

func decodeAllowedExtensions(raw json.RawMessage) ([]string, error) {
	var value []string
	if err := json.Unmarshal(raw, &value); err != nil || value == nil {
		return nil, fmt.Errorf("允许上传的文件类型必须为字符串列表")
	}
	extensions := make([]string, 0)
	seen := make(map[string]bool)
	for _, item := range value {
		extension := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(item)), ".")
		if !extensionPattern.MatchString(extension) {
			return nil, fmt.Errorf("文件类型 %q 无效，每项请填写 jpg、png、pdf 这样的单个扩展名", item)
		}
		if !seen[extension] {
			extensions = append(extensions, extension)
			seen[extension] = true
		}
	}
	return extensions, nil
}

func readAllowedExtensions() ([]string, error) {
	raw, err := bizconfig.ReadValue(allowedExtensionsKey)
	if err != nil {
		return nil, err
	}
	return decodeAllowedExtensions(raw)
}

func validateExtension(name string, extensions []string) error {
	if len(extensions) == 0 {
		return nil
	}
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(name)), ".")
	for _, allowed := range extensions {
		if extension == allowed {
			return nil
		}
	}
	return fmt.Errorf("不允许上传此文件类型，仅支持：%s", strings.Join(extensions, "、"))
}

func validateCurrentExtension(name string) error {
	extensions, err := readAllowedExtensions()
	if err != nil {
		return err
	}
	return validateExtension(name, extensions)
}
