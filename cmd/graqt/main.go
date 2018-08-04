package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/serinuntius/graqt/viewer"
	"golang.org/x/sync/errgroup"
)

func main() {
	app := viewer.NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var eg errgroup.Group

	var rp *viewer.RequestParser
	var qp *viewer.QueryParser

	eg.Go(func() error {
		file, err := os.Open(viewer.Option.RequestFile)
		if err != nil {
			return err
		}
		defer file.Close()

		rp = viewer.NewRequestParser(file)
		if err := rp.Parse(); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		file, err := os.Open(viewer.Option.QueryFile)
		if err != nil {
			return errors.Wrap(err, "Failed to os.Open()")
		}
		defer file.Close()

		qp = viewer.NewQueryParser(file)
		if err := qp.Parse(); err != nil {
			return errors.Wrap(err, "Failed to qp.Parse().")
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

}
