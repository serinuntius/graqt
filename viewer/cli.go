package viewer

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/urfave/cli"
)

const (
	Version = "v0.1.0"

	ExitCodeOK             = iota
	ExitCodeParseFlagError
)

var (
	Option option
)

// options
type option struct {
	RequestFile string
	QueryFile   string
	Max         bool
	Min         bool
	Avg         bool
	Sum         bool
	Count       bool
	Uri         bool
	Method      bool
	MaxBody     bool
	MinBody     bool
	AvgBody     bool
	SumBody     bool
	P1          bool
	P50         bool
	P99         bool
	Stddev      bool
}

type CLI struct {
	OutSteam, ErrStream io.Writer
}

func (c *CLI) Run(args []string) int {
	var version bool
	flags := flag.NewFlagSet("hoge", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)
	flags.BoolVar(&version, "version", false, "Print version")

	if err := flags.Parse(args[1:]); err != nil {
		log.Println(err)
		return ExitCodeParseFlagError
	}
	if version {
		fmt.Fprintf(c.ErrStream, "graqt version %s\n", Version)
	}

	return ExitCodeOK
}
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "graqt"
	app.Version = "v0.0.1"
	app.Usage = "A access log and query log profiler"

	flags := newFlags(&Option)
	app.Flags = *flags
	return app
}

func newFlags(option *option) (*[]cli.Flag) {
	return &[]cli.Flag{
		&cli.StringFlag{
			Name:        "request-file,rf",
			Value:    "log/request.log",
			Usage:       "path of request.log",
			Destination: &option.RequestFile,
		},
		&cli.StringFlag{
			Name:        "query-file,qf",
			Value:    "log/query.log",
			Usage:       "path of query.log",
			Destination: &option.QueryFile,
		},
		&cli.BoolFlag{
			Name:        "max",
			Usage:       "sort by max response time",
			Destination: &option.Max,
		},
		&cli.BoolFlag{
			Name:        "min",
			Usage:       "sort by min response time",
			Destination: &option.Min,
		},
		&cli.BoolFlag{
			Name:        "avg",
			Usage:       "sort by avg response time",
			Destination: &option.Avg,
		},
		&cli.BoolFlag{
			Name:        "Sum",
			Usage:       "sort by Sum response time",
			Destination: &option.Sum,
		},
		&cli.BoolFlag{
			Name:        "count",
			Usage:       "sort by count",
			Destination: &option.Count,
		},
		&cli.BoolFlag{
			Name:        "uri",
			Usage:       "sort by uri",
			Destination: &option.Uri,
		},
		&cli.BoolFlag{
			Name:        "method",
			Usage:       "sort by method",
			Destination: &option.Method,
		},
		&cli.BoolFlag{
			Name:        "max-body",
			Usage:       "sort by max body size",
			Destination: &option.MaxBody,
		},
		&cli.BoolFlag{
			Name:        "min-body",
			Usage:       "sort by min body size",
			Destination: &option.MinBody,
		},
		&cli.BoolFlag{
			Name:        "avg-body",
			Usage:       "sort by avg body size",
			Destination: &option.AvgBody,
		},
		&cli.BoolFlag{
			Name:        "sum-body",
			Usage:       "sort by sum body size",
			Destination: &option.SumBody,
		},
		&cli.BoolFlag{
			Name:        "p1",
			Usage:       "sort by 1 percentile response time",
			Destination: &option.P1,
		},
		&cli.BoolFlag{
			Name:        "p50",
			Usage:       "sort by 50 percentile response time",
			Destination: &option.P50,
		},
		&cli.BoolFlag{
			Name:        "p99",
			Usage:       "sort by 99 percentile response time",
			Destination: &option.P99,
		},
		&cli.BoolFlag{
			Name:        "stddev",
			Usage:       "sort by standard deviation response time",
			Destination: &option.Stddev,
		},
	}
}
