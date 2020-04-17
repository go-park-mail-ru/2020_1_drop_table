package middleware

import (
	"2020_1_drop_table/internal/pkg/responses"
	"context"
	"fmt"
	"gopkg.in/boj/redistore.v1"
	"net/http"
)

const sessionCookieName = "authCookie"

type sessionMiddleware struct {
	sessionRepo *redistore.RediStore
}

func (s *sessionMiddleware) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionRepo.Get(r, sessionCookieName)
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
				return
			}

		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))

		next.ServeHTTP(w, r)
	})
}
