package middleware

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(store *SessionStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, ok := store.Get(cookie.Value)
		log.Printf("RequireAuth cookie=%q, found=%v", cookie.Value, ok)
		if !ok {
			log.Printf("session not found: %s", cookie.Value)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("session found: userID=%d", userID)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// GetUserID extracts the authenticated user ID from the request context.
func GetUserID(r *http.Request) (int, bool) {
	v := r.Context().Value(userIDKey)
	id, ok := v.(int)
	return id, ok
}
