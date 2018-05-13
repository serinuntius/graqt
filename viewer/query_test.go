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
	for _, req := range rp.Queries {
		if req.Time == 0.0 {
			t.Errorf("Time should not be 0.0")
		}
		if req.Ts == 0.0 {
			t.Errorf("Time should not be 0.0")
		}
		if req.Caller == "" {
			t.Errorf("Caller should not be empty")
		}
		if req.Level != "info" {
			t.Errorf("Level should be info")
		}
		if req.Msg == "" {
			t.Errorf("Msg should not be empty")
		}
		if req.RequestID == "" {
			t.Errorf("RequestID should not be empty")
		}
		if req.Query == "" {
			t.Errorf("Query should not be empty")
		}

		for _, arg := range req.Args {
			if arg.Name != "" {
				t.Errorf("arg.Name should not be empty")
			}
			switch v := arg.Value.(type) {
			case int:
				if v == 0 {
					t.Errorf("arg.Value should not be 0")
				}
			case float64:
				if v == 0.0 {
					t.Errorf("arg.Value should not be 0.0")
				}
			case string:
				if v == "" {
					t.Errorf("arg.Value should not be empty")
				}
			default:
				t.Errorf("arg.Value should have type!")
			}

		}

	}
}
