package viewer

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dustin/go-humanize"
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
	Body      uint64  `json:"content-length"`
}

type RequestMinimum struct {
	Time      float64
	RequestID string
}

type RequestParser struct {
	File           io.Reader
	RequestIndexes RequestIndexes
	Setting        *Setting
}

type RequestIndexes []RequestIndex

type RequestIndex struct {
	Requests []RequestMinimum
	Max      time.Duration
	Min      time.Duration
	Avg      time.Duration
	Sum      time.Duration
	P1       time.Duration
	P50      time.Duration
	P99      time.Duration
	Stddev   time.Duration
	Count    int
	Uri      string
	Method   string
	MaxBody  uint64
	MinBody  uint64
	AvgBody  uint64
	SumBody  uint64
}

func (ri *RequestIndex) String() string {
	return fmt.Sprintf("\t%d\t%s\t%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%s\t%s\t%s\t%s",
		ri.Count,
		ri.Method,
		ri.Uri,
		ri.Max.Seconds(),
		ri.Min.Seconds(),
		ri.Avg.Seconds(),
		ri.Sum.Seconds(),
		ri.P1.Seconds(),
		ri.P50.Seconds(),
		ri.P99.Seconds(),
		ri.Stddev.Seconds(),
		humanize.Bytes(ri.MaxBody),
		humanize.Bytes(ri.MinBody),
		humanize.Bytes(ri.AvgBody),
		humanize.Bytes(ri.SumBody))
}

type RequestMap map[string]*RequestIndex

func NewRequestParser(file io.Reader, setting *Setting) *RequestParser {
	return &RequestParser{
		File:    file,
		Setting: setting,
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

		t, err := time.ParseDuration(fmt.Sprintf("%fs", r.Time))
		if err != nil {
			return errors.Wrap(err, "Failed to parse time")
		}


		path := r.Path

		if len(rp.Setting.AggregateRegexps) > 0 {
			for _, re := range rp.Setting.AggregateRegexps {
				if ok := re.Match([]byte(path)); !ok {
					continue
				} else {
					path = re.String()
					break
				}
			}
		}

		ri, ok := rm[path]
		if ok {
			ri.Requests = append(ri.Requests, RequestMinimum{Time: r.Time, RequestID: r.RequestID})

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
			rm[path] = &RequestIndex{
				Requests: []RequestMinimum{{Time: r.Time, RequestID: r.RequestID}},
				Max:      t,
				Min:      t,
				Sum:      t,
				Count:    1,
				Uri:      path,
				Method:   r.Method,
				MaxBody:  r.Body,
				MinBody:  r.Body,
				SumBody:  r.Body,
			}
		}
	}

	ris := make(RequestIndexes, len(rm))

	idx := 0
	for _, ri := range rm {
		ris[idx] = *ri
		idx++
	}

	rp.RequestIndexes = ris

	// TODO avg,p,stddev系は後で数える
	return nil

}
