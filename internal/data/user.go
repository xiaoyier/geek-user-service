package data

import (
	"geek-user-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
)

func init() {
	RegisterAutoMigrates(&UserInfo{})
}

// UserInfo User PO
type UserInfo struct {
	ID         uint      `gorm:"primary_key;column:id"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime(6)"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime(6)"`
	UserID     int64     `gorm:"unique_index;column:user_id;type:bigint(20)"`
	UserName   string    `gorm:"unique_index;column:user_name;type:varchar(50)"`
	UserPasswd string    `gorm:"column:user_passwd;type:varchar(50)"`
	Email      string    `gorm:"column:email;type:varchar(30)"`
	Avatar     string    `gorm:"column:avatar;type:varchar(100)"`
}

func (u *UserInfo) TableName() string {
	return "user"
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data.user")),
	}
}

func (r *userRepo) GetUserByUserName(userName string) (*biz.User, error) {
	user := new(UserInfo)
	d := r.data.db.Where("user_name = ?", userName).First(&user)
	if d.RecordNotFound() {
		return nil, errors.Wrapf(biz.ErrRecordNotFound, "GetUserByUserName: user: %s", userName)
	}

	if d.Error != nil {
		return nil, errors.Wrapf(d.Error, "GetUserByUserName: user: %s", userName)
	}
	// PO -> DO
	return &biz.User{
		UserID: user.UserID,
		UserName: user.UserName,
		Password: user.UserPasswd,
	}, nil
}

func (r *userRepo) GetUserByUserID(userId int64) (*biz.User, error) {
	user := new(UserInfo)
	d := r.data.db.Where("user_id = ?", userId).First(&user)
	if d.RecordNotFound() {
		return nil, errors.Wrapf(biz.ErrRecordNotFound, "GetUserByUserName: user: %d", userId)
	}

	if d.Error != nil {
		return nil, errors.Wrapf(d.Error, "GetUserByUserName: user: %f", userId)
	}
	// PO -> DO
	return &biz.User{
		UserID: user.UserID,
		UserName: user.UserName,
		Password: user.UserPasswd,
	}, nil
}

func (r *userRepo) CreateUser(user *biz.User) error {
	// DO -> PO
	userInfo := UserInfo{
		UserID: user.UserID,
		UserName: user.UserName,
		UserPasswd: user.Password,
	}
	d := r.data.db.Save(&userInfo)
	if d.Error != nil {
		return errors.Wrapf(d.Error, "CreateUser: user: %d_%s", user.UserID, user.UserName)
	}
	return nil
}

func (r *userRepo) UpdateUser(user *biz.User) error {
	// DO -> PO
	userInfo := UserInfo{
		UserID: user.UserID,
		UserName: user.UserName,
		UserPasswd: user.Password,
	}
	d := r.data.db.Model(userInfo).Where(UserInfo{UserID: userInfo.UserID}).Update(&userInfo)
	if d.Error != nil {
		return errors.Wrapf(d.Error, "UpdateUser: user: %d", userInfo.UserID)
	}
	return nil
}


func (r *userRepo) GetUserByUserNameFromRedis(userName string) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) GetUserByUserIDFromRedis(userId int64) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) CreateUserFromRedis(user *biz.User) error {
	return nil
}

func (r *userRepo) UpdateUserFromRedis(user *biz.User) error {
	return nil
}

