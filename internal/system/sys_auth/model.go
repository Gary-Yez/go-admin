package sys_auth

type loginJson struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangeInfoBody 使用指针区分未提交字段，只更新本次提交的资料。
type ChangeInfoBody struct {
	Nickname     *string `json:"nickname"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email"`
	AvatarFileId *uint   `json:"avatar_file_id"`
}
