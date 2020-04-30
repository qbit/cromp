package main

import (
	"log"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

type state struct {
	entryList layout.List
	editor    widget.Editor
}

func main() {
	appState.entryList.Axis = layout.Vertical
	appState.entryList.Alignment = layout.Start

	appState.editor.Submit = true

	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

var appState state

func loop(w *app.Window) error {
	gofont.Register()
	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)

			render(gtx, th)

			e.Frame(gtx.Ops)
		}
	}
}

func render(gtx *layout.Context, th *material.Theme) {
}
