package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Client struct {
	Window xproto.Window
	Tag    int
}

var stack []Client

func UnmanageWindow(window xproto.Window) {
	var newStack []Client
	for _, c := range stack {
		if c.Window == window {
			continue
		}
		newStack = append(newStack, c)
	}
	stack = newStack
}

func Pop(tag int, conn *xgb.Conn) {
	if len(stack) < 2 {
		return
	}
	if len(stack) < tag || tag < 0 {
		return
	}

	var mask uint16 = xproto.ConfigWindowStackMode
	values := []uint32{xproto.StackModeAbove}
	var newStack []Client
	var clientToPutOnTop Client
	for _, c := range stack {
		if c.Tag == tag {
			xproto.ConfigureWindowChecked(conn, c.Window, mask, values)
			clientToPutOnTop = c
			continue
		}
		newStack = append(newStack, c)
	}
	stack = append(newStack, clientToPutOnTop)
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

func KillSelectedTag(conn *xgb.Conn) {
	if len(stack) < 1 {
		return
	}
	windowToKill := stack[len(stack)-1].Window
	err := xproto.KillClientChecked(conn, uint32(windowToKill)).Check()
	if err != nil {
		return
	}
	stack = stack[:len(stack)-1]
}

func HandleKeyPress(ev xproto.KeyPressEvent, conn *xgb.Conn) {
	shiftActive := (ev.State & xproto.ModMaskShift) != 0
	superActive := (ev.State & xproto.ModMask4) != 0
	if !superActive {
		return
	}

	switch ev.Detail {
	case 33: // 'p'
		exec.Command("dmenu_run").Run()
	case 54: // 'c'
		if shiftActive {
			KillSelectedTag(conn)
		}
	case 58:
		if shiftActive {
			conn.Close()
			os.Exit(0)
		}
	case 10: // '1'
		Pop(1, conn)
	case 11: // '2'
		Pop(2, conn)
	case 12: // '3'
		Pop(3, conn)
	case 13: // '4'
		Pop(4, conn)
	case 14: // '5'
		Pop(5, conn)
	case 15: // '6'
		Pop(6, conn)
	case 16: // '7'
		Pop(7, conn)
	case 17: // '8'
		Pop(8, conn)
	case 18: // '1'
		Pop(9, conn)
	case 19: // '10'
		Pop(10, conn)
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

	stack = append(stack, Client{ev.Window, len(stack) + 1})
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
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 10, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 11, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 12, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 13, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 14, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 15, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 16, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 17, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 18, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 19, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 54, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 58, xproto.GrabModeAsync, xproto.GrabModeAsync)

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
		case xproto.UnmapNotifyEvent:
			UnmanageWindow(event.Window)
		case xproto.DestroyNotifyEvent:
			UnmanageWindow(event.Window)
		}
	}
}
