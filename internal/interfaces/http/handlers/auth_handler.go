package handlers

import (
	"fmt"
	"net/http"

	"github.com/signcrushmotorist/vent-here/internal/application/auth"
	"github.com/signcrushmotorist/vent-here/internal/application/user"
	"github.com/signcrushmotorist/vent-here/internal/interfaces/http/middleware"
)

type AuthHandler struct {
	Service     *auth.AuthService
	UserService *user.UserService
}

func NewAuthHandler(service *auth.AuthService, userService *user.UserService) *AuthHandler {
	return &AuthHandler{Service: service, UserService: userService}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")

	user, err := h.Service.Register(email, password, username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprintf(w, "user registered successfully! ID: %d, public_id: %s, alias: %s", user.ID, user.PublicID, user.PublicAlias)
}

func (h *AuthHandler) LoginHandler(store *middleware.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.Service.Login(email, password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials"))
			return
		}
		sessionID := store.Create(user.ID)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
		})
		w.Write([]byte("Login successful"))
	}

}
func (h *AuthHandler) ChangeAliasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newAlias := r.FormValue("new_alias")

	err := h.UserService.ChangeAlias(userID, newAlias)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Alias changed successfully: %s", newAlias)
}
