package api

import (
	"context"
	"fmt"
	"net/http"
	"raspberrysour/dao"
)

type key int

const VERSION_KEY key = 1
const USER_KEY key = 2

func (env *RequestEnvironment) VersionMiddleWare(r *http.Request) (*http.Request, error) {
	return r.WithContext(context.WithValue(r.Context(), VERSION_KEY, 0.1)), nil
}

func (env *RequestEnvironment) UserMiddleWare(r *http.Request) (*http.Request, error) {
	userDAO := dao.NewUserDAO(env.db)

	authorization := r.Header.Get("Authorization")
	fmt.Printf("Authorization %s\n", authorization)
	if len(authorization) == 0 {
		return r, nil
	}

	fmt.Println("1")
	userClaims, err := parseUserJWT(authorization, env.jwtVerifyKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("2")
	user, err := userDAO.Get(userClaims.Id)
	if err != nil {
		return nil, err
	}

	fmt.Println(user)
	return r.WithContext(context.WithValue(r.Context(), USER_KEY, user)), nil
}

type Middleware = func(r *http.Request) (*http.Request, error)
