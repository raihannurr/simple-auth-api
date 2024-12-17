package utils

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Session struct {
	IssuedAt  int64     `json:"issued_at"`
	CSRFToken CSRFToken `json:"csrf_token"`
}

type CSRFToken struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

func (c CSRFToken) IsExpired() bool {
	return c.ExpiredAt < time.Now().Unix()
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateCSRFToken(cfg config.SessionConfig) CSRFToken {
	return CSRFToken{
		Token:     GenerateRandomString(cfg.CSRFTokenLength),
		ExpiredAt: time.Now().Add(cfg.CSRFTokenExp).Unix(),
	}
}

func DecryptSessionValue(encrypted string, cfg config.SessionConfig) (Session, error) {
	decoded, _ := base64.StdEncoding.DecodeString(encrypted)

	block, err := aes.NewCipher([]byte(cfg.SecretKey))
	PanicIfError(err)

	gcm, err := cipher.NewGCM(block)
	PanicIfError(err)

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return Session{}, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	PanicIfError(err)

	session := Session{}
	err = json.Unmarshal(decrypted, &session)
	PanicIfError(err)

	return session, nil
}

func EncryptSessionValue(session Session, cfg config.SessionConfig) string {
	json, err := json.Marshal(session)
	PanicIfError(err)

	block, err := aes.NewCipher([]byte(cfg.SecretKey))
	PanicIfError(err)

	gcm, err := cipher.NewGCM(block)
	PanicIfError(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = crand.Read(nonce)
	PanicIfError(err)

	encrypted := gcm.Seal(nonce, nonce, json, nil)

	return base64.StdEncoding.EncodeToString(encrypted)
}

func SetSessionCookie(w http.ResponseWriter, session Session, cfg config.SessionConfig) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    EncryptSessionValue(session, cfg),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func GetSessionCookie(r *http.Request, cfg config.SessionConfig) Session {
	cookie, _ := r.Cookie("session")
	session, _ := DecryptSessionValue(cookie.Value, cfg)
	return session
}
