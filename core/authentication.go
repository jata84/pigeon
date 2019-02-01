package core

import (
	"net/http"
	"net/url"
	"strings"
)

func authenticate() bool {
	return true
}

type Credentials struct {
}

type Authentication struct {
	credentials    interface{}
	authentication func() bool
}

type RemoteBasicAuthenticationCredentials struct {
	credentials Credentials
	remote_url  url.URL
	username    string
	password    string
}

func NewRemoteBasicAuthentication(credentials RemoteBasicAuthenticationCredentials) Authentication {
	return Authentication{
		credentials: credentials,
		authentication: func() bool {
			resp, err := http.PostForm(credentials.remote_url.String(),
				url.Values{"username": {credentials.username}, "password": {credentials.password}})
			return (err == nil) && (resp.Status == "200")

		},
	}
}

type TokenAuthentication struct {
	Authentication *Authentication
}

type JWTAuthentication struct {
	Authentication *Authentication
}

type IAuthenticationMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

type TokenAuthenticationMiddleware struct {
	IAuthenticationMiddleware
}

func (t *TokenAuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		authorized_token := Configuration.Authorization.Authorization_token
		if strings.HasPrefix(authorization, "BASIC") {
			authorization = strings.TrimSpace(strings.TrimPrefix(authorization, "BASIC"))
		}

		if authorized_token == authorization {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

type UserAgentMiddleware struct {
	IAuthenticationMiddleware
}

func (u *UserAgentMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer func() {}()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)

	})
}
