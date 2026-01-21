package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
	Verify(map[string]string) (bool, error)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
	login, err := s.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, err
}

func (s DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return false, err
	} else {
		if jwtToken.Valid {
			mapClaims := jwtToken.Claims.(jwt.MapClaims)
			if claims, err := domain.BuildClaimsFromJwtMapClaims(mapClaims); err != nil {
				return false, nil
			} else {

				if claims.IsUserRole() {
					if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
						return false, nil
					}
				}

				isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
				return isAuthorized, nil

			}
		} else {
			return false, errors.New("Invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewLoginService(repo domain.AuthRepository, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo, permissions}
}
