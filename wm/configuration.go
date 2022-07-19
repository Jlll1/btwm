package wm

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type command func(*xgb.Conn)
type key struct {
	modMask, keyCode int
}

var keybindings map[key]command = make(map[key]command)

func BindKey(modMask, keyCode int, action command) {
	k := key{modMask, keyCode}
	keybindings[k] = action
}

func GrabKeys(conn *xgb.Conn, rootWindow xproto.Window) {
	for k := range keybindings {
		xproto.GrabKey(conn, true, rootWindow, uint16(k.modMask), xproto.Keycode(k.keyCode), xproto.GrabModeAsync, xproto.GrabModeAsync)
	}
}

func keyToCommand(modMask, keyCode int) command {
	return keybindings[key{modMask, keyCode}]
}
