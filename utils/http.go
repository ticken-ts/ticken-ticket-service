package utils

import (
	"io"
	"net/http"
)

func ReadHTTPResponseBody(res *http.Response) string {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	return string(resBody)
}
