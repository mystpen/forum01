package user

import (
	"forum/internal/types"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo types.UserRepo // userDB
}

func NewUserService(repo types.UserRepo) *UserService {
	return &UserService{repo}
}

func (u *UserService) CreateUser(userData *types.CreateUserData) error {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &types.User{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: string(hashedPW),
	}

	u.repo.CreateUserDB(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) CheckUserExists(user *types.CreateUserData) (bool, types.ErrText) {
	errs := types.ErrText{}

	NameExists := u.repo.GetUserNameDB(user.Username)
	EmailExists := u.repo.GetUserEmailDB(user.Email)

	existBool := false
	if NameExists == nil {
		errs.Username = "Username already exists"
		existBool = true
	}
	if EmailExists == nil {
		errs.Email = "Email already exists"
		existBool = true
	}
	if !existBool {
		err := u.CreateUser(user)
		if err != nil {
			log.Fatalln(err)
		}

	}
	return existBool, errs
}

func (u *UserService) CheckLogin(user *types.GetUserData) (int, error) {
	return u.repo.CheckLoginDB(user)
}

func (u *UserService) AddToken(userid int, cookie string) error {
	return u.repo.AddTokenDB(userid, cookie)
}

func (u *UserService) RemoveToken(token string) error {
	return u.repo.RemoveTokenDB(token)
}

func (u *UserService) GetUserByToken(token string) (user *types.User, err error){
	return u.repo.GetUserByToken(token)
}