package appmiddleware

import (
	"api/internal/requestctx"
	"net/http"

	"github.com/google/uuid"
)

func DevelopmentIdentity(userID uuid.UUID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			ctx := requestctx.WithUserID(
				r.Context(),
				userID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
