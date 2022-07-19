package clients

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var clients []*client
var FocusedClient *client
var IsInSplitMode bool
var ScreenWidth, ScreenHeight uint32

func findClient(predicate func(*client) bool) *client {
	for _, c := range clients {
		if predicate(c) {
			return c
		}
	}
	return nil
}

func MakeTagSplit(tag int, conn *xgb.Conn) {
	if IsInSplitMode {
		ExitSplitMode(conn, ScreenWidth)
	}
	splitMasterClient := FocusedClient
	splitSlaveClient := FindByTag(tag)
	if splitSlaveClient == nil {
		return
	}

	err := splitSlaveClient.PutBelow(splitMasterClient.Window, conn)
	if err != nil {
		return
	}

	splitWidth := ScreenWidth / 2
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
	IsInSplitMode = true
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
	IsInSplitMode = false
}
func ManageWindow(windowToManage xproto.Window, windowX, windowY, windowWidth, windowHeight uint32, conn *xgb.Conn) error {
	err := xproto.ChangeWindowAttributesChecked(
		conn, windowToManage, xproto.CwEventMask, []uint32{xproto.EventMaskEnterWindow}).Check()
	if err != nil {
		return err
	}

	err = xproto.MapWindowChecked(conn, windowToManage).Check()
	if err != nil {
		return err
	}

	c := NewClient(windowToManage, windowX, windowY, windowWidth, windowHeight)
	err = c.Reconfigure(conn)
	if err != nil {
		return err
	}
	clients = append(clients, c)
	return nil
}

func UnmanageWindow(windowToUnmanage xproto.Window) {
	var newClients []*client
	for _, c := range clients {
		if c.Window != windowToUnmanage {
			newClients = append(newClients, c)
		}
	}
	clients = newClients
}

func FindByWindow(window xproto.Window) *client {
	return findClient(func(c *client) bool { return c.Window == window })
}

func FindByTag(tag int) *client {
	return clients[tag-1]
}

func RemoveClientByWindow(window xproto.Window) {
	var newClients []*client
	for _, c := range clients {
		if c.Window != window {
			newClients = append(newClients, c)
		}
	}
	clients = newClients
}
