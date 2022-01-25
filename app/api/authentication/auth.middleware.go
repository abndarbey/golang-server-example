package authentication

import (
	"context"
	"net/http"
	"orijinplus/app/models"
	"orijinplus/utils/authtoken"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			cookie, _ := r.Cookie("jwt")
			// if err != nil {
			// 	if err == http.ErrNoCookie {
			// 		w.WriteHeader(http.StatusUnauthorized)
			// 		return
			// 	}
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	return
			// }

			if cookie == nil {
				tokenString = r.Header.Get("authorization")
			} else {
				tokenString = cookie.Value
			}

			if tokenString != "" {
				auther, decodeErr := authtoken.Decode(tokenString)
				if decodeErr != nil {
					http.Error(w, "token error", http.StatusForbidden)
					return
				}

				// put it in context
				ctx := context.WithValue(r.Context(), userCtxKey, auther)

				// and call the next with our new context
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				r = r.WithContext(r.Context())
				next.ServeHTTP(w, r)
			}
		})
	}
}

// AutherFromContext finds the user from the context. REQUIRES Middleware to have run.
func AutherFromContext(ctx context.Context) *models.Auther {
	auther, _ := ctx.Value(userCtxKey).(*models.Auther)
	return auther
}
