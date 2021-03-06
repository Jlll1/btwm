package wm

import (
	"fmt"

	"github.com/Jlll1/btwm/clients"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func HandleEvent(event xgb.Event, conn *xgb.Conn) {
	switch ev := event.(type) {
	case xproto.ConfigureRequestEvent:
		handleConfigureRequest(ev, conn)
	case xproto.DestroyNotifyEvent:
		handleDestroyNotify(ev, conn)
	case xproto.EnterNotifyEvent:
		handleEnterNotify(ev, conn)
	case xproto.KeyPressEvent:
		handleKeyPress(ev, conn)
	case xproto.MapRequestEvent:
		handleMapRequest(ev, conn)
	case xproto.UnmapNotifyEvent:
		handleUnmapNotify(ev, conn)
	}
}

func handleConfigureRequest(ev xproto.ConfigureRequestEvent, conn *xgb.Conn) {
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

func handleDestroyNotify(ev xproto.DestroyNotifyEvent, conn *xgb.Conn) {
	unmanageWindow(ev.Window, conn)
}

func handleEnterNotify(ev xproto.EnterNotifyEvent, conn *xgb.Conn) {
	targetClient := clients.FindByWindow(ev.Event)
	if targetClient == nil {
		return
	}
	xproto.SetInputFocusChecked(conn, xproto.InputFocusPointerRoot, targetClient.Window, xproto.TimeCurrentTime)
	focusedClient = targetClient
}

func handleKeyPress(ev xproto.KeyPressEvent, conn *xgb.Conn) {
	mask := ev.State & (xproto.ModMaskShift | xproto.ModMask4)
	if command := keyToCommand(int(mask), int(ev.Detail)); command != nil {
		command(conn)
	}
}

func handleMapRequest(ev xproto.MapRequestEvent, conn *xgb.Conn) (err error) {
	if clients.FindByWindow(ev.Window) != nil {
		return fmt.Errorf("attempted to remap a window that has already been mapped")
	}
	err = xproto.ChangeWindowAttributesChecked(
		conn, ev.Window, xproto.CwEventMask, []uint32{xproto.EventMaskEnterWindow}).Check()
	if err != nil {
		return err
	}

	err = xproto.MapWindowChecked(conn, ev.Window).Check()
	if err != nil {
		return err
	}

	client := clients.NewClient(ev.Window, 0, 0, screenWidth, screenHeight)
	err = client.Reconfigure(conn)
	if err != nil {
		return err
	}
	clients.Add(client)

	if isInSplitMode {
		exitSplitMode(conn)
	}

	return nil
}

func handleUnmapNotify(ev xproto.UnmapNotifyEvent, conn *xgb.Conn) {
	unmanageWindow(ev.Window, conn)
}
