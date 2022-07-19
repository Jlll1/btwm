package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func HandleConfigureRequest(ev xproto.ConfigureRequestEvent, conn *xgb.Conn) {
	fmt.Println("configure")
	configureEvent := xproto.ConfigureNotifyEvent{
		Event:            ev.Window,
		Window:           ev.Window,
		AboveSibling:     0,
		X:                ev.X,
		Y:                ev.Y,
		Width:            ev.Width,
		Height:           ev.Height,
		BorderWidth:      ev.BorderWidth,
		OverrideRedirect: false,
	}
	xproto.SendEventChecked(
		conn, false, ev.Window, xproto.EventMaskStructureNotify, string(configureEvent.Bytes()))
}

func HandleKeyPress(ev xproto.KeyPressEvent) {
	// 'p'
	if ev.Detail == 33 {
		exec.Command("dmenu_run").Run()
	}
}

func HandleMapNotify(ev xproto.MapNotifyEvent, conn *xgb.Conn, screenWidth uint16, screenHeight uint16) {
	fmt.Println("Map")
	attributes, err := xproto.GetWindowAttributes(conn, ev.Window).Reply()
	if err != nil {
		log.Fatal(err)
	}
	if attributes.OverrideRedirect {
		return
	}

	values := []uint32{0, 0, uint32(screenWidth), uint32(screenHeight)}
	xproto.ChangeWindowAttributesChecked(conn, ev.Window, xproto.CwEventMask, values).Check()
}

func main() {
	conn, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	connInfo := xproto.Setup(conn)
	if connInfo == nil {
		log.Fatal("couldn't parse connection info")
	}

	screen := connInfo.DefaultScreen(conn)
	root := screen.Root

	mask := []uint32{
		xproto.EventMaskKeyPress |
			xproto.EventMaskStructureNotify |
			xproto.EventMaskSubstructureRedirect,
	}

	err = xproto.ChangeWindowAttributesChecked(
		conn, root, xproto.CwEventMask, mask).Check()
	if err != nil {
		log.Fatal(err)
	}

	for {
		ev, err := conn.WaitForEvent()
		if err != nil {
			continue
		}
		if ev == nil && err == nil {
			break
		}

		switch event := ev.(type) {
		case xproto.KeyPressEvent:
			HandleKeyPress(event)
		case xproto.ConfigureRequestEvent:
			HandleConfigureRequest(event, conn)
		case xproto.MapNotifyEvent:
			HandleMapNotify(event, conn, screen.WidthInPixels, screen.HeightInPixels)
		}
	}
}
