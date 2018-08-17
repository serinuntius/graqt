package gui

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type ParameterWidget struct {
	name string
	x, y int
	w, h int
}

func NewParameterWidget(name string, x, y, w, h int) *ParameterWidget {
	// add initialize
	return &ParameterWidget{name: name, x: x, y: y, w: w, h: h}
}

func (w *ParameterWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.Wrap(err, "Failed to SetView")
		}
		v.Title = "Query Parameter"
	}
	return nil
}
