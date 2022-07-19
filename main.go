package main

import (
	"log"
	"os/exec"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func HandleKeyPress(ev xproto.KeyPressEvent) {
	// 'p'
	if ev.Detail == 33 {
		exec.Command("dmenu_run").Run()
	}
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

	root := connInfo.DefaultScreen(conn).Root

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

		switch ev.(type) {
		case xproto.KeyPressEvent:
			HandleKeyPress(ev.(xproto.KeyPressEvent))
		}
	}
}
