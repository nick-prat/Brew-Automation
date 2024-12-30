package api

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"raspberrysour/dao"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type UserLogin struct {
	Email    string
	Password string
}

type UserClaims struct {
	Id   int32  `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type UserRegister = UserLogin

func generateUserJWT(userId int32, key *rsa.PrivateKey) (string, error) {
	claims := UserClaims{
		userId,
		"user",
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func parseUserJWT(tokenString string, key *rsa.PublicKey) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT")
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid JWT Claims")
	}
}

func (env *RequestEnvironment) Register(_ http.ResponseWriter, r *http.Request) (string, error) {
	userDAO := dao.NewUserDAO(env.db)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", BadRequestError(err)
	}

	userRegister := UserRegister{}
	err = json.Unmarshal(body, &userRegister)
	if err != nil {
		return "", BadRequestError(err)
	}

	userId, err := userDAO.Register(userRegister.Email, userRegister.Password)
	if err != nil {
		return "", BadRequestError(err)
	}

	jwt, err := generateUserJWT(userId, env.jwtSigningKey)
	if err != nil {
		return "", InternalServerError(err)
	}
	loginResponse := LoginResponse{jwt}

	resp, err := json.Marshal(loginResponse)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) Login(_ http.ResponseWriter, r *http.Request) (string, error) {
	userDAO := dao.NewUserDAO(env.db)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", BadRequestError(err)
	}

	userLogin := UserLogin{}
	err = json.Unmarshal(body, &userLogin)
	if err != nil {
		return "", BadRequestError(err)
	}

	userId, err := userDAO.Login(userLogin.Email, userLogin.Password)
	if err != nil {
		return "", BadRequestError(err)
	}

	jwt, err := generateUserJWT(userId, env.jwtSigningKey)
	if err != nil {
		return "", InternalServerError(err)
	}
	response := LoginResponse{jwt}

	resp, err := json.Marshal(response)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) GetUsers(_ http.ResponseWriter, r *http.Request) (string, error) {
	fmt.Println("foobar")
	if !IsAuthenticated(r) {
		return "", UnauthorizedError()
	}
	userDAO := dao.NewUserDAO(env.db)

	users, err := userDAO.Select(0, 10)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := json.Marshal(users)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}
