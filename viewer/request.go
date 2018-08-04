package viewer

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

const (
	mapBuffer = 1024
)

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	Level     string  `json:"level"`
	Ts        float64 `json:"ts"`
	Caller    string  `json:"caller"`
	Msg       string  `json:"msg"`
	Time      float64 `json:"time"`
	RequestID string  `json:"request_id"`
	Path      string  `json:"path"`
	Method    string  `json:"method"`
	Body      int // TODO
}

type RequestParser struct {
	File       io.Reader
	RequestMap RequestMap
}

type RequestIndex struct {
	RequestIDs []string
	Max        time.Duration
	Min        time.Duration
	Avg        time.Duration
	Sum        time.Duration
	P1         time.Duration
	P50        time.Duration
	P99        time.Duration
	Stddev     time.Duration
	Count      int
	Uri        string
	Method     string
	MaxBody    int
	MinBody    int
	AvgBody    int
	SumBody    int
}

type RequestMap map[string]*RequestIndex

func NewRequestParser(file io.Reader) *RequestParser {
	return &RequestParser{
		File: file,
	}
}

func (rp *RequestParser) Parse() error {
	dec := json.NewDecoder(rp.File)
	rm := make(RequestMap, mapBuffer)

	for {
		var r Request
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "Failed to Decode json.")
		}

		//log.Println(r)

		t, err := time.ParseDuration(fmt.Sprintf("%fs", r.Time))
		if err != nil {
			return errors.Wrap(err, "Failed to parse time")
		}

		ri, ok := rm[r.Path]
		if ok {
			ri.RequestIDs = append(ri.RequestIDs, r.RequestID)

			if ri.Max < t {
				ri.Max = t
			}

			if ri.Min > t || ri.Min == 0 {
				ri.Min = t
			}

			if ri.MaxBody < r.Body {
				ri.MaxBody = r.Body
			}

			if ri.MinBody > r.Body || ri.MinBody == 0 {
				ri.MinBody = r.Body
			}

			ri.Count += 1
			ri.Sum += t
		} else {
			rm[r.Path] = &RequestIndex{
				RequestIDs: []string{r.RequestID},
				Max:        t,
				Min:        t,
				Sum:        t,
				Count:      1,
				Uri:        r.Path,
				Method:     r.Method,
				// TODO Body size
				MaxBody: r.Body,
				MinBody: r.Body,
				SumBody: r.Body,
			}
		}
	}

	rp.RequestMap = rm
	// TODO avg,p,stddev系は後で数える
	return nil

}
