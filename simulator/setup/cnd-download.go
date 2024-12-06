package functions

import (
	"net/http"

	"github.com/ermes-labs/api-go/api"
	rc "github.com/ermes-labs/storage-redis/packages/go"
)

func cnd_download(
	w http.ResponseWriter,
	r *http.Request,
	sessionToken *api.SessionToken,
) (err error) {
	// Check that a session exists.
	if sessionToken == nil {
		http.Error(w, "No session available", http.StatusBadRequest)
		return nil
	}

	// Read the file name from the request query parameters.
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing ‘filename’ parameter", http.StatusBadRequest)
		return nil
	}

	// Acquire the session. Ermes will handle the returned error with a
	// “500 Internal Server Error” response or a Redirect in case the error
	// is an “ErrMigratedTo” instance.
	return Node.AcquireSession(
		r.Context(),
		*sessionToken,
		api.NewAcquireSessionOptionsBuilder().AllowOffloading().Build(),
		func() error {
			ks := rc.NewErmesKeySpaces(sessionToken.SessionId)
			// Get the file from the session.
			file, err := redisClient.Get(
				r.Context(), ks.Session(filename)).Result()
			if err != nil {
				return err
			}

			if file == "" {
				// Return an error if the file is not found.
				http.Error(w, "File not found", http.StatusNotFound)
			} else {
				// Derive the IO usage from the file size. Early returns and
				// unset resources will default to the average(!) usage of 1.
				Node.SetResourcesTotalUsage(r.Context(), *sessionToken, map[string]float64{
					"io": deriveIOUsage(len(file)),
				})

				// Write the file in the response.
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write([]byte(file))
			}

			return nil
		})
}
