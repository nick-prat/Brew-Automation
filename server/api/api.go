package api

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type MessageBody struct {
	Field string `json:"message"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

type PKBody struct {
	PK int32 `json:"pk"`
}

type HTTPError struct {
	StatusCode int
	Err        error
}

type RequestEnvironment struct {
	db            *sqlx.DB
	jwtSigningKey *rsa.PrivateKey
	jwtVerifyKey  *rsa.PublicKey
}

func NewRequestEnvironment(db *sqlx.DB, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RequestEnvironment {
	return &RequestEnvironment{db, privateKey, publicKey}
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

func InternalServerError(err error) *HTTPError {
	return &HTTPError{StatusCode: http.StatusInternalServerError, Err: err}
}

func BadRequestError(err error) *HTTPError {
	return &HTTPError{StatusCode: http.StatusBadRequest, Err: err}
}

func UnauthorizedError() *HTTPError {
	return &HTTPError{StatusCode: http.StatusUnauthorized, Err: fmt.Errorf("Unauthorized")}
}

func PanicResponse() string {
	return "{\"error\": \"Unknown Error!\"}"
}

func GenerateError(error string) string {
	val, err := json.Marshal(ErrorBody{Error: error})
	if err != nil {
		return PanicResponse()
	}

	return string(val)
}

func ResponseHandler(f func(w http.ResponseWriter, r *http.Request) (string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := f(w, r)
		if err != nil {
			httpErr, ok := err.(*HTTPError)
			if ok {
				w.WriteHeader(httpErr.StatusCode)
				io.WriteString(w, GenerateError(err.Error()))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, GenerateError(err.Error()))
			}
		} else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, resp)
		}
	}
}

func PKResponse(pk int32) (string, error) {
	resp := PKBody{PK: pk}
	val, err := json.Marshal(&resp)
	if err != nil {
		return "", err
	} else {
		return string(val), nil
	}
}

func IsAuthenticated(r *http.Request) bool {
	return r.Context().Value(USER_KEY) != nil
	// return r.Context().Value(USER_KEY) != nil
}
