package permissions

import (
	"2020_1_drop_table/internal/pkg/responses"
	"github.com/gorilla/sessions"
	"net/http"
)

func CheckAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			session := r.Context().Value("session").(*sessions.Session)

			staffID, found := session.Values["userID"]
			if !found || staffID == -1 {
				responses.SendForbidden(w)
				return
			}

			next.ServeHTTP(w, r)
			return
		})

}
