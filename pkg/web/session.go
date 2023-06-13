package webapp

// import (
// 	"encoding/json"
// 	"net/http"
// 	"time"
// )

// const (
// 	sessionName = "session"
// )

// func (env *ViewEnv) LoadSession() (map[string]any, error) {
// 	var decKey string
// 	if v, ok := env.option["secret_key"]; ok {
// 		decKey = v.(string)
// 	}

// 	cooike, err := env.GetEncryptionCookie(decKey, sessionName)
// 	if err != nil && err != http.ErrNoCookie {
// 		return nil, err
// 	}

// 	var data map[string]any
// 	if cooike != nil && len(cooike.Value) != 0 {
// 		// unmarshal json
// 		if err := json.Unmarshal([]byte(cooike.Value), &data); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return data, nil
// }

// func (env *ViewEnv) setSession() error {
// 	var encKey string
// 	if v, ok := env.option["secret_key"]; ok {
// 		encKey = v.(string)
// 	}

// 	session := &http.Cookie{Name: sessionName}
// 	if v, ok := env.option["session_path"]; ok {
// 		session.Path = v.(string)
// 	}

// 	if v, ok := env.option["session_domain"]; ok {
// 		session.Domain = v.(string)
// 	}

// 	if v, ok := env.option["session_expires"]; ok {
// 		session.Expires = v.(time.Time)
// 	}

// 	if v, ok := env.option["session_maxage"]; ok {
// 		session.MaxAge = v.(int)
// 	}

// 	// marshal session to json
// 	j, err := json.Marshal(env.Session)
// 	if err != nil {
// 		return err
// 	}

// 	session.Value = string(j)

// 	return env.SetEncryptionCookie(encKey, session)
// }

// func (env *ViewEnv) DeleteSession() error {
// 	return env.DelCookie(sessionName)
// }
