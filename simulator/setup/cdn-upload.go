package functions

import (
	"io"
	"net/http"
	"time"

	"github.com/ermes-labs/api-go/api"
	rc "github.com/ermes-labs/storage-redis/packages/go"
)

func cdn_upload(
	w http.ResponseWriter,
	r *http.Request,
	sessionToken *api.SessionToken,
) (err error) {
	// Read the file name from the request query parameters.
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing ‘filename’ parameter", http.StatusBadRequest)
		return nil
	}

	// Read the file from the request body.
	file, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(file) == 0 {
		http.Error(w, "Missing body in http message", http.StatusBadRequest)
		return nil
	}

	// Create a session if it does not exists and acquire it. Ermes
	// will handle the returned error with a “500 Internal Server Error”
	// response or a Redirect in case the error is an “ErrMigratedTo”
	// instance.
	return Node.CreateAndAcquireSession(
		r.Context(),
		&sessionToken,
		api.CreateAndAcquireSessionOptions{},
		func(sessionToken api.SessionToken) error {
			ks := rc.NewErmesKeySpaces(sessionToken.SessionId)
			// Derive the IO usage from the file size. Early returns and unset
			// resources will default to the average(!) usage of 1.
			Node.SetResourcesTotalUsage(r.Context(), sessionToken, map[string]float64{
				"io": deriveIOUsage(len(file)),
			})

			// Set the file in the session.
			return redisClient.Set(
				r.Context(), ks.Session(filename), file, time.Hour).Err()
		})
}

func deriveIOUsage(size int) float64 {
	// Derive the IO usage from the file size. Early returns and unset
	// resources will default to the average(!) usage of 1.
	return 0.03 * (1 + float64(size)/1024/5)
}
