package viewer

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type Query struct {
	Level     string  `json:"level"`
	Ts        float64 `json:"ts"`
	Caller    string  `json:"caller"`
	Msg       string  `json:"msg"`
	Time      float64 `json:"time"`
	RequestID string  `json:"request_id"`
	Query     string  `json:"query"`
	Args      []args    `json:"args"`
}

type args struct {
	Name    string      `json:"name"`
	Ordinal int         `json:"ordinal"`
	Value   interface{} `json:"value"`
}

// {"level":"info","ts":1525134349.979565,"caller":"graqt/tracer.go:20","msg":"Exec",
// "query":"INSERT INTO `user` (email,age) VALUES (?, ?)",
// "args":[{"Name":"","Ordinal":1,"Value":"hoge1525134349977993408@hoge.com"},{"Name":"","Ordinal":2,"Value":8}],
// "time":0.001446192,
// "request_id":"d8523ba8-9948-4171-90ae-57f3b4649efa"}

type QueryParser struct {
	Queries []Query
	file    io.Reader
}

func NewQueryParser(file io.Reader) *QueryParser {
	return &QueryParser{
		file: file,
	}
}

func (r *QueryParser) Parse() error {
	scanner := bufio.NewScanner(r.file)

	var buf bytes.Buffer
	buf.WriteByte('[')
	lineCount := 0

	for scanner.Scan() {
		b := scanner.Bytes()

		buf.Write(b)
		buf.WriteByte(',')

		lineCount += 1
	}

	data := buf.Bytes()
	data[len(data)-1] = ']'

	query := make([]Query, lineCount)

	if err := json.Unmarshal(data, &query); err != nil {
		return errors.Wrap(err, "Failed to Unmarshal json.")
	}

	r.Queries = query

	return nil
}
