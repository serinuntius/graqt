package graqt

import (
	"context"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type key int

const RequestIDKey key = 0

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		id := newRequestID()
		ctx := setRequestID(r.Context(), id)

		next.ServeHTTP(w, r.WithContext(ctx))

		RLogger.Info("",
			zap.Duration("time", time.Since(t1)),
			zap.String("request_id", id),
			zap.String("path", r.RequestURI),
			zap.String("method", r.Method),
		)
	})
}

func newRequestID() string {
	return uuid.NewV4().String()
}

func setRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}
