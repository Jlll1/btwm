package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type position struct{ X, Y uint32 }
type size struct{ Width, Height uint32 }
type Client struct {
	Window          xproto.Window
	currentPosition position
	currentSize     size
	Position        position
	Size            size
}

func NewClient(window xproto.Window, posX, posY, width, height uint32) *Client {
	return &Client{
		window,
		position{posX, posY},
		size{width, height},
		position{posX, posY},
		size{width, height},
	}
}

func (c *Client) Reconfigure(conn *xgb.Conn) (err error) {
	var mask uint16
	var values []uint32
	if c.currentPosition.X != c.Position.X {
		mask = mask | xproto.ConfigWindowX
		values = append(values, c.Position.X)
	}
	if c.currentPosition.Y != c.Position.Y {
		mask = mask | xproto.ConfigWindowY
		values = append(values, c.Position.Y)
	}
	if c.currentSize.Width != c.Size.Width {
		mask = mask | xproto.ConfigWindowWidth
		values = append(values, c.Size.Width)
	}
	if c.currentSize.Height != c.Size.Height {
		mask = mask | xproto.ConfigWindowHeight
		values = append(values, c.Size.Height)
	}

	if len(values) > 0 {
		err := xproto.ConfigureWindowChecked(conn, c.Window, mask, values).Check()
		if err == nil {
			c.currentSize = c.Size
			c.currentPosition = c.Position
		}
	}
	return err
}

var clients []*Client

func FindClient(predicate func(*Client) bool) *Client {
	for _, c := range clients {
		if predicate(c) {
			return c
		}
	}
	return nil
}

func FindClientByWindow(window xproto.Window) *Client {
	return FindClient(func(c *Client) bool { return c.Window == window })
}

func FindClientByTag(tag int) *Client {
	return clients[tag-1]
}

var screenWidth, screenHeight uint32
var focusedClient *Client
var isInSplitMode bool

func HandleEnterNotifyEvent(ev xproto.EnterNotifyEvent) {
	targetClient := FindClientByWindow(ev.Event)
	if targetClient == nil {
		return
	}
	focusedClient = targetClient
}

func UnmanageWindow(conn *xgb.Conn, window xproto.Window) {
	var newStack []*Client
	for _, c := range clients {
		if c.Window == window {
			continue
		}
		newStack = append(newStack, c)
	}
	clients = newStack

	if isInSplitMode {
		ExitSplitMode(conn, screenWidth)
	}
}

func Pop(tag int, conn *xgb.Conn) {
	if len(clients) < 2 {
		return
	}
	if len(clients) < tag || tag < 0 {
		return
	}

	clientToPutOnTop := FindClientByTag(tag)
	if clientToPutOnTop == nil {
		return
	}
	var mask uint16 = xproto.ConfigWindowStackMode
	values := []uint32{xproto.StackModeAbove}
	err := xproto.ConfigureWindowChecked(conn, clientToPutOnTop.Window, mask, values).Check()
	if err != nil {
		return
	}

	if isInSplitMode {
		ExitSplitMode(conn, screenWidth)
	}
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

func MakeTagSplit(tag int, conn *xgb.Conn) {
	if len(clients) < 2 {
		return
	}
	if isInSplitMode {
		ExitSplitMode(conn, screenWidth)
	}
	splitMasterClient := focusedClient
	splitSlaveClient := FindClientByTag(tag)
	if splitSlaveClient == nil {
		return
	}

	var mask uint16 = xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode
	values := []uint32{uint32(splitMasterClient.Window), xproto.StackModeBelow}
	err := xproto.ConfigureWindowChecked(conn, splitSlaveClient.Window, mask, values).Check()
	if err != nil {
		return
	}

	splitWidth := screenWidth / 2
	splitMasterClient.Position.X = 0
	splitMasterClient.Size.Width = splitWidth
	err = splitMasterClient.Reconfigure(conn)
	if err != nil {
		return
	}
	splitSlaveClient.Position.X = splitWidth
	splitSlaveClient.Size.Width = splitWidth
	err = splitSlaveClient.Reconfigure(conn)
	if err != nil {
		return
	}
	isInSplitMode = true
}

func ExitSplitMode(conn *xgb.Conn, screenWidth uint32) {
	for _, c := range clients {
		if c.Size.Width != screenWidth {
			c.Position.X = 0
			c.Size.Width = screenWidth
			err := c.Reconfigure(conn)
			if err != nil {
				return
			}
		}
	}
	isInSplitMode = false
}

func KillSelectedTag(conn *xgb.Conn) {
	if len(clients) < 1 {
		return
	}
	err := xproto.KillClientChecked(conn, uint32(focusedClient.Window)).Check()
	if err != nil {
		return
	}

	UnmanageWindow(conn, focusedClient.Window)
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
		if shiftActive {
			MakeTagSplit(1, conn)
			break
		}
		Pop(1, conn)
	case 11: // '2'
		if shiftActive {
			MakeTagSplit(2, conn)
			break
		}
		Pop(2, conn)
	case 12: // '3'
		if shiftActive {
			MakeTagSplit(3, conn)
			break
		}
		Pop(3, conn)
	case 13: // '4'
		if shiftActive {
			MakeTagSplit(4, conn)
			break
		}
		Pop(4, conn)
	case 14: // '5'
		if shiftActive {
			MakeTagSplit(5, conn)
			break
		}
		Pop(5, conn)
	case 15: // '6'
		if shiftActive {
			MakeTagSplit(6, conn)
			break
		}
		Pop(6, conn)
	case 16: // '7'
		if shiftActive {
			MakeTagSplit(7, conn)
			break
		}
		Pop(7, conn)
	case 17: // '8'
		if shiftActive {
			MakeTagSplit(8, conn)
			break
		}
		Pop(8, conn)
	case 18: // '1'
		if shiftActive {
			MakeTagSplit(9, conn)
			break
		}
		Pop(9, conn)
	case 19: // '10'
		if shiftActive {
			MakeTagSplit(10, conn)
			break
		}
		Pop(10, conn)
	}
}

func HandleMapRequest(ev xproto.MapRequestEvent, conn *xgb.Conn, screenWidth uint32, screenHeight uint32) (err error) {
	err = xproto.ChangeWindowAttributesChecked(conn, ev.Window, xproto.CwEventMask, []uint32{xproto.EventMaskEnterWindow}).Check()
	if err != nil {
		return err
	}

	err = xproto.MapWindowChecked(conn, ev.Window).Check()
	if err != nil {
		return err
	}
	var mask uint16 = xproto.ConfigWindowX |
		xproto.ConfigWindowY |
		xproto.ConfigWindowWidth |
		xproto.ConfigWindowHeight
	values := []uint32{0, 0, screenWidth, screenHeight}
	err = xproto.ConfigureWindowChecked(conn, ev.Window, mask, values).Check()
	if err != nil {
		return err
	}
	clients = append(clients, NewClient(ev.Window, 0, 0, screenWidth, screenHeight))

	if isInSplitMode {
		ExitSplitMode(conn, screenWidth)
	}

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
	screenWidth = uint32(screen.WidthInPixels)
	screenHeight = uint32(screen.HeightInPixels)
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
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 10, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 11, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 11, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 12, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 12, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 13, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 13, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 14, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 14, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 15, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 15, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 16, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 16, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 17, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 17, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 18, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 18, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4, 19, xproto.GrabModeAsync, xproto.GrabModeAsync)
	xproto.GrabKey(conn, true, root, xproto.ModMask4|xproto.ModMaskShift, 19, xproto.GrabModeAsync, xproto.GrabModeAsync)

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
			HandleMapRequest(event, conn, uint32(screenWidth), uint32(screenHeight))
		case xproto.UnmapNotifyEvent:
			UnmanageWindow(conn, event.Window)
		case xproto.DestroyNotifyEvent:
			UnmanageWindow(conn, event.Window)
		case xproto.EnterNotifyEvent:
			HandleEnterNotifyEvent(event)
		}
	}
}
