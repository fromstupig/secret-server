package models

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smapig/secret-server/helpers"
	"golang.org/x/crypto/bcrypt"
)

type Secret struct {
	gorm.Model
	Hash           string     `json:"hash"`
	SecretText     string     `json:"secretText"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	RemainingViews int        `json:"remainingViews"`
}

func (secret *Secret) GenerateHash() {
	hased, _ := bcrypt.GenerateFromPassword([]byte(secret.SecretText), bcrypt.DefaultCost)
	secret.Hash = string(hased)
}

func (secret *Secret) SecretKey() string {
	secretKey := os.Getenv("SECRET_KEY")
	return secretKey + secret.Hash
}

func (secret *Secret) EncryptSecret() {
	secret.SecretText = string(helpers.Encrypt([]byte(secret.SecretText), secret.SecretKey()))
}

func (secret *Secret) DecryptSecret() {
	secret.SecretText = string(helpers.Decrypt([]byte(secret.SecretText), secret.SecretKey()))
}

func (secret *Secret) Create() error {
	secret.GenerateHash()
	secret.EncryptSecret()
	err := DB().Create(secret).Error
	return err
}

func (secret *Secret) DecreaseRemainingViews() error {
	if secret.RemainingViews <= 0 {
		return nil
	}

	secret.RemainingViews = secret.RemainingViews - 1
	err := DB().Save(secret).Error
	return err
}

func (secret *Secret) IsAlive() bool {
	if secret.RemainingViews <= 0 {
		return false
	}

	expiresAt := *(*secret).ExpiresAt
	if !expiresAt.IsZero() && expiresAt.Before(time.Now()) {
		return false
	}

	return true
}

func GetSecretByHash(hash string) *Secret {
	secret := &Secret{}

	if hash == "" {
		return nil
	}

	err := DB().Where(&Secret{Hash: hash}).First(secret).Error

	if err != nil {
		return nil
	}

	return secret
}
