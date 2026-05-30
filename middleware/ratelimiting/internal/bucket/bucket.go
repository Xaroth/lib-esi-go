package bucket

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware/authentication"
	"github.com/xaroth/lib-esi-go/request"
)

func GetRequestBucket(req *http.Request) int64 {
	ctx := req.Context()

	requiredScope := request.GetRequiredScope(ctx)

	if token, ok := authentication.GetToken(ctx); ok {
		owner := token.Owner()

		if len(requiredScope) == 0 {
			return int64(-2) // Request is bucketed to the application.
		}

		return owner
	}

	// No authentication token provided.
	// Request is shared across all applications from the same IP.
	return int64(-1)
}
