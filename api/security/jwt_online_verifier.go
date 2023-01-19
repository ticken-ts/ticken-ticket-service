package security

import (
	"context"
	"crypto/tls"
	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/env"
	"time"
)

type JWTOnlineVerifier struct {
	oidcClientCtx context.Context
	oidcProvider  *oidc.Provider

	issuer   string
	clientID string

	verifier *oidc.IDTokenVerifier
}

func NewJWTOnlineVerifier(issuer string, clientID string) *JWTOnlineVerifier {
	jwtVerifier := new(JWTOnlineVerifier)

	jwtVerifier.oidcClientCtx = initOIDCClientContext()
	jwtVerifier.oidcProvider = initOIDCProvider(jwtVerifier.oidcClientCtx, issuer)

	jwtVerifier.issuer = issuer
	jwtVerifier.clientID = clientID

	oidcConfig := oidc.Config{
		ClientID: jwtVerifier.clientID,

		// in stage, if we are running with docker, the issues is emited
		// with localhost:8080, but the url inside the docker network is
		// keycloak:8080
		// TODO: solve this
		SkipIssuerCheck: env.TickenEnv.IsStage(),
	}

	jwtVerifier.verifier = jwtVerifier.oidcProvider.Verifier(&oidcConfig)

	return jwtVerifier
}

func (jwtVerifier *JWTOnlineVerifier) Verify(rawJWT string) (*JWT, error) {
	jwt, err := jwtVerifier.verifier.Verify(jwtVerifier.oidcClientCtx, rawJWT)
	if err != nil {
		return nil, err
	}

	claims := new(Claims)
	err = jwt.Claims(claims)
	if err != nil {
		return nil, err
	}

	uuidSubject, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	token := &JWT{
		Subject:  uuidSubject,
		Email:    claims.Email,
		Username: claims.PreferredUsername,
	}

	return token, nil
}

func initOIDCClientContext() context.Context {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}

	return oidc.ClientContext(context.Background(), client)
}

func initOIDCProvider(oidcClientCtx context.Context, issuer string) *oidc.Provider {
	provider, err := oidc.NewProvider(oidcClientCtx, issuer)
	if err != nil {
		panic(err)
	}

	return provider
}
