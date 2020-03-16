package testsUtils

import (
	"2020_1_drop_table/owners"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func GetAuthCookie(ownerID int) (http.Cookie, error) {
	r := httptest.NewRequest("GET", "/", nil)
	rr := &httptest.ResponseRecorder{}
	session, err := owners.CookieStore.New(r, owners.CookieName)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("auth error: %s", err)
	}
	session.Values["userID"] = ownerID
	fmt.Println("Values ", session.Values, "Session: ", session.ID)

	err = session.Save(r, rr)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("session save error %s", err)

	}

	allCookies := rr.Result().Cookies()
	for _, cookie := range allCookies {
		if cookie.Name == owners.CookieName {
			return *cookie, nil
		}
	}

	return http.Cookie{}, fmt.Errorf(
		"no cookie with given name %s", owners.CookieName)
}
