package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type HelpWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewHelpWidget(name string, x, y, w, h int) *HelpWidget {
	// add initialize
	return &HelpWidget{name: name, x: x, y: y, w: w, h: h}
}

func (w *HelpWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.Wrap(err, "Failed to SetView")
		}
		v.Title = "Help"
		fmt.Fprintf(v, "%d:%d:%d:%d", w.x, w.y, w.x+w.w, w.y+w.h)
	}
	return nil
}
