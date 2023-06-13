package webapp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func encode(key string, plainText []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("the encryption key has not been set")
	}

	// Create a new AES cipher using the key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	// iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	// Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

func decode(key string, secure string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("the decryption key has not been set")
	}

	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return "", err
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("length of the cipher text is less than 16 Bytes")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
