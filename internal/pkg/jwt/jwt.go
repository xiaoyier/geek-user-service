package jwt

import (
	"github.com/pkg/errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	UserID int64
}

const ExpireTimeLoginToken = time.Hour * 24 * 7
const ExpireTimeRefreshToken = time.Hour

var snowflakeSecret = []byte("golang 大法好!")

func GenToken(userID int64) (loginToken, refreshToken string, err error) {

	loginToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ExpireTimeLoginToken).Unix(),
			Issuer:    "clover",
		},
		userID,
	}).SignedString(snowflakeSecret)

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ExpireTimeRefreshToken).Unix(),
		Issuer:    "clover",
	}).SignedString(snowflakeSecret)

	if err != nil {
		err = errors.Wrap(err, "GenToken error")
	}
	return
}

func ParseUserID(loginToken string) (userId int64, err error) {

	var token *jwt.Token
	claims := &Claims{}
	token, err = jwt.ParseWithClaims(loginToken, claims, func(token *jwt.Token) (interface{}, error) {
		return snowflakeSecret, nil
	})

	if err != nil {
		err = errors.Wrapf(err, "ParseUserID: login_token: %s", loginToken)
		return
	}

	userId = claims.UserID
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

func CheckValid(refreshToken string) bool {
	token, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return snowflakeSecret, nil
	})

	return token != nil && token.Valid
}
