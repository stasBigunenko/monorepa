package jwt

import (
	"time"

	"github.com/stasBigunenko/monorepa/model"
)

type UserClaims model.JWTUserClaims

func (cl *UserClaims) addExpTime(duration int) {
	cl.ExpiresAt = time.Now().Add(time.Duration(duration) * time.Minute).Unix()
}
