package middleware

import (
	"fmt"
	"ghotos/server/app"
	"ghotos/server/service"
	"net/http"
	"strings"
)

func JWTAuth(a *app.App) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			token := strings.Split(bearerToken, " ")
			if len(token) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error.message": "%v"}`, "missing token")
				return
			}
			jwt := service.Jwt{}
			user, err := jwt.ValidateToken(token[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error.message": "%v"}`, "invalid token")
				return
			}
			a.SetUser(user)
			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}

}
