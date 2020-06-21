package middleware

import (
	"context"
	"encoding/json"
	"net/http"
)

// validate the client jwt token
func IsClientAllowed(auther JwtAuther) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, ok := isTokenValid(auther, r)
			if ok {
				next.ServeHTTP(w, req)
				return
			}
			unAuthorized(w)
		})
	}
}

func isTokenValid(auther JwtAuther, r *http.Request) (*http.Request, bool) {
	if jwtToken := r.Header.Get("Authorization"); len(jwtToken) >= 8 && jwtToken[:7] == "Bearer " {
		token := &AccessToken{}
		if err := auther.RequireLogin(token, jwtToken[7:]); err != nil {
			return nil, false
		}
		req := updateContext(token, r)
		return req, true
	}
	return nil, false
}

func updateContext(token *AccessToken, r *http.Request) *http.Request {
	ctx := r.Context()
	// updatedCtx := context.WithValue(ctx, "UserId", token.UserID)
	updatedCtx := context.WithValue(ctx, "AccessToken", token)
	return r.WithContext(updatedCtx)

}

func unAuthorized(w http.ResponseWriter) {
	code := http.StatusUnauthorized
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(http.StatusText(code)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
