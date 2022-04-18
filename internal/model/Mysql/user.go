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

	RealName  string  `json:"realName"`
	NickName  string  `json:"nickName" gorm:"unique"`
	Sex       int     `json:"sex"`
	Phone     string  `json:"phone" gorm:"unique"`
	Email     string  `json:"email"`
	Signature string  `json:"signature"`
	Status    string  `json:"status"`
	Avatar    string  `json:"avatar"`
	Oauth     []Oauth `json:"-"`

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

// https://learnku.com/articles/20074
type Oauth struct {
	Base
	UserID     string `gorm:"size:191"`
	OauthType  string `gorm:"comment:第三方登陆类型 weibo、qq、wechat 等"`
	OauthId    string `gorm:"unique;第三方 uid openid 等"`
	UnionId    string `gorm:"comment:QQ / 微信同一主体下 Unionid 相同"`
	Credential string `gorm:"commrnt:密码凭证 /access_token (目前更多是存储在缓存里)"`
}

func (u *Oauth) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4().String()
	return
}
