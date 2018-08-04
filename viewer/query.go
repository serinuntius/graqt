package viewer

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

const queryMapBuffer = 1024

type Query struct {
	Level     string  `json:"level"`
	Ts        float64 `json:"ts"`
	Caller    string  `json:"caller"`
	Msg       string  `json:"msg"`
	Time      float64 `json:"time"`
	RequestID string  `json:"request_id"`
	Query     string  `json:"query"`
	Args      []args  `json:"args"`
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
	QueryMap QueryMap
	file     io.Reader
}

type QueryIndex struct {
	Queries []Query
	Count   int
}

type QueryMap map[string]*QueryIndex

func NewQueryParser(file io.Reader) *QueryParser {
	return &QueryParser{
		file: file,
	}
}

func (qp *QueryParser) Parse() error {
	dec := json.NewDecoder(qp.file)
	qm := make(QueryMap, queryMapBuffer)

	for {
		var q Query
		if err := dec.Decode(&q); err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "Failed to Decode json.")
		}

		qi, ok := qm[q.RequestID]
		if ok {
			qi.Count++
			qi.Queries = append(qi.Queries, q)
		} else {
			qm[q.RequestID] = &QueryIndex{
				Queries: []Query{q},
				Count:   1,
			}
		}
	}
	qp.QueryMap = qm
	return nil
}
