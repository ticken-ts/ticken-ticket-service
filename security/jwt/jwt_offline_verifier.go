package jwt

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type OfflineVerifier struct {
	key *rsa.PrivateKey
}

func NewOfflineVerifier(key *rsa.PrivateKey) *OfflineVerifier {
	return &OfflineVerifier{key: key}
}

func (jwtVerifier *OfflineVerifier) Verify(rawJWT string) (*Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// here we are assuming that the Token
		// is generated with the correct key
		return &jwtVerifier.key.PublicKey, nil
	}

	claims := new(Claims)
	_, err := jwt.ParseWithClaims(rawJWT, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	uuidSubject, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	token := &Token{
		Subject:  uuidSubject,
		Email:    claims.Email,
		Username: claims.PreferredUsername,
	}

	return token, nil
}
