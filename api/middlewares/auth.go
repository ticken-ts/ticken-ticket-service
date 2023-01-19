package middlewares

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"ticken-ticket-service/api/security"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/log"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
)

type AuthMiddleware struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
	jwtVerifier     security.JWTVerifier
}

func NewAuthMiddleware(serviceProvider services.IProvider, serverConfig *config.ServerConfig, devConfig *config.DevConfig) *AuthMiddleware {
	middleware := new(AuthMiddleware)

	middleware.validator = validator.New()
	middleware.serviceProvider = serviceProvider

	// we only want to try to connect to the real identity
	// provider in prod or stage environments. For test and
	// dev purposes, fake token is going to be used
	if env.TickenEnv.IsDev() || env.TickenEnv.IsTest() {
		middleware.jwtVerifier = security.NewJWTOfflineVerifier(loadDevJWTRSA(devConfig))
	} else {
		middleware.jwtVerifier = security.NewJWTOnlineVerifier(serverConfig.IdentityIssuer, serverConfig.ClientID)
	}

	return middleware
}

func (middleware *AuthMiddleware) Setup(router gin.IRouter) {
	router.Use(middleware.isJWTAuthorized())
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

		jwt, err := middleware.jwtVerifier.Verify(rawAccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "authorisation failed while verifying the token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("jwt", jwt)
	}
}

func loadDevJWTRSA(devConfig *config.DevConfig) *rsa.PrivateKey {
	rsaKey, err := utils.LoadRSA(devConfig.JWTPrivateKey, devConfig.JWTPublicKey)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	return rsaKey
}
