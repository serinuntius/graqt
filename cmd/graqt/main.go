package main

import (
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"github.com/serinuntius/graqt/gui"
	"github.com/serinuntius/graqt/viewer"
	"golang.org/x/sync/errgroup"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	app := viewer.NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// setup log

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


	// setup gocui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue

	winX, winY := g.Size()

	request := gui.NewRequestWidget("request", 0, 0, winX-1, winY-3)
	help := gui.NewHelpWidget("help", 0, winY-3, winX-1, 2)
	//parameter := NewParameterWidget("parameter", winX/2, winY/2, winX/2-1, winY/2-2)

	g.SetManager(help, request) //, query, parameter)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
