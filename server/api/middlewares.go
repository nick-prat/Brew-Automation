package api

import (
	"context"
	"net/http"
	"raspberrysour/dao"
)

type key int

const VERSION_KEY key = 1

// const SESSION_KEY key = 2
const USER_KEY key = 3

func (env *RequestEnvironment) VersionMiddleWare(r *http.Request) (*http.Request, error) {
	return r.WithContext(context.WithValue(r.Context(), VERSION_KEY, 0.1)), nil
}

func (env *RequestEnvironment) UserMiddleWare(r *http.Request) (*http.Request, error) {
	userDAO := dao.NewUserDAO(env.db)

	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		return r, nil
	}

	userClaims, err := parseUserJWT(authorization, env.jwtVerifyKey)
	if err != nil {
		return nil, err
	}

	user, err := userDAO.GetByPK(userClaims.Id)
	if err != nil {
		return nil, err
	}

	return r.WithContext(context.WithValue(r.Context(), USER_KEY, user)), nil
}

type Middleware = func(r *http.Request) (*http.Request, error)
