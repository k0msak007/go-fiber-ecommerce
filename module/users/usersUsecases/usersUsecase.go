package usersUsecases

import (
	"fmt"

	"github.com/k0msak007/go-fiber-ecommerce/config"
	"github.com/k0msak007/go-fiber-ecommerce/module/users"
	"github.com/k0msak007/go-fiber-ecommerce/module/users/usersRepositories"
	"github.com/k0msak007/go-fiber-ecommerce/pkg/kawaiiauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassPort, error)
	GetPassPort(req *users.UserCredential) (*users.UserPassPort, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassPort, error)
	DeleteOauth(oauthId string) error
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassPort, error)
	GetUserProfile(userId string) (*users.User, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassPort, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassPort, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetPassPort(req *users.UserCredential) (*users.UserPassPort, error) {
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	accessToken, _ := kawaiiauth.NewKawaiiAuth(kawaiiauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, _ := kawaiiauth.NewKawaiiAuth(kawaiiauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	passport := &users.UserPassPort{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassPort, error) {
	// Parse Token
	claims, err := kawaiiauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := kawaiiauth.NewKawaiiAuth(
		kawaiiauth.Access,
		u.cfg.Jwt(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}
	refreshToken, _ := kawaiiauth.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassPort{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	return u.usersRepository.DeleteOauth(oauthId)
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	return u.usersRepository.GetProfile(userId)
}
