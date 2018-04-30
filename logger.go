package graqt

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

var QLogger *zap.Logger // Query logger
var RLogger *zap.Logger // Request logger

const defaultQueryLogPath = "query.log"
const defaultRequestLogPath = "request.log"

func NewLogger(outputPath string) *zap.Logger {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{outputPath}
	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func init() {
	QLogger = NewLogger(defaultQueryLogPath)
	RLogger = NewLogger(defaultRequestLogPath)
}

func SetQueryLogger(outputPath string) {
	makeDir(outputPath)

	QLogger = NewLogger(outputPath)
}

func SetRequestLogger(outputPath string) {
	makeDir(outputPath)
	RLogger = NewLogger(outputPath)
}

func makeDir(outputPath string) {
	dirs := strings.Split(outputPath, "/")
	mkdirs := dirs[:len(dirs)-1]

	if len(mkdirs) > 0 {
		path := strings.Join(mkdirs, "/")
		os.MkdirAll(path, os.ModePerm)
	}
}
