package permissions

import (
	"2020_1_drop_table/internal/pkg/responses"
	"fmt"
	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"time"
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

func generateCsrfLogic(w http.ResponseWriter) {
	csrf, err := uuid.NewV4()
	if err != nil {
		responses.SendForbidden(w)
		return
	}
	expiresDate := time.Now().Add(time.Hour)
	cookie1 := &http.Cookie{Name: "csrf", Value: csrf.String(), HttpOnly: true, Expires: expiresDate}
	http.SetCookie(w, cookie1)
	w.Header().Set("csrf", csrf.String())
	fmt.Println(cookie1)
}

func SetCSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
			return
		})

}

func CheckCSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			csrf := r.Header.Get("X-Csrf-Token")
			csrfCookie, err := r.Cookie("csrf")
			fmt.Println(csrfCookie, csrf, err)
			if err != nil || csrf == "" || csrfCookie.Value == "" || csrfCookie.Value != csrf {
				responses.SendSingleError("csrf-protection", w)
				return
			}
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
			return
		})

}
