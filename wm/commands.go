package wm

import (
	"os"
	"os/exec"

	"github.com/jezek/xgb"
)

func ExecCommand(program string) command {
	return func(conn *xgb.Conn) {
		if conn == nil {
			return
		}
		exec.Command(program).Run()
	}
}

func ExitWmCommand(conn *xgb.Conn) {
	conn.Close()
	os.Exit(0)
}

func FocusTagCommand(tag int) command {
	return func(conn *xgb.Conn) {
		focusTag(tag, conn)
	}
}

func KillFocusedClientCommand(conn *xgb.Conn) {
	killFocusedClient(conn)
}

func SplitOnTagCommand(tag int) command {
	return func(conn *xgb.Conn) {
		enterSplitMode(tag, conn)
	}
}
