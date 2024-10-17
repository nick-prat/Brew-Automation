package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "admin"
	password = "1234"
	dbname   = "admin"
)

type key int

const VERSION_KEY key = 1
const SESSION_KEY key = 2
const USER_KEY key = 3

type MessageBody struct {
	Field string `json:"message"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

type PKResponse struct {
	PK int `json:"pk"`
}

type HTTPError struct {
	StatusCode int
	Err        error
}

type LoginResponse struct {
	Token string `json:"token"`
}

func InternalServerError(err error) *HTTPError {
	return &HTTPError{StatusCode: http.StatusInternalServerError, Err: err}
}

func BadRequestError(err error) *HTTPError {
	return &HTTPError{StatusCode: http.StatusBadRequest, Err: err}
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

type RequestEnvironment struct {
	db *sql.DB
}

func generateError(error string) string {
	val, err := json.Marshal(ErrorBody{Error: error})
	if err != nil {
		return panicResponse()
	}

	return string(val)
}

func responseHandler(f func(w http.ResponseWriter, r *http.Request) (string, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := f(w, r)
		fmt.Printf("%s", err)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, generateError(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, resp)
		}
	}
}

func pkResponse(pk int) (string, error) {
	resp := PKResponse{PK: pk}
	val, err := json.Marshal(&resp)
	if err != nil {
		return "", err
	} else {
		return string(val), nil
	}
}

func panicResponse() string {
	return "{\"error\": \"Unknown Error!\"}"
}

func (env *RequestEnvironment) login(_ http.ResponseWriter, r *http.Request) (string, error) {
	userDAO := UserDAO{db: env.db}
	sessionDAO := SessionDAO{db: env.db}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", BadRequestError(err)
	}

	userLogin := UserLogin{}
	err = json.Unmarshal(body, &userLogin)
	if err != nil {
		return "", BadRequestError(err)
	}

	user, err := userDAO.Login(userLogin.Email, userLogin.Password)
	if err != nil {
		return "", BadRequestError(err)
	}

	session, err := sessionDAO.Create(user.UserID)
	if err != nil {
		return "", InternalServerError(err)
	}

	loginResponse := LoginResponse{Token: session}

	resp, err := json.Marshal(loginResponse)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) register(_ http.ResponseWriter, r *http.Request) (string, error) {
	userDAO := UserDAO{db: env.db}
	sessionDAO := SessionDAO{db: env.db}

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

	session, err := sessionDAO.Create(userId)
	if err != nil {
		return "", InternalServerError(err)
	}

	loginResponse := LoginResponse{Token: session}

	resp, err := json.Marshal(loginResponse)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) getTempLog(_ http.ResponseWriter, r *http.Request) (string, error) {
	tempLogDAO := TempLogDAO{db: env.db}
	pk, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return "", BadRequestError(err)
	}

	tempLog, err := tempLogDAO.Get(pk)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := json.Marshal(tempLog)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) getTempLogs(_ http.ResponseWriter, _ *http.Request) (string, error) {
	tempLogDAO := TempLogDAO{db: env.db}

	tempLogs, err := tempLogDAO.Select(10)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := json.Marshal(tempLogs)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) postTempLog(_ http.ResponseWriter, r *http.Request) (string, error) {
	tempLogDAO := TempLogDAO{db: env.db}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", BadRequestError(err)
	}

	tempLog := TempLog{}
	err = json.Unmarshal(body, &tempLog)
	if err != nil {
		return "", BadRequestError(err)
	}

	pk, err := tempLogDAO.Create(&tempLog)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := pkResponse(pk)
	if err != nil {
		return "", InternalServerError(err)
	}

	return resp, nil
}

func initDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	return db
}

func (env *RequestEnvironment) versionMiddleWare(r *http.Request) (*http.Request, error) {
	return r.WithContext(context.WithValue(r.Context(), VERSION_KEY, 0.1)), nil
}

func (env *RequestEnvironment) sessionMiddleWare(r *http.Request) (*http.Request, error) {
	authorization := r.Header.Get("Authorization")

	if len(authorization) == 0 {
		return r, nil
	}

	session, err := NewSessionDAO(env.db).Get(authorization)
	if err != nil {
		return nil, err
	}

	fmt.Println("Foobar!")

	return r.WithContext(context.WithValue(r.Context(), SESSION_KEY, session)), nil
}

func (env *RequestEnvironment) userMiddleWare(r *http.Request) (*http.Request, error) {
	userDAO := NewUserDAO(env.db)
	ctxSession := r.Context().Value(SESSION_KEY)
	if ctxSession == nil {
		return r, nil
	}

	session := ctxSession.(*Session)

	user, err := userDAO.GetByPK(session.UserId)
	if err != nil {
		return nil, err
	}

	return r.WithContext(context.WithValue(r.Context(), USER_KEY, user)), nil
}

type Middleware = func(r *http.Request) (*http.Request, error)

func main() {
	db := initDB()
	defer db.Close()

	env := RequestEnvironment{db: db}

	middlewares := []Middleware{
		env.versionMiddleWare,
		env.sessionMiddleWare,
		env.userMiddleWare,
	}

	middleware := func(f func(w http.ResponseWriter, r *http.Request) (string, error)) func(w http.ResponseWriter, r *http.Request) (string, error) {
		return func(w http.ResponseWriter, r *http.Request) (string, error) {
			var req = r
			for _, m := range middlewares {
				fmt.Println("step")
				tempreq, err := m(req)
				if err != nil {
					return "", nil
				}
				req = tempreq
			}
			str, err := f(w, req)
			return str, err
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /temp-log", responseHandler(middleware(env.postTempLog)))
	mux.HandleFunc("GET /temp-log", responseHandler(middleware(env.getTempLogs)))
	mux.HandleFunc("GET /temp-log/{id}", responseHandler(middleware(env.getTempLog)))
	mux.HandleFunc("POST /login", responseHandler(middleware(env.login)))
	mux.HandleFunc("POST /register", responseHandler(middleware(env.register)))

	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("Error starting server %s\n", err)
			os.Exit(1)
		}

		cancelCtx()
	}()

	<-ctx.Done()
}
