package sys_role

import (
	"gitee.com/mxcker/go-admin/server/global"
	"strconv"
)

func InitData() error {
	db := global.DB.Model(&SysRole{})
	count := int64(0)
	if err := db.Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		data := SysRole{
			Name:    "超级管理员",
			Default: true,
		}
		if err := db.Create(&data).Error; err != nil {
			return err
		}
		_, err := Enforcer.AddPolicy(strconv.Itoa(int(data.Id)), "*", "*")
		if err != nil {
			return err
		}
	}
	return nil
}
