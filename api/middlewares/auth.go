package middlewares

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
	"time"
)

type AuthMiddleware struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
	oidcClientCtx   context.Context
	oidcProvider    *oidc.Provider

	clientID       string
	identityIssuer string
}

func NewAuthMiddleware(serviceProvider services.IProvider, serverConfig *config.ServerConfig) *AuthMiddleware {
	middleware := new(AuthMiddleware)

	middleware.validator = validator.New()
	middleware.serviceProvider = serviceProvider

	middleware.clientID = serverConfig.ClientID
	middleware.identityIssuer = serverConfig.IdentityIssuer

	// we only want to try to connect to the real identity
	// provider in prod or stage environments. For test and
	// dev purposes, fake token is going to be used
	if env.TickenEnv.IsProd() || env.TickenEnv.IsStage() {
		middleware.oidcClientCtx = initOIDCClientContext()
		middleware.oidcProvider = initOIDCProvider(middleware.oidcClientCtx, middleware.identityIssuer)
	}

	return middleware
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

func (middleware *AuthMiddleware) Setup(router gin.IRouter) {
	if env.TickenEnv.IsDev() {
		router.Use(middleware.isJWTAuthorizedForDev())
	} else {
		router.Use(middleware.isJWTAuthorized())
	}

}

func isFreeURI(uri string) bool {
	return uri == "/healthz"
}

func (middleware *AuthMiddleware) isJWTAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isFreeURI(c.Request.URL.Path) {
			return
		}

		rawAccessToken := c.GetHeader("Authorization")

		oidcConfig := oidc.Config{
			ClientID: middleware.clientID,
		}

		verifier := middleware.oidcProvider.Verifier(&oidcConfig)
		jwt, err := verifier.Verify(middleware.oidcClientCtx, rawAccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "authorisation failed while verifying the token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("jwt", jwt)
	}
}

func (middleware *AuthMiddleware) isJWTAuthorizedForDev() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isFreeURI(c.Request.URL.Path) {
			return
		}

		rawAccessToken := c.GetHeader("Authorization")
		jwt, err := parseAndVerifyTestJWT(rawAccessToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "authorisation failed while verifying the token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("jwt", jwt)
	}
}

func parseAndVerifyTestJWT(jwt string) (*oidc.IDToken, error) {
	payload, err := getTestJWTPayload(jwt)
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt: %v", err)
	}

	var token oidc.IDToken
	if err := json.Unmarshal(payload, &token); err != nil {
		return nil, fmt.Errorf("oidc: failed to unmarshal claims: %v", err)
	}
	return &token, nil
}

func getTestJWTPayload(jwt string) ([]byte, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) < 3 {
		return nil, fmt.Errorf("test jwt malformed, expected 3 parts got %d", len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("test jwt malformed: %v", err)
	}
	return payload, nil
}
