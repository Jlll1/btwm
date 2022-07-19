package main

import (
	"log"
	"os/exec"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Client struct {
	Tag    int
	Window *xproto.Window
}

var stack []Client

func Pop(tag int, conn *xgb.Conn) {
	if len(stack) < 2 {
		return
	}
	if len(stack) < tag {
		return
	}

	topClient := stack[len(stack)-1]
	var selectedClient Client
	var newStack []Client
	for _, c := range stack {
		if c.Tag == tag {
			selectedClient = c
			var mask uint16 = xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode
			values := []uint32{uint32(*topClient.Window), xproto.StackModeAbove}
			err := xproto.ConfigureWindowChecked(conn, *selectedClient.Window, mask, values).Check()
			if err != nil {
				return
			}
			continue
		}
		newStack = append(newStack, c)
	}
	stack = append(newStack, selectedClient)
}

func HandleConfigureRequest(ev xproto.ConfigureRequestEvent, conn *xgb.Conn) {
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

func HandleKeyPress(ev xproto.KeyPressEvent, conn *xgb.Conn) {
	switch ev.Detail {
	case 33: // 'p'
		exec.Command("dmenu_run").Run()
	case 67: // 'F1'
		Pop(1, conn)
	}
}

func HandleMapRequest(ev xproto.MapRequestEvent, conn *xgb.Conn, screenWidth uint32, screenHeight uint32) (err error) {
	err = xproto.MapWindowChecked(conn, ev.Window).Check()
	if err != nil {
		return err
	}
	var mask uint16 = xproto.ConfigWindowX |
		xproto.ConfigWindowY |
		xproto.ConfigWindowWidth |
		xproto.ConfigWindowHeight
	values := []uint32{0, 0, screenWidth, screenHeight}
	xproto.ConfigureWindowChecked(conn, ev.Window, mask, values)

	stack = append(stack, Client{len(stack), &ev.Window})
	return nil
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

	xproto.GrabKey(conn, true, root, xproto.ModMask4, 33, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 67, xproto.GrabModeAsync, xproto.GrabModeAsync)

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
			HandleKeyPress(event, conn)
		case xproto.ConfigureRequestEvent:
			HandleConfigureRequest(event, conn)
		case xproto.MapRequestEvent:
			HandleMapRequest(event, conn, uint32(screen.WidthInPixels), uint32(screen.HeightInPixels))
		}
	}
}
