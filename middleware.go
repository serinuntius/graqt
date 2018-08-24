package graqt

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const RequestIDKey = "RequestID"

var bodySizeWriterPool = sync.Pool{
	New: func() interface{} {
		return &bodySizeWriter{}
	},
}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		id := newRequestID()
		ctx := setRequestID(r.Context(), id)
		bcw := newBodyCopyWriter(w)
		defer bcw.close()

		next.ServeHTTP(bcw, r.WithContext(ctx))

		RLogger.Info("",
			zap.Duration("time", time.Since(t1)),
			zap.String("request_id", id),
			zap.String("path", r.RequestURI),
			zap.String("method", r.Method),
			zap.Uint64("content_length", bcw.Size()),
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

type bodySizeWriter struct {
	http.ResponseWriter
	size uint64
}

func newBodyCopyWriter(w http.ResponseWriter) *bodySizeWriter {
	bsw := bodySizeWriterPool.Get().(*bodySizeWriter)
	bsw.ResponseWriter = w
	bsw.size = 0
	return bsw
}

func (bcw *bodySizeWriter) close() {
	bodySizeWriterPool.Put(bcw)
}

func (bcw *bodySizeWriter) Write(b []byte) (int, error) {
	size, err := bcw.ResponseWriter.Write(b)
	atomic.AddUint64(&bcw.size, uint64(size))
	return size, err
}

func (bcw *bodySizeWriter) Size() uint64 {
	u := atomic.LoadUint64(&bcw.size)
	return u
}

func newRequestID() string {
	return uuid.NewV4().String()
}

func setRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}
