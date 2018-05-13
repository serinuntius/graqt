package viewer

import (
	"bytes"
	"testing"
)

func TestRequestParser_Parse(t *testing.T) {
	f := bytes.NewBufferString(
		`{"level":"info","ts":1525134407.1862984,"caller":"graqt/middleware.go:25","msg":"","time":0.002084588,"request_id":"5cbceb20-a503-494d-8f0e-334b87c972c4","path":"/user","method":"GET"}
{"level":"info","ts":1525134407.186596,"caller":"graqt/middleware.go:25","msg":"","time":0.008757381,"request_id":"25cebf1e-d924-492f-bb94-997ced74f8ce","path":"/user","method":"GET"}
`)

	rp := NewRequestParser(f)
	if err := rp.Parse(); err != nil {
		t.Fatal(err)
	}
	for _, req := range rp.Requests {
		if req.Method != "GET" {
			t.Errorf("Method should be GET. Expect %s. But Got %s.","GET", req.Method)
		}
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
		if req.Msg != "" {
			t.Errorf("Msg should be empty")
		}
		if req.RequestID == "" {
			t.Errorf("RequestID should not be empty")
		}
		if req.Path != "/user" {
			t.Errorf("Path should be /user")
		}
	}
}
