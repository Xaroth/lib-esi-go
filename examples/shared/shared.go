package shared

import (
	"fmt"
	"net/http"
	"os"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/request"
	"github.com/xaroth/lib-esi-go/transport"
)

var Client = &http.Client{
	Transport: transport.New("lib-esi-go examples", "0.0.0", []string{}, defaults.CompatibilityDate),
}

func Error(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func ShowErrors[T any](resp *request.Response[T]) bool {
	var errorMessage string
	if resp.ErrorData != nil {
		errorMessage = resp.ErrorData.ErrorMessage
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return false
	case http.StatusNotFound:
		Error("Server is currently offline: %s", errorMessage)
	case http.StatusServiceUnavailable:
		Error("Server is currently under maintenance: %s", errorMessage)
	case http.StatusTooManyRequests:
		Error("Too many requests: %s", errorMessage)
	case http.StatusInternalServerError:
		Error("Internal server error: %s", errorMessage)
	case http.StatusBadGateway:
		Error("Bad gateway: %s", errorMessage)
	case http.StatusGatewayTimeout:
		Error("Gateway timeout: %s", errorMessage)
	default:
		fmt.Printf("Unknown status code: %d\n", resp.StatusCode)
	}
	return true
}
