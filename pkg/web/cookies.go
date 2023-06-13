package webapp

// import (
// 	"net/http"
// )

// func (env *ViewEnv) GetCookie(name string) (*http.Cookie, error) {
// 	return env.Request.Cookie(name)
// }

// func (env *ViewEnv) SetCookie(c *http.Cookie) {
// 	env.cooikes[c.Name] = c
// }

// func (env *ViewEnv) GetEncryptionCookie(
// 	key string, name string) (*http.Cookie, error) {
// 	cookie, err := env.GetCookie(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	dec, err := decode(key, cookie.Value)
// 	if err != nil {
// 		return nil, err
// 	}
// 	cookie.Value = dec
// 	return cookie, nil
// }

// func (env *ViewEnv) SetEncryptionCookie(key string, c *http.Cookie) error {
// 	enc, err := encode(key, []byte(c.Value))
// 	if err != nil {
// 		return err
// 	}
// 	c.Value = enc
// 	env.SetCookie(c)
// 	return nil
// }

// func (env *ViewEnv) DelCookie(name string) error {
// 	cookie, err := env.GetCookie(name)
// 	if err != nil {
// 		return err
// 	}
// 	cookie.MaxAge = -1
// 	env.cooikes[cookie.Name] = cookie
// 	return nil
// }
