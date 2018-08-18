package viewer

import (
	"bytes"
	"testing"
)

func TestQueryParser_Parse(t *testing.T) {
	f := bytes.NewBufferString(
		`{"level":"info","ts":1525134349.979565,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO user (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525134349977993408@hoge.com"},{"Name":"","Ordinal":2,"Value":8}],"time":0.001446192,"request_id":"d8523ba8-9948-4171-90ae-57f3b4649efa"}
{"level":"info","ts":1525134349.9831538,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO user (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525134349980926251@hoge.com"},{"Name":"","Ordinal":2,"Value":65}],"time":0.002207108,"request_id":"378cfdee-08cd-477b-8f5b-eeeb06a1f6ce"}
`)

	rp := NewQueryParser(f)
	if err := rp.Parse(); err != nil {
		t.Fatal(err)
	}
}
