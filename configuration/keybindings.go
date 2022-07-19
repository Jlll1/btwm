package configuration

import (
	"github.com/Jlll1/btwm/wm"
	"github.com/jezek/xgb/xproto"
)

const modKey = xproto.ModMask4

func Load() {
	wm.BindKey(modKey, 33, wm.ExecCommand("dmenu_run"))                                 // Mod+P
	wm.BindKey(modKey|xproto.ModMaskShift, 58, wm.ExitWmCommand)                        // Mod+Shift+M
	wm.BindKey(modKey|xproto.ModMaskShift, 54, wm.KillFocusedClientCommand)             // Mod+Shift+C
  wm.BindKey(modKey|xproto.ModMaskShift, 45, wm.MoveFocusedClientOneTagUpCommand())   // Mod+Shift+K
  wm.BindKey(modKey|xproto.ModMaskShift, 44, wm.MoveFocusedClientOneTagDownCommand()) // Mod+Shift+J
	wm.BindKey(modKey, 10, wm.FocusTagCommand(1))                                       // Mod+1
	wm.BindKey(modKey, 11, wm.FocusTagCommand(2))                                       // Mod+2
	wm.BindKey(modKey, 12, wm.FocusTagCommand(3))                                       // Mod+3
	wm.BindKey(modKey, 13, wm.FocusTagCommand(4))                                       // Mod+4
	wm.BindKey(modKey, 14, wm.FocusTagCommand(5))                                       // Mod+5
	wm.BindKey(modKey, 15, wm.FocusTagCommand(6))                                       // Mod+6
	wm.BindKey(modKey, 16, wm.FocusTagCommand(7))                                       // Mod+7
	wm.BindKey(modKey, 17, wm.FocusTagCommand(8))                                       // Mod+8
	wm.BindKey(modKey, 18, wm.FocusTagCommand(9))                                       // Mod+9
	wm.BindKey(modKey, 19, wm.FocusTagCommand(10))                                      // Mod+0
	wm.BindKey(modKey|xproto.ModMaskShift, 10, wm.SplitOnTagCommand(1))                 // Mod+Shift+1
	wm.BindKey(modKey|xproto.ModMaskShift, 11, wm.SplitOnTagCommand(2))                 // Mod+Shift+2
	wm.BindKey(modKey|xproto.ModMaskShift, 12, wm.SplitOnTagCommand(3))                 // Mod+Shift+3
	wm.BindKey(modKey|xproto.ModMaskShift, 13, wm.SplitOnTagCommand(4))                 // Mod+Shift+4
	wm.BindKey(modKey|xproto.ModMaskShift, 14, wm.SplitOnTagCommand(5))                 // Mod+Shift+5
	wm.BindKey(modKey|xproto.ModMaskShift, 15, wm.SplitOnTagCommand(6))                 // Mod+Shift+6
	wm.BindKey(modKey|xproto.ModMaskShift, 16, wm.SplitOnTagCommand(7))                 // Mod+Shift+7
	wm.BindKey(modKey|xproto.ModMaskShift, 17, wm.SplitOnTagCommand(8))                 // Mod+Shift+8
	wm.BindKey(modKey|xproto.ModMaskShift, 18, wm.SplitOnTagCommand(9))                 // Mod+Shift+9
	wm.BindKey(modKey|xproto.ModMaskShift, 19, wm.SplitOnTagCommand(0))                 // Mod+Shift+0
}
