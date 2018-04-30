package graqt

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/go-sql-proxy"
	"go.uber.org/zap"
)

func init() {
	sql.Register("mysql-tracer", proxy.NewProxyContext(&mysql.MySQLDriver{}, &proxy.HooksContext{
		PreExec: func(_ context.Context, _ *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			return time.Now(), nil
		},
		PostExec: func(c context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Result, _ error) error {
			QLogger.Info("Exec",
				zap.String("query", stmt.QueryString),
				zap.Any("args", args),
				zap.Duration("time", time.Since(ctx.(time.Time))),
				zap.String("request_id", c.Value(RequestIDKey).(string)),
			)
			return nil
		},
		PreQuery: func(_ context.Context, _ *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			return time.Now(), nil
		},
		PostQuery: func(c context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Rows, _ error) error {
			QLogger.Info("Query",
				zap.String("query", stmt.QueryString),
				zap.Any("args", args),
				zap.Duration("time", time.Since(ctx.(time.Time))),
				zap.String("request_id", c.Value(RequestIDKey).(string)),
			)
			return nil
		},
	}))
}
