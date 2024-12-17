package utils_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestCSRFToken_IsExpired(t *testing.T) {
	csrf := utils.CSRFToken{
		Token:     "dummy-token",
		ExpiredAt: time.Now().Add(-10 * time.Minute).Unix(),
	}
	assert.True(t, csrf.IsExpired())
}

func TestEncryptAndDecryptSessionValue(t *testing.T) {
	cfg := config.SessionConfig{
		SecretKey:       "dummy-secret-key",
		CSRFTokenLength: 10,
	}

	session := utils.Session{
		CSRFToken: utils.GenerateCSRFToken(cfg),
		IssuedAt:  time.Now().Unix(),
	}

	assert.Equal(t, cfg.CSRFTokenLength, len(session.CSRFToken.Token))

	encrypted := utils.EncryptSessionValue(session, cfg)
	assert.NotNil(t, encrypted)

	_, err := utils.DecryptSessionValue("short-invalid-value", cfg)
	assert.NotNil(t, err)

	decrypted, err := utils.DecryptSessionValue(encrypted, cfg)
	assert.Equal(t, session.CSRFToken, decrypted.CSRFToken)
	assert.Equal(t, session.IssuedAt, decrypted.IssuedAt)
	assert.Nil(t, err)
}

func TestGetAndSetSessionCookie(t *testing.T) {
	cfg := config.SessionConfig{
		SecretKey:       "dummy-secret-key",
		CSRFTokenLength: 10,
	}

	session := utils.Session{
		CSRFToken: utils.GenerateCSRFToken(cfg),
		IssuedAt:  time.Now().Unix(),
	}

	w := httptest.NewRecorder()
	utils.SetSessionCookie(w, session, cfg)

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, "session", cookie.Name)

	decrypted, err := utils.DecryptSessionValue(cookie.Value, cfg)
	assert.Nil(t, err)
	assert.Equal(t, session.CSRFToken, decrypted.CSRFToken)
	assert.Equal(t, session.IssuedAt, decrypted.IssuedAt)
	assert.Equal(t, true, cookie.HttpOnly)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)

	session = utils.GetSessionCookie(req, cfg)
	assert.Equal(t, session.CSRFToken, decrypted.CSRFToken)
	assert.Equal(t, session.IssuedAt, decrypted.IssuedAt)
}
