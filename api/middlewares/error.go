package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/log"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/commonerr"
	"ticken-ticket-service/tickenerr/eventerr"
	"ticken-ticket-service/tickenerr/ticketerr"
	"ticken-ticket-service/tickenerr/usererr"
)

var statusCodesByError = map[uint32]int{
	// user err
	usererr.UserAlreadyExistErrorCode: http.StatusBadRequest,
	usererr.UserNotFoundErrorCode:     http.StatusNotFound,

	// ticket err
	ticketerr.TicketNotFoundErrorCode:             http.StatusNotFound,
	ticketerr.BuyResellErrorCode:                  http.StatusBadRequest,
	ticketerr.CreateResellErrorCode:               http.StatusBadRequest,
	ticketerr.ResellCurrencyNotSupportedErrorCode: http.StatusBadRequest,
	ticketerr.TicketResellNotFoundErrorCode:       http.StatusNotFound,

	// event err
	eventerr.EventNotFoundErrorCode: http.StatusNotFound,

	// ticket err
	commonerr.ElementNotFoundInDatabase: http.StatusNotFound,
}

type ErrorMiddleware struct {
}

func NewErrorMiddleware() *ErrorMiddleware {
	return new(ErrorMiddleware)
}

func (middleware *ErrorMiddleware) Setup(router gin.IRouter) {
	router.Use(middleware.ErrorHandler())
}

func (middleware *ErrorMiddleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			switch v := err.Err.(type) {
			case tickenerr.TickenError:
				if v.UnderlyingErr != nil {
					log.TickenLogger.Error().Msg(v.UnderlyingErr.Error())
				}
				c.JSON(getStatusCode(v.Code), res.Error{Code: v.Code, Message: v.Message})
			default:
				c.JSON(http.StatusInternalServerError, res.Error{Code: 0, Message: "An error occurred"})
			}
		}
	}
}

func getStatusCode(errCode uint32) int {
	statusCode, ok := statusCodesByError[errCode]
	if !ok {
		return http.StatusInternalServerError
	}
	return statusCode
}
