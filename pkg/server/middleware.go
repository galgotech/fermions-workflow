package server

import (
	"net/http"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

// Finally we can use gin context in the auth middleware of centrifuge.
func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.NewRandom()
		if err != nil {
			return
		}
		ctx := r.Context()
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: userID.String(),
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

func Logger(log log.Logger, skipPaths []string) gin.HandlerFunc {
	var skip map[string]struct{}
	if length := len(skipPaths); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			timeStamp := time.Now()
			latency := timeStamp.Sub(start)

			if raw != "" {
				path = path + "?" + raw
			}

			log.Info("request",
				"clientIP", c.ClientIP(),
				"timeStamp", timeStamp.Format(time.RFC1123),
				"method", c.Request.Method,
				"path", path,
				"method", c.Request.Proto,
				"statusCode", c.Writer.Status(),
				"latency", latency,
				"userAgent", c.Request.UserAgent(),
				"errorMessage", c.Errors.ByType(gin.ErrorTypePrivate).String(),
				"bodySize", c.Writer.Size(),
			)
		}
	}
}
