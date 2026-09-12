package sys_login_log

import "time"

const (
	Success = "success"
	Failed  = "failed"
	Blocked = "blocked"
)

type SysLoginLog struct {
	Id        uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"index;comment:登录时间"`
	UserId    uint      `json:"user_id" gorm:"index;comment:成功登录的管理员ID"`
	Username  string    `json:"username" gorm:"size:191;index;comment:登录账号"`
	IP        string    `json:"ip" gorm:"size:45;index;comment:连接来源IP"`
	Status    string    `json:"status" gorm:"size:16;index;comment:登录结果"`
	Message   string    `json:"message" gorm:"size:512;comment:结果说明"`
	UserAgent string    `json:"user_agent" gorm:"size:512;comment:客户端信息"`
	Duration  int64     `json:"duration" gorm:"comment:登录耗时毫秒"`
}

type CleanupBody struct {
	Days int `json:"days" binding:"required,min=1,max=36500"`
}
