package gui

import (
	"text/tabwriter"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type QueryWidget struct {
	name string
	x, y int
	w, h int
	view *gocui.View
	tw   *tabwriter.Writer
}

func NewQueryWidget(name string, x, y, w, h int) *QueryWidget {
	// add initialize
	return &QueryWidget{name: name, x: x, y: y, w: w, h: h}
}

func (w *QueryWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.Wrap(err, "Failed to SetView")
		}
		v.Title = "Query Index"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		if err := w.KeyBindings(g); err != nil {
			return errors.Wrap(err, "Failed to Set KeyBindings")
		}

		if _, err := g.SetCurrentView(w.name); err != nil {
			if err != gocui.ErrUnknownView {
				return errors.Wrap(err, "Failed to SetCurrentView")
			}
		}

		w.InitTabWriter(v)

	}
	return nil
}

func (w *QueryWidget) InitTabWriter(v *gocui.View) error {
	w.view = v
	w.tw = tabwriter.NewWriter(v, 0, 8, 2, ' ', 0)
	return nil
}

func (w *QueryWidget) closeWidget(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	// SetCurrentView the request before DeleteView
	if _, err := g.SetCurrentView("request"); err != nil {
		if err != gocui.ErrUnknownView {
			return errors.Wrap(err, "Failed to SetCurrentView")
		}
	}

	if err := g.DeleteView(w.name); err != nil {
		if err != gocui.ErrUnknownView {
			return errors.Wrap(err, "Failed to DeleteView")
		}
	}

	return nil
}

func (w *QueryWidget) KeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeWidget); err != nil {
		return err
	}

	return nil
}
