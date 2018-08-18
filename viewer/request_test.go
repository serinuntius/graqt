package viewer

import (
	"bytes"
	"strings"
	"testing"
)

func TestRequestParser_Parse(t *testing.T) {
	f := bytes.NewBufferString(`{"level":"info","ts":1525134407.1862984,"caller":"graqt/middleware.go:25","msg":"","time":0.002084588,"request_id":"5cbceb20-a503-494d-8f0e-334b87c972c4","path":"/user","method":"GET"}`)
	fs := f.String()

	rp := NewRequestParser(f)
	if err := rp.Parse(); err != nil {
		t.Fatal(err)
	}

	rm := rp.RequestMap
	for path, ri := range rm {
		idx := strings.Index(fs, path)
		if idx == -1 {
			t.Errorf("Path should be exsist.")
		}
		if ri.Method != "GET" {
			t.Errorf("Method should be GET. Expect %s. But Got %s.", "GET", ri.Method)
		}
		if ri.Count != 1 {
			t.Errorf("Count should be 1. Expect %d. But Got %d.", 1, ri.Count)
		}
		if len(ri.RequestIDs) == 1 && ri.RequestIDs[0] != "5cbceb20-a503-494d-8f0e-334b87c972c4" {
			t.Errorf("RequestIDs should be exsist. Expect %s. But Got %s.", ri.RequestIDs[0], "5cbceb20-a503-494d-8f0e-334b87c972c4")
		}
	}
}
