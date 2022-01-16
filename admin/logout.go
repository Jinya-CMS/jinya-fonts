package admin

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == nil {
		cookie.Expires = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
