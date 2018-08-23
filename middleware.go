package graqt

import (
	"bufio"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const RequestIDKey = "RequestID"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		id := newRequestID()
		ctx := setRequestID(r.Context(), id)

		next.ServeHTTP(w, r.WithContext(ctx))
		ww := bufio.NewReader(r.Response.Body)

		RLogger.Info("",
			zap.Duration("time", time.Since(t1)),
			zap.String("request_id", id),
			zap.String("path", r.RequestURI),
			zap.String("method", r.Method),
			zap.Uint64("content_length", uint64(ww.Size())),
		)
	})
}

func RequestIdForGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		t1 := time.Now()

		id := newRequestID()

		c.Set(RequestIDKey, id)

		c.Next()

		RLogger.Info("",
			zap.Duration("time", time.Since(t1)),
			zap.String("request_id", id),
			zap.String("path", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.Uint64("content_length", uint64(c.Writer.Size())),
		)

	}
}

func newRequestID() string {
	return uuid.NewV4().String()
}

func setRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}
