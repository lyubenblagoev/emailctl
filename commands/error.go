package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lyubenblagoev/goprsc"
)

func checkErr(err error) {
	if err != nil {
		if e, ok := err.(*goprsc.ErrorResponse); ok {
			statusCode := e.Response.StatusCode
			if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
				CleanAuth()
			}
		}
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
