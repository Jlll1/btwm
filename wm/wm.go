package wm

import (
	"log"

	"github.com/Jlll1/btwm/clients"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var isInSplitMode bool
var screenWidth, screenHeight uint32
var focusedClient *clients.Client
var atomWMProtocols, atomWMDeleteWindow xproto.Atom

func InitAtoms(conn *xgb.Conn) {
	getAtom := func(name string) xproto.Atom {
		reply, err := xproto.InternAtom(conn, false, uint16(len(name)), name).Reply()
		if err != nil {
			log.Fatal(err)
		}
		if reply == nil {
			return 0
		}
		return reply.Atom
	}

	atomWMProtocols = getAtom("WM_PROTOCOLS")
	atomWMDeleteWindow = getAtom("WM_DELETE_WINDOW")
}

func SetScreenDimensions(width, height uint32) {
	screenWidth = width
	screenHeight = height
}

func enterSplitMode(tagToSplitOn int, conn *xgb.Conn) {
	if isInSplitMode {
		exitSplitMode(conn)
	}
	splitMasterClient := focusedClient
	splitSlaveClient := clients.FindByTag(tagToSplitOn)
	if splitSlaveClient == nil {
		return
	}

	err := splitSlaveClient.PutBelow(splitMasterClient.Window, conn)
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

func exitSplitMode(conn *xgb.Conn) {
	clientsToReconfigure := clients.FindMany(func(c *clients.Client) bool {
		return c.Size.Width != screenWidth
	})
	for _, c := range clientsToReconfigure {
		c.Position.X = 0
		c.Size.Width = screenWidth
		err := c.Reconfigure(conn)
		if err != nil {
			return
		}
	}
	isInSplitMode = false
}

func focusTag(tag int, conn *xgb.Conn) {
	clientToPutOnTop := clients.FindByTag(tag)
	if clientToPutOnTop == nil {
		return
	}
	if err := clientToPutOnTop.Focus(conn); err != nil {
		return
	}

	if isInSplitMode {
		exitSplitMode(conn)
	}
}

func killFocusedClient(conn *xgb.Conn) {
	if err := focusedClient.Kill(conn, atomWMProtocols, atomWMDeleteWindow); err == nil {
		unmanageWindow(focusedClient.Window, conn)
	}
}

func moveFocusedClientOneTagDown() {
  clients.MoveOneTagDown(focusedClient)
}

func moveFocusedClientOneTagUp() {
  clients.MoveOneTagUp(focusedClient)
}

func unmanageWindow(window xproto.Window, conn *xgb.Conn) {
	clientToRemove := clients.FindByWindow(window)
	if clientToRemove == nil {
		return
	}
	clients.Remove(clientToRemove)
	if isInSplitMode {
		exitSplitMode(conn)
	}
}
