package permissions

import (
	"2020_1_drop_table/responses"
	"github.com/gorilla/sessions"
	"net/http"
)

func CheckAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			session := r.Context().Value("session").(*sessions.Session)

			ownerID, found := session.Values["userID"]
			if !found || ownerID == -1 {
				responses.SendForbidden(w)
				return
			}

			next.ServeHTTP(w, r)
			return
		})

}