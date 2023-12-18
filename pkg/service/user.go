package service

import (
	"errors"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepo
}

func NewUserServiceImpl(repo repository.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (u *UserServiceImpl) ChangeUsername(userID int64, username string) error {
	return u.repo.UpdateUsername(userID, username)
}

func (u *UserServiceImpl) ChangePassword(userId int64, newPassword, oldPassword string) error {
	currentPassword, err := u.repo.SelectPassword(userId)
	if err != nil {
		return err
	}
	if currentPassword != oldPassword {
		return errors.New("wrong old password")
	}
	return u.repo.UpdatePassword(userId, newPassword)
}

func (u *UserServiceImpl) ChangePicture(userID int64, picture string) error {
	return u.repo.UpdatePicture(userID, picture)
}

func (u *UserServiceImpl) ChangeEmail(userID int64, email string) error {
	return u.repo.UpdateEmail(userID, email)
}

func (u *UserServiceImpl) GetAllUsers() ([]music.User, error) {
	return u.repo.SelectAllUsers()
}

func (u *UserServiceImpl) DeleteUser(idUser int64, reason string) error {
	return u.repo.DeleteUser(idUser, reason)
}

func (u *UserServiceImpl) GetUserByID(idUser int64) (music.User, error) {
	return u.repo.SelectUserById(idUser)
}
