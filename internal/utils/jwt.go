package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/errors"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
	IssuedAt int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

func ParseJwtClaims(claims jwt.Claims) (Jwt, error) {
	j := Jwt{}
	iss, errIss := claims.GetIssuer()
	subject, errSubject := claims.GetSubject()
	exp, errExp := claims.GetExpirationTime()
	iat, errIat := claims.GetIssuedAt()

	err := errors.Join(errIss, errSubject, errExp, errIat)
	if err == nil {
		j = Jwt{
			Issuer:   iss,
			Subject:  subject,
			Exp:      exp.Unix(),
			IssuedAt: iat.Unix(),
		}
	}

	return j, err
}

func GenerateToken(user entity.User, cfg config.JWTConfig) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Lifetime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	signedToken, err := token.SignedString([]byte(cfg.PrivateKey))
	PanicIfError(err)

	return signedToken
}

func VerifyToken(token string, cfg config.JWTConfig) (Jwt, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// expiration check also handled here
		return []byte(cfg.PrivateKey), nil
	})
	if err != nil {
		return Jwt{}, err
	}

	t, err := ParseJwtClaims(parsedToken.Claims)
	if err != nil {
		return Jwt{}, err
	}

	if t.Issuer != cfg.Issuer {
		return Jwt{}, errors.ErrInvalidToken
	}

	return t, nil
}
