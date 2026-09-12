package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gary-Yez/go-admin/internal/bizconfig"
	"time"
	"unicode/utf8"
)

const JWTExpireMinutesKey = "jwt.expire_minutes"
const PasswordMaxBytes = 72

// decodeLoginNumber 只解析一次，同时检查整数类型和业务范围。
func decodeLoginNumber(key string, value json.RawMessage) (int64, error) {
	var number int64
	err := json.Unmarshal(value, &number)
	if key == JWTExpireMinutesKey {
		if err != nil || number <= 0 || number > int64((1<<63-1)/time.Minute) {
			return 0, errors.New("登录有效期必须为正整数，且不能超出时间范围")
		}
		return number, nil
	}
	if err != nil || number <= 0 || int64(int(number)) != number {
		return 0, errors.New("登录配置必须为正整数")
	}
	if key == "login.password_min_length" && number > PasswordMaxBytes {
		return 0, errors.New("密码最小长度不能超过 72")
	}
	if (key == "login.failure_window_minutes" || key == "login.lock_minutes") && number > int64((1<<63-1)/time.Minute) {
		return 0, errors.New("登录配置超出时间范围")
	}
	return number, nil
}

func ValidateLoginValue(key string, value json.RawMessage) error {
	switch key {
	case JWTExpireMinutesKey, "login.max_failures", "login.failure_window_minutes", "login.lock_minutes", "login.password_min_length":
		_, err := decodeLoginNumber(key, value)
		return err
	}
	return nil
}

func ReadJWTLifetime() (time.Duration, error) {
	value, err := bizconfig.ReadValue(JWTExpireMinutesKey)
	if err != nil {
		return 0, err
	}
	minutes, err := decodeLoginNumber(JWTExpireMinutesKey, value)
	if err != nil {
		return 0, err
	}
	return time.Duration(minutes) * time.Minute, nil
}

func ReadLoginInt(key string) (int, error) {
	value, err := bizconfig.ReadValue(key)
	if err != nil {
		return 0, err
	}
	number, err := decodeLoginNumber(key, value)
	return int(number), err
}

func ValidatePassword(password string) error {
	minimum, err := ReadLoginInt("login.password_min_length")
	if err != nil {
		return err
	}
	if !utf8.ValidString(password) || utf8.RuneCountInString(password) < minimum {
		return fmt.Errorf("密码至少需要 %d 个字符", minimum)
	}
	if len(password) > PasswordMaxBytes {
		return errors.New("密码不能超过 72 字节")
	}
	return nil
}
