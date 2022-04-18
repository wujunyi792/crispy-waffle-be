package Mysql

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Base
	LastLoginToken string    `json:"-"`
	LastLoginTime  time.Time `json:"-"`
	Password       string    `json:"-"`
	Salt           string    `json:"-"`

	RealName    string `json:"realName"`
	NickName    string `json:"nickName" gorm:"unique"`
	Sex         int    `json:"sex"`
	Phone       string `json:"phone" gorm:"unique"`
	Email       string `json:"email"`
	Signature   string `json:"signature"`
	Status      string `json:"status"`
	Avatar      string `json:"avatar"`
	GithubID    int64  `json:"githubID" gorm:"unique"`
	GitHubToken string `json:"-"`
	GithubInfo  string `json:"-"`

	Permission []*Permission `gorm:"many2many:users_permission;" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4().String()
	u.LastLoginTime = time.Now()
	return
}

type Permission struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PermissionName string  `gorm:"primarykey"`
	Remark         string  `json:"remark"`
	ApplyUser      []*User `gorm:"many2many:users_permission;" json:"applyUser"`
}
