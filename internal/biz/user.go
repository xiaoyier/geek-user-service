package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	sysErrors "errors"
	"geek-user-service/internal/pkg/jwt"
	snowflake "geek-user-service/internal/pkg/snowflake"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

// User User DO
type User struct {
	UserID int64
	UserName string
	Password string
	LoginToken string
	RefreshToken string
}

// UserRepo repo
type UserRepo interface {
	// mysql
	GetUserByUserName(userName string) (*User, error)
	GetUserByUserID(userId int64) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error

	// redis
	GetUserByUserNameFromRedis(userName string) (*User, error)
	GetUserByUserIDFromRedis(userId int64) (*User, error)
	CreateUserFromRedis(user *User) error
	UpdateUserFromRedis(user *User) error
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log: log.NewHelper(log.With(logger, "module", "user.biz"))}
}

func (u *UserUsecase) SignUp(ctx context.Context, user *User) error {
	// TODO: redis
	result, err := u.repo.GetUserByUserName(user.UserName)
	if err != nil && !sysErrors.Is(err, ErrRecordNotFound) {
		return err
	}

	if result != nil {
		return errors.Wrapf(ErrRecordAlreadyExist, "SignUp: user: %s", user.UserName)
	}

	userId := snowflake.GenSnowflakeID()
	if userId == 0 {
		return errors.Wrapf(ErrUserIDGen, "SignUp: user: %s", user.UserName)
	}

	user.UserID = int64(userId)
	user.Password = MD5Encrypt(user.Password)
	err = u.repo.CreateUser(user)
	return err
}

func (u *UserUsecase) Login(ctx context.Context, user *User) (result *User, err error) {

	// TODO: redis
	if user.UserID != 0 {
		result, err = u.repo.GetUserByUserID(user.UserID)
	} else if len(user.UserName) > 0 {
		result, err = u.repo.GetUserByUserName(user.UserName)
	}

	if err != nil {
		return
	}

	if result == nil {
		err = errors.Wrapf(ErrRecordNotFound, "Login: user: %d", user.UserID)
		return
	}

	newPassword := MD5Encrypt(user.Password)
	if result.Password != newPassword {
		err = errors.Wrapf(ErrPasswordWrong, "Login: user_id: %d, password: %s", user.UserID, newPassword)
		return
	}

	login, refresh, err := jwt.GenToken(user.UserID)
	if err != nil {
		return
	}

	result.LoginToken = login
	result.RefreshToken = refresh
	return
}

func (u *UserUsecase) Update(ctx context.Context, user *User, newName string) (result *User, err error) {
	// TODO: redis
	result, err = u.repo.GetUserByUserName(user.UserName)
	if err != nil {
		return
	}

	if result == nil {
		err = errors.Wrapf(ErrRecordNotFound, "Login: user: %d", user.UserID)
		return
	}

	if result.Password != MD5Encrypt(user.Password) {
		err = errors.Wrapf(ErrPasswordWrong, "Login: user_id: %d, password: %s", user.UserID, user.Password)
		return
	}

	result.UserName = newName
	err = u.repo.UpdateUser(result)
	if err != nil {
		return
	}

	login, refresh, err := jwt.GenToken(user.UserID)
	if err != nil {
		return
	}

	result.LoginToken = login
	result.RefreshToken = refresh
	return
}

func MD5Encrypt(src string) string {

	sha := md5.New()
	sha.Write([]byte(src))
	return hex.EncodeToString(sha.Sum(nil))
}
