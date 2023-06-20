package web

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// session cookie factory, saves session data in browser cookies
type SessionCookieFactory struct {
	SecretKey      string
	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieMaxAge   int
	CookieSecure   bool
	CookieHttpOnly bool
	CookieSameSite http.SameSite
}

func NewSessionCookieFactory(opts map[string]any) *SessionCookieFactory {
	fact := &SessionCookieFactory{
		SecretKey:      "",
		CookieName:     "session",
		CookiePath:     "/",
		CookieDomain:   "",
		CookieMaxAge:   604800, // 7 days: 7*24*60*60
		CookieSecure:   false,
		CookieHttpOnly: false,
		CookieSameSite: http.SameSiteStrictMode,
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

	return fact
}

func (f *SessionCookieFactory) Create(
	r *http.Request, w http.ResponseWriter) SessionStore {
	return &SessionCookieStore{
		BaseSessionStore: NewBaseSessionStore(),
		factory:          f,
		rHttp:            r,
		wHttp:            w,
	}
}

type SessionCookieStore struct {
	*BaseSessionStore
	factory *SessionCookieFactory
	rHttp   *http.Request
	wHttp   http.ResponseWriter
}

func (s *SessionCookieStore) Load() error {
	cookie, _ := s.rHttp.Cookie(s.factory.CookieName)
	if cookie == nil || len(cookie.Value) == 0 {
		return nil
	}
	val, err := base64.RawStdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}

	// decrypt session data if secret_key is defined
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

func (s *SessionCookieStore) Save() error {
	// do nothing if empty data
	if len(s.DataBuffer) == 0 {
		return nil
	}

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

	cookie := &http.Cookie{
		Name:     s.factory.CookieName,
		Value:    base64.RawStdEncoding.EncodeToString(val),
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

func (s *SessionCookieStore) Purge() error {
	s.Reset()
	cookie := &http.Cookie{
		Name:   s.factory.CookieName,
		MaxAge: -1,
	}
	s.wHttp.Header().Add("Set-Cookie", cookie.String())
	return nil
}
