package web

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
)

type SessionFactory interface {
	// create new session store handler
	Create(*http.Request, http.ResponseWriter) SessionStore
}

type SessionStore interface {
	// load data from session actual store into data buffer
	Load() error
	// save data from data buffer into session actual store
	Save() error
	// delete all data buffer and session actual store
	Purge() error

	// return the whole internal data buffer
	Buffer() map[string]any
	// get item by key from data buffer
	Get(string) (any, bool)
	// set item by key in data buffer
	Set(string, any)
	// delete item by key from data buffer
	Del(string)
	// reset data buffer, delete all keys
	Reset()
}

type BaseSessionStore struct {
	DataBuffer map[string]any
}

func NewBaseSessionStore() *BaseSessionStore {
	return &BaseSessionStore{
		DataBuffer: make(map[string]any),
	}
}

func (s *BaseSessionStore) Buffer() map[string]any {
	return s.DataBuffer
}

func (s *BaseSessionStore) Get(key string) (any, bool) {
	if s.DataBuffer != nil {
		if value, ok := s.DataBuffer[key]; ok {
			return value, true
		}
	}
	return nil, false
}

func (s *BaseSessionStore) Set(key string, value any) {
	if s.DataBuffer == nil {
		s.DataBuffer = make(map[string]any)
	}
	s.DataBuffer[key] = value
}

func (s *BaseSessionStore) Del(key string) {
	if s.DataBuffer != nil {
		delete(s.DataBuffer, key)
	}
}

func (s *BaseSessionStore) Reset() {
	s.DataBuffer = make(map[string]any)
}

func (s *BaseSessionStore) encrypt(
	key []byte, plainText []byte) ([]byte, error) {

	// create new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func (s *BaseSessionStore) decrypt(
	key []byte, cipherText []byte) ([]byte, error) {

	// length of cipherText must be larger than 16 bytes
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext length is less than 16 bytes")
	}

	// create new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(cipherText)-aes.BlockSize)
	iv := cipherText[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plainText, cipherText[aes.BlockSize:])

	return plainText, nil
}
