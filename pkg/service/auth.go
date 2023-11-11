package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
	"unicode/utf8"
)

const (
	salt       = "djlspocnacdgd"
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

func (a *AuthService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type")
	}
	return claims.UserId, nil
}

func (a *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	passwordHash := hex.EncodeToString(hash.Sum([]byte(salt)))
	if !utf8.ValidString(passwordHash) {
		log.Printf("password valid: %s      ", passwordHash)
	}
	return fmt.Sprintf("%s", passwordHash)
}

func (a *AuthService) FillHtml() ([]music.Genre, error) {
	//tmpl, err := template.New("recommend").ParseFiles("./public/recommendation.html")
	//if err != nil {
	//	return err
	//}
	genres, err := a.repo.GetAllGenres()
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func (a *AuthService) InsertRecommendation(genres []music.Genre, userId int64) error {
	if err := a.repo.InsertUserGenre(genres, userId); err != nil {
		return err
	}
	return nil
}
