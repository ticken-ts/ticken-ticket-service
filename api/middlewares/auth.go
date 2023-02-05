package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
)

type AuthMiddleware struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
	jwtVerifier     jwt.Verifier
}

func NewAuthMiddleware(serviceProvider services.IProvider, jwtVerifier jwt.Verifier) *AuthMiddleware {
	middleware := new(AuthMiddleware)

	middleware.validator = validator.New()
	middleware.serviceProvider = serviceProvider
	middleware.jwtVerifier = jwtVerifier

	return middleware
}

func (middleware *AuthMiddleware) Setup(router gin.IRouter) {
	router.Use(middleware.isJWTAuthorized())
}

func isPublicURI(uri string) bool {
	return uri == "/healthz"
}

func (middleware *AuthMiddleware) isJWTAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicURI(c.Request.URL.Path) {
			return
		}

		rawAccessToken := c.GetHeader("Authorization")

		token, err := middleware.jwtVerifier.Verify(rawAccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "authorization failed while verifying the token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("jwt", token)
	}
}
