package utils_test

import (
	"testing"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var cfg = config.JWTConfig{
	PrivateKey: "dummy-key",
	Lifetime:   1 * time.Second,
	Issuer:     "test",
}

var invalidConfig = config.JWTConfig{
	PrivateKey: "dummy-key",
	Lifetime:   1 * time.Second,
	Issuer:     "invalid-issuer",
}

func generateImpersonatorToken(user entity.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss": 12345,
		"sub": user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(cfg.Lifetime)),
		"iat": jwt.NewNumericDate(time.Now()),
	})

	signedToken, _ := token.SignedString([]byte(cfg.PrivateKey))

	return signedToken
}
func TestGenerateAndVerifyToken(t *testing.T) {
	validToken1 := utils.GenerateToken(entity.User{ID: 11111, Username: "user1"}, cfg)
	validToken2 := utils.GenerateToken(entity.User{ID: 22222, Username: "user2"}, cfg)
	validToken3 := utils.GenerateToken(entity.User{ID: 33333, Username: "user3"}, invalidConfig)
	impersonatorToken := generateImpersonatorToken(entity.User{ID: 44444, Username: "user4"})

	testCases := []struct {
		title   string
		token   string
		want    utils.Jwt
		success bool
		delay   time.Duration
	}{
		{"Verify Valid Token 1", validToken1, utils.Jwt{Subject: "11111"}, true, 0},
		{"Verify Valid Token 2", validToken2, utils.Jwt{Subject: "22222"}, true, 0},
		{"Verify Invalid Token", "invalidTokenValue", utils.Jwt{}, false, 0},
		{"Verify Expired Token", validToken1, utils.Jwt{Subject: "11111"}, false, 1 * time.Second},
		{"Verify Invalid Issuer", validToken3, utils.Jwt{}, false, 0},
		{"Verify Impersonator Token", impersonatorToken, utils.Jwt{Subject: "44444"}, false, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			time.Sleep(tc.delay)
			verified, err := utils.VerifyToken(tc.token, cfg)

			if tc.success {
				assert.Greater(t, verified.Exp, time.Now().Unix())
				assert.Equal(t, tc.want.Subject, verified.Subject)
				assert.Equal(t, cfg.Issuer, verified.Issuer)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
