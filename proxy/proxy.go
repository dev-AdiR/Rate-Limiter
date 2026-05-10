package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rate-limiter/Env"
	ratelimit "rate-limiter/rate_limit"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthData struct {
	UserId   float64 `json:"userId"`
	Username string  `json:"username"`
}

func RegisterProxy() {

	env := Env.NewEnv()

	target, _ := url.Parse(env.TargetUrl)
	proxy := httputil.NewSingleHostReverseProxy(target)

	limiter := ratelimit.NewTokenBucket(5, 10)

	mux := http.NewServeMux()

	// public routes (no auth)
	mux.Handle("/api/auth/", LogRequest(rateLimit(limiter, http.HandlerFunc(proxy.ServeHTTP))))

	// protected routes
	handler := verifyToken(env, LogRequest(rateLimit(limiter, http.HandlerFunc(proxy.ServeHTTP))))
	mux.Handle("/api/", handler)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Running reverse proxy on :8081")
	server.ListenAndServe()
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func verifyToken(env *Env.Environment, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, 401, "Missing/invalid auth header")
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(env.JWTSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				writeJSON(w, 401, "Token expired")
			} else {
				writeJSON(w, 401, "Invalid token")
			}
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userId := claims["userId"].(float64)
		username := claims["username"].(string)

		authData := &AuthData{
			UserId:   userId,
			Username: username,
		}

		r.Header.Set("X-User-Id", fmt.Sprintf("%.0f", authData.UserId))
		r.Header.Set("X-Username", authData.Username)

		next.ServeHTTP(w, r)
	})
}

func rateLimit(limiter *ratelimit.TokenBucket, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !limiter.Allow() {
			writeJSON(w, 429, "Rate limit exceeded")
			return
		}
		fmt.Printf("Allowed request: %s %s\n", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
