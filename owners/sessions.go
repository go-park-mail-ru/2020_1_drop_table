package owners

import (
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	UserID      int
	CookieToken string
	ExpiresDate time.Time
}

type sessionStorage struct {
	sync.Mutex
	sessions []Session
}

func NewSessionsStorage() *sessionStorage {
	return &sessionStorage{}
}

func (s *sessionStorage) Count() int {
	s.Lock()
	defer s.Unlock()
	return len(s.sessions)
}

func (s *sessionStorage) get(index int) Session {
	if len(s.sessions) > index {
		item := s.sessions[index]
		return item
	}
	return Session{}
}

func (s *sessionStorage) Get(index int) Session {
	s.Lock()
	defer s.Unlock()
	return s.get(index)
}

func (s *sessionStorage) createNewSession(userID int, expiresDate time.Time) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	session := Session{
		UserID:      userID,
		CookieToken: u.String(),
		ExpiresDate: expiresDate,
	}
	s.sessions = append(s.sessions, session)
	return u.String(), nil
}

func (s *sessionStorage) CreateNewSession(value Owner, expiresDate time.Time) (string, error) {
	s.Lock()
	defer s.Unlock()
	return s.createNewSession(value.OwnerID, expiresDate)
}

func (s *sessionStorage) Login(email string, password string, expiresDate time.Time) (string, error) {
	existed, owner, _ := Storage.Existed(email, password)
	if !existed {
		err := errors.New("user with given login and password does not exist")
		return "", err
	}
	sessionToken, err := s.CreateNewSession(owner, expiresDate)
	return sessionToken, err
}

func (s *sessionStorage) GetOwnerByCookie(cookie string) (Owner, error) {
	for i := 0; i < s.Count(); i++ {
		session := s.Get(i)
		timeDiff := session.ExpiresDate.Sub(time.Now())
		if session.CookieToken == cookie && timeDiff > 0 {
			return Storage.Get(session.UserID)
		}
	}
	return Storage.Get(-1)
}

func GetAuthCookie(email, password string) (http.Cookie, error) {
	expiresDate := time.Now().Add(time.Hour * 24 * 100).UTC()
	token, err := StorageSession.Login(email, password, expiresDate)

	if err != nil {
		err := errors.New("user with given email and password does not exist")
		return http.Cookie{}, err
	}
	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    token,
		Expires:  expiresDate,
		Path:     "/",
		HttpOnly: true,
	}
	return cookie, nil
}

var StorageSession = NewSessionsStorage()
