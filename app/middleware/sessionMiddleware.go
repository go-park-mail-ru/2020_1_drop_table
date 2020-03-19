package middleware

import (
	"2020_1_drop_table/owners"
	"2020_1_drop_table/responses"
	"context"
	"fmt"
	"net/http"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := owners.CookieStore.Get(r, owners.CookieName)

		if err != nil {
			errMessage := fmt.Sprintf("err: %s, while getting cookie", err.Error())
			responses.SendServerError(errMessage, w)
			return
		}

		if _, found := session.Values["userID"]; !found {
			session.Values["userID"] = -1
			err = session.Save(r, w)
			if err != nil {
				responses.SendServerError(err.Error(), w)
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))

		next.ServeHTTP(w, r)
	})
}
