package sys_file

import "time"

const (
	uploading = "uploading"
	ready     = "ready"
)

type SysFile struct {
	Id           uint       `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time  `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	Name         string     `json:"name" gorm:"size:255;comment:文件名"`
	Size         int64      `json:"size" gorm:"comment:文件大小"`
	ContentType  string     `json:"content_type" gorm:"size:255;comment:文件类型"`
	Engine       string     `json:"engine" gorm:"size:32;index;comment:存储引擎"`
	Status       string     `json:"status" gorm:"size:16;index;comment:上传状态"`
	UserId       uint       `json:"user_id" gorm:"index;comment:上传者"`
	Username     string     `json:"username" gorm:"size:191;comment:上传账号"`
	Multipart    bool       `json:"multipart" gorm:"comment:分片上传"`
	ExpiresAt    *time.Time `json:"expires_at" gorm:"index;comment:会话过期时间"`
	LastModified int64      `json:"last_modified" gorm:"comment:原文件修改时间"`
	ObjectKey    string     `json:"-" gorm:"size:191;uniqueIndex"`
	StorageId    uint       `json:"storage_id" gorm:"index;comment:存储账号"`
	StorageName  string     `json:"storage_name" gorm:"-"`
	UploadJSON   string     `json:"-" gorm:"type:text;comment:存储上传会话"`
	IsOwner      bool       `json:"is_owner" gorm:"-"`
}

type SysFilePart struct {
	FileId uint   `json:"-" gorm:"primaryKey;autoIncrement:false"`
	Number int    `json:"number" gorm:"primaryKey;autoIncrement:false"`
	Size   int64  `json:"size"`
	ETag   string `json:"-" gorm:"size:512"`
}

type BeginBody struct {
	StorageId    uint   `json:"storage_id"`
	Name         string `json:"name" binding:"required,max=255"`
	Size         int64  `json:"size" binding:"required,gt=0"`
	LastModified int64  `json:"last_modified" binding:"gte=0"`
}

type IdBody struct {
	Id uint `json:"id" form:"id" binding:"required,gt=0"`
}

type Session struct {
	File     *SysFile      `json:"file"`
	Parts    []SysFilePart `json:"parts"`
	PartSize int64         `json:"part_size"`
}
