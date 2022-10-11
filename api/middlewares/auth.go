package middlewares

import (
	"context"
	"crypto/tls"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
	"time"
)

type AuthMiddleware struct {
	validator       *validator.Validate
	serviceProvider services.Provider
	oidcClientCtx   context.Context
	oidcProvider    *oidc.Provider

	clientID       string
	identityIssuer string
	isDev          bool
}

func NewAuthMiddleware(serviceProvider services.Provider, config *utils.TickenConfig) *AuthMiddleware {
	middleware := new(AuthMiddleware)

	middleware.validator = validator.New()
	middleware.serviceProvider = serviceProvider
	middleware.oidcClientCtx = initOIDCClientContext()

	middleware.isDev = config.IsDev()
	middleware.clientID = config.Config.Server.ClientID
	middleware.identityIssuer = config.Config.Server.IdentityIssuer

	middleware.oidcProvider = initOIDCProvider(middleware.oidcClientCtx, middleware.identityIssuer)

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
	if middleware.isDev {
		router.Use(middleware.isJWTAuthorizedForDev())
	} else {
		router.Use(middleware.isJWTAuthorized())
	}

}

func (middleware *AuthMiddleware) isJWTAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		rawAccessToken := c.GetHeader("Authorization")

		oidcConfig := oidc.Config{
			ClientID:          middleware.clientID,
			SkipClientIDCheck: true,
			SkipIssuerCheck:   true,
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
