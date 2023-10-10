package service

import (
	"crypto/sha1"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"time"
	"unicode/utf8"
)

const (
	salt       = "11"
	signingKey = "dvkdvkdfvsfvhg"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(user music.User) (int64, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := a.repo.GetUser(username, a.generatePasswordHash(password))
	if err != nil {
		return "", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (a *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	passwordHash, _ := utf8.DecodeRune(hash.Sum([]byte(salt)))
	return fmt.Sprintf("%s", passwordHash)
}
