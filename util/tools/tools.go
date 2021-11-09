package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var DATE_FORMAT_COMPARE = "20060102150405000"

//var commonIV = []byte("GEbJOVHUONrWInXe")
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MakeTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
func MakeTime(millisec int64) time.Time {
	return time.Unix(0, millisec*int64(time.Millisecond))
}

func Encrypt(plaintext []byte, key []byte, iv string) (ciphertext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	log.Printf("%v", len(iv))

	nonce := make([]byte, 12)

	if iv == "" {
		_, err = io.ReadFull(rand.Reader, nonce)
		if err != nil {
			return nil, err
		}
	} else {
		nonce = []byte(iv[15:27])

	}

	//log.Printf("%v", len(iv[15:27]))
	//nonce := make([]byte, gcm.NonceSize())
	//nonce := commonIV
	/*
		_, err = io.ReadFull(rand.Reader, nonce)
		if err != nil {
			return nil, err
		}
	*/

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < 12 {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:12],
		ciphertext[12:],
		nil,
	)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
