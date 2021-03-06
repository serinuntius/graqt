package gui

import (
	"fmt"
	"text/tabwriter"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"github.com/serinuntius/graqt/viewer"
)

type RequestWidget struct {
	name           string
	x, y           int
	w, h           int
	index          int
	view           *gocui.View
	tw             *tabwriter.Writer
	RequestIndexes viewer.RequestIndexes
}

func NewRequestWidget(name string, x, y, w, h int, ri viewer.RequestIndexes) *RequestWidget {
	// add initialize
	return &RequestWidget{name: name, x: x, y: y, w: w, h: h, index: 0, RequestIndexes: ri}
}

func (w *RequestWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Request Index"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		v.SetCursor(0, 1)

		if err := w.KeyBindings(g); err != nil {
			return errors.Wrap(err, "Failed to Set KeyBindings")
		}

		if _, err := g.SetCurrentView(w.name); err != nil {
			return errors.Wrapf(err, "Failed to SetCurrentView. name: %s", w.name)
		}

		w.InitTabWriter(v)

		if _, err := w.PrintHeader(); err != nil {
			return err
		}

		for _, ri := range w.RequestIndexes {
			fmt.Fprintln(w.tw, ri.String())
		}

		w.tw.Flush()
		fmt.Fprintln(w.view)
	}

	return nil
}

func (w *RequestWidget) InitTabWriter(v *gocui.View) error {
	w.view = v
	w.tw = tabwriter.NewWriter(v, 0, 8, 2, ' ', 0)
	return nil
}

func (w *RequestWidget) PrintHeader() (int, error) {
	return fmt.Fprintln(w.tw, "\tCount\tMethod\tPath\tMax\tMin\tAvg\tSum\tP1\tP50\tP99\tStddev\tMaxBody\tMinBody\tAvgBody\tSumBody")
}

func (w *RequestWidget) enter(g *gocui.Gui, v *gocui.View) error {
	winX, winY := g.Size()
	query := NewQueryWidget("query", winX/2, 0, winX-1, winY-3)
	if err := query.Layout(g); err != nil {
		if err != gocui.ErrUnknownView {
			return errors.Wrap(err, "Failed to query.Layout")
		}
	}

	return nil
}

func (w *RequestWidget) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy == 1 {
			// if bottom of header
			v.Clear()

			if _, err := w.PrintHeader(); err != nil {
				return errors.Wrap(err, "Failed to PrintHeader")
			}

			if w.index != 0 {
				w.index--
			}
			for i := w.index; i < len(w.RequestIndexes); i++ {
				fmt.Fprintln(w.tw, w.RequestIndexes[i].String())
			}

			w.tw.Flush()
			fmt.Fprintln(w.view)
		} else {
			if err := v.SetCursor(cx, cy-1); err != nil {
				v.Clear()

				if _, err := w.PrintHeader(); err != nil {
					return errors.Wrap(err, "Failed to PrintHeader")
				}

				if w.index != 0 {
					w.index--
				}

				for i := w.index; i < len(w.RequestIndexes); i++ {
					fmt.Fprintln(w.tw, w.RequestIndexes[i].String())
				}
				w.tw.Flush()
				fmt.Fprintln(w.view)
			}
		}

	}
	return nil
}

func (w *RequestWidget) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			// is bottom ? start
			l, err := v.Line(cy + 1)
			if err != nil {
				return errors.Wrap(err, "Failed to get Line")
			}

			if l == "" {
				// if bottom, set before position
				return v.SetCursor(cy, cy)
			}
			// is bottom ? end

			v.Clear()
			if _, err := w.PrintHeader(); err != nil {
				return errors.Wrap(err, "Failed to PrintHeader")
			}

			if len(w.RequestIndexes) != w.index {
				w.index++
			}

			for i := w.index; i < len(w.RequestIndexes); i++ {
				fmt.Fprintln(w.tw, w.RequestIndexes[i].String())
			}

			w.tw.Flush()
			fmt.Fprintln(w.view)
		}

	}
	return nil
}

func (w *RequestWidget) KeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.enter); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.cursorDown); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.cursorUp); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}
	if err := g.SetKeybinding(w.name, gocui.KeyCtrlN, gocui.ModNone, w.cursorDown); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}
	if err := g.SetKeybinding(w.name, gocui.KeyCtrlP, gocui.ModNone, w.cursorUp); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}
	if err := g.SetKeybinding(w.name, 'j', gocui.ModNone, w.cursorDown); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}
	if err := g.SetKeybinding(w.name, 'k', gocui.ModNone, w.cursorUp); err != nil {
		return errors.Wrap(err, "Failed to SetKeybinding")
	}

	return nil
}
