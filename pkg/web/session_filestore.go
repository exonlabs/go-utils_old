package web

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
)

// session file factory, saves session data in server side files
type SessionFileFactory struct {
	SecretKey      string
	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieMaxAge   int
	CookieSecure   bool
	CookieHttpOnly bool
	CookieSameSite http.SameSite
	StorePath      string
}

func NewSessionFileFactory(opts map[string]any) *SessionFileFactory {
	fact := &SessionFileFactory{
		SecretKey:      "",
		CookieName:     "sessionid",
		CookiePath:     "/",
		CookieDomain:   "",
		CookieMaxAge:   604800, // 7 days: 7*24*60*60
		CookieSecure:   false,
		CookieHttpOnly: false,
		CookieSameSite: http.SameSiteStrictMode,
		StorePath:      "",
	}

	if v, ok := opts["secret_key"]; ok {
		fact.SecretKey = v.(string)
	}
	if v, ok := opts["cookie_name"]; ok {
		fact.CookieName = v.(string)
	}
	if v, ok := opts["cookie_path"]; ok {
		fact.CookiePath = v.(string)
	}
	if v, ok := opts["cookie_domain"]; ok {
		fact.CookieDomain = v.(string)
	}
	if v, ok := opts["cookie_maxage"]; ok {
		fact.CookieMaxAge = v.(int)
	}
	if v, ok := opts["cookie_secure"]; ok {
		fact.CookieSecure = v.(bool)
	}
	if v, ok := opts["cookie_httponly"]; ok {
		fact.CookieHttpOnly = v.(bool)
	}
	if v, ok := opts["cookie_samesite"]; ok {
		fact.CookieSameSite = v.(http.SameSite)
	}
	if v, ok := opts["store_path"]; ok {
		fact.StorePath = v.(string)
	}

	return fact
}

func (f *SessionFileFactory) Create(
	r *http.Request, w http.ResponseWriter) SessionStore {
	return &SessionFileStore{
		BaseSessionStore: NewBaseSessionStore(),
		factory:          f,
		rHttp:            r,
		wHttp:            w,
	}
}

type SessionFileStore struct {
	*BaseSessionStore
	factory   *SessionFileFactory
	sessionId string
	rHttp     *http.Request
	wHttp     http.ResponseWriter
}

func (s *SessionFileStore) Load() error {
	// read session-id cookie
	cookie, _ := s.rHttp.Cookie(s.factory.CookieName)
	if cookie != nil && len(cookie.Value) != 0 {
		s.sessionId = cookie.Value
	}
	// create new random session-id if empty after reading cookie
	if s.sessionId == "" {
		token := make([]byte, 16)
		if _, err := rand.Read(token); err != nil {
			return errors.New("unable to create session-id tokens")
		}
		s.sessionId = hex.EncodeToString(token)
	}

	// check storage path
	if f, err := os.Stat(s.factory.StorePath); err != nil || !f.IsDir() {
		return errors.New("invalid store path")
	}

	// read session file if exists
	fpath := filepath.Join(s.factory.StorePath, s.sessionId)
	if _, err := os.Stat(fpath); err != nil {
		return nil
	}
	val, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}

	// decrypt session data if secret_key defined
	if s.factory.SecretKey != "" {
		val, err = s.decrypt([]byte(s.factory.SecretKey), val)
		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(val, &s.DataBuffer); err != nil {
		return err
	}
	return nil
}

func (s *SessionFileStore) Save() error {
	// do nothing if no session-id defined
	if s.sessionId == "" {
		return nil
	}

	if len(s.DataBuffer) > 0 {
		// check storage path
		if f, err := os.Stat(s.factory.StorePath); err != nil || !f.IsDir() {
			return errors.New("invalid store path")
		}

		// serialize session data
		val, err := json.Marshal(s.DataBuffer)
		if err != nil {
			return err
		}

		// encrypt session data if secret_key defined
		if s.factory.SecretKey != "" {
			val, err = s.encrypt([]byte(s.factory.SecretKey), val)
			if err != nil {
				return err
			}
		}

		// write session file
		fpath := filepath.Join(s.factory.StorePath, s.sessionId)
		if err := os.WriteFile(fpath, val, 0664); err != nil {
			os.RemoveAll(fpath)
			return err
		}
	}

	// write session-id cookie
	cookie := &http.Cookie{
		Name:     s.factory.CookieName,
		Value:    s.sessionId,
		Path:     s.factory.CookiePath,
		Domain:   s.factory.CookieDomain,
		MaxAge:   s.factory.CookieMaxAge,
		Secure:   s.factory.CookieSecure,
		HttpOnly: s.factory.CookieHttpOnly,
		SameSite: s.factory.CookieSameSite,
	}
	s.wHttp.Header().Add("Set-Cookie", cookie.String())
	return nil
}

func (s *SessionFileStore) Purge() error {
	// delete session file
	fpath := filepath.Join(s.factory.StorePath, s.sessionId)
	if f, err := os.Stat(s.factory.StorePath); err == nil && f.IsDir() {
		os.Remove(fpath)
	}
	s.sessionId = ""
	s.Reset()

	// delete session-id cookie
	cookie := &http.Cookie{
		Name:   s.factory.CookieName,
		MaxAge: -1,
	}
	s.wHttp.Header().Add("Set-Cookie", cookie.String())
	return nil
}
