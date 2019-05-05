package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Message(success bool, message string) map[string]interface{} {
	return map[string]interface{}{"success": success, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, contentType string, err error) {
	var validateErr error

	switch contentType {
	case "application/json":
		w.Header().Add("Content-Type", contentType)
		validateErr = json.NewEncoder(w).Encode(data)
	case "application/xml":
		w.Header().Add("Content-Type", contentType)
		validateErr = xml.NewEncoder(w).Encode(data)
	default:
		validateErr = errors.New("Invalid content type! Only accept application/json or application/xml.")
	}

	if err != nil || validateErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err, validateErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
