package main

import (
	"log"

	"github.com/Jlll1/btwm/configuration"
	"github.com/Jlll1/btwm/wm"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

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

	err = xproto.ChangeWindowAttributesChecked(
		conn, root, xproto.CwEventMask, []uint32{
			xproto.EventMaskKeyPress |
				xproto.EventMaskStructureNotify |
				xproto.EventMaskSubstructureRedirect,
		}).Check()
	if err != nil {
		log.Fatal(err)
	}

	configuration.Load()
	wm.SetScreenDimensions(uint32(screen.WidthInPixels), uint32(screen.HeightInPixels))
	wm.GrabKeys(conn, root)

	for {
		event, err := conn.WaitForEvent()
		if err != nil {
			continue
		}
		if event == nil && err == nil {
			break
		}
		wm.HandleEvent(event, conn)
	}
}
