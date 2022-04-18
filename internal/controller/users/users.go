package users

import (
	"errors"
	"github.com/wujunyi792/crispy-waffle-be/internal/db"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var dbManage *UserDBManage = nil

func init() {
	logger.Info.Println("[ USER ]start init Table ...")
	dbManage = GetManage()
}

type UserDBManage struct {
	mDB     *db.MainGORM
	sDBLock sync.RWMutex
}

func (m *UserDBManage) getGOrmDB() *gorm.DB {
	return m.mDB.GetDB()
}

func (m *UserDBManage) atomicDBOperation(op func()) {
	m.sDBLock.Lock()
	op()
	m.sDBLock.Unlock()
}

func GetManage() *UserDBManage {
	if dbManage == nil {
		var userDb = db.MustCreateGorm()
		err := userDb.GetDB().AutoMigrate(&Mysql.Permission{}, &Mysql.User{})
		if err != nil {
			logger.Error.Fatalln(err)
			return nil
		}
		dbManage = &UserDBManage{mDB: userDb}
	}
	return dbManage
}

func SetLoginLog(id string, token string) {
	GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{"last_login_token": token, "last_login_time": time.Now()})
}

func CheckPhoneExist(phone string) bool {
	res := false
	GetManage().atomicDBOperation(func() {
		res = GetManage().getGOrmDB().Model(&Mysql.User{}).Where("phone = ?", phone).Find(&Mysql.User{}).RowsAffected > 0
	})
	return res
}

func CheckUserNameExist(username string) bool {
	res := false
	GetManage().atomicDBOperation(func() {
		res = GetManage().getGOrmDB().Model(&Mysql.User{}).Where("nick_name = ?", username).Find(&Mysql.User{}).RowsAffected > 0
	})
	return res
}

func GetEntity(entity *Mysql.User) *Mysql.User {
	GetManage().atomicDBOperation(func() {
		GetManage().getGOrmDB().Where(entity).Find(entity)
	})
	return entity
}

func GetUserByID(id string) (entity *Mysql.User) {
	entity = &Mysql.User{}
	GetManage().atomicDBOperation(func() {
		GetManage().getGOrmDB().Where("id = ?", entity.ID).Find(entity)
	})
	return entity
}

func RegisterUser(user *Mysql.User) (err error) {
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Create(user).Error
	})
	return
}

func UpdateAvatar(id string, avatar string) (err error) {
	//tx := GetManage().getGOrmDB().Begin()
	//defer tx.Commit()
	//
	//entity := Mysql.User{
	//	Base: Mysql.Base{
	//		ID: id,
	//	},
	//}
	//
	//if tx.Where("id = ?", entity.ID).Find(&entity).RowsAffected == 0 {
	//	tx.Rollback()
	//	return errors.New("用户不存在")
	//}
	//
	//entity.Avatar = avatar
	//
	//err = tx.Model(&entity).Update("avatar", avatar).Error
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("avatar", avatar)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdatePhone(id string, phone string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("phone", phone)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdatePasswordAndSalt(phone string, passwordHashed string, salt string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("phone = ?", phone).Updates(map[string]interface{}{"password": passwordHashed, "salt": salt})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdateUserName(id string, username string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("nick_name", username)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdateSex(id string, sex int) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("sex", sex)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdateSignature(id string, signature string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("signature", signature)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdateEmail(id string, email string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("email", email)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdateStatus(id string, status string) error {
	res := GetManage().getGOrmDB().Model(&Mysql.User{}).Where("id = ?", id).Update("status", status)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func PermissionAdd(id string, permissionString string) (err error) {
	var u Mysql.User
	u.ID = id
	_ = GetEntity(&u)
	if u.ID == "" {
		return errors.New("用户不存在")
	}
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Model(&Mysql.User{
			Base: Mysql.Base{
				ID: u.ID,
			},
		}).
			Association("Permission").Append(&Mysql.Permission{
			PermissionName: permissionString,
		})
	})
	return err
}

func PermissionDel(id string, permissionString string) (err error) {
	var u Mysql.User
	u.ID = id
	_ = GetEntity(&u)
	if u.ID == "" {
		return errors.New("用户不存在")
	}
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Model(&Mysql.User{
			Base: Mysql.Base{
				ID: u.ID,
			},
		}).
			Association("Permission").Delete(&Mysql.Permission{
			PermissionName: permissionString,
		})
	})
	return err
}

func PermissionClear(id string) (err error) {
	var u Mysql.User
	u.ID = id
	_ = GetEntity(&u)
	if u.ID == "" {
		return errors.New("用户不存在")
	}
	GetManage().atomicDBOperation(func() {
		err = GetManage().getGOrmDB().Model(&Mysql.User{
			Base: Mysql.Base{
				ID: u.ID,
			},
		}).
			Association("Permission").Clear()
	})
	return err
}

func PermissionCheck(id string, permission string) (exist bool) {
	var permissionEntity Mysql.Permission
	err := GetManage().getGOrmDB().Model(&Mysql.User{
		Base: Mysql.Base{
			ID: id,
		},
	}).Where("permission_name = ?", permission).Association("Permission").Find(&permissionEntity)
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	if permissionEntity.PermissionName != "" {
		return true
	}
	return
}
