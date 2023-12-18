package music

import (
	"encoding/base64"
	"strconv"
	"time"
)

type User struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name" binding:"required"`
	UserName     string    `json:"username" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	Picture      string    `json:"picture"`
	CreateDate   time.Time `json:"createDate"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

type UserJSON struct {
	Id           string `json:"id"`
	Name         string `json:"name" binding:"required"`
	UserName     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Picture      string `json:"picture"`
	CreateDate   string `json:"create_date"`
	ModifiedDate string `json:"modified_date"`
}

func (u *User) SetFromJSON(user UserJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(user.Id)
	if err != nil {
		return err
	}
	u.Id, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(user.Name)
	if err != nil {
		return err
	}
	u.Name = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(user.UserName)
	if err != nil {
		return err
	}
	u.UserName = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(user.Email)
	if err != nil {
		return err
	}
	u.Email = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return err
	}
	u.Password = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(user.Picture)
	if err != nil {
		return err
	}
	u.Picture = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(user.CreateDate)
	if err != nil {
		return err
	}
	u.CreateDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(user.ModifiedDate)
	if err != nil {
		return err
	}
	u.ModifiedDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}
	return err
}
