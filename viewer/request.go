package viewer

import (
	"bufio"
	"bytes"
	"io"

	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	Level     string  `json:"level"`
	Ts        float64 `json:"ts"`
	Caller    string  `json:"caller"`
	Msg       string  `json:"msg"`
	Time      float64 `json:"time"`
	RequestID string  `json:"request_id"`
	Path      string  `json:"path"`
	Method    string  `json:"method"`
}

type RequestParser struct {
	Requests []Request
	file     io.Reader
}

func NewRequestParser(file io.Reader) *RequestParser {
	return &RequestParser{
		file: file,
	}
}

func (r *RequestParser) Parse() error {
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

	requests := make([]Request, lineCount)

	if err := json.Unmarshal(data, &requests); err != nil {
		return errors.Wrap(err, "Failed to Unmarshal json.")
	}

	r.Requests = requests

	return nil
}
