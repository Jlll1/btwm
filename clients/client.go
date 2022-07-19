package clients

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type position struct{ X, Y uint32 }
type size struct{ Width, Height uint32 }
type client struct {
	Window          xproto.Window
	currentPosition position
	currentSize     size
	Position        position
	Size            size
}

func NewClient(window xproto.Window, posX, posY, width, height uint32, conn *xgb.Conn) (c *client, err error) {
	c = &client{
		window,
		position{posX, posY},
		size{width, height},
		position{posX, posY},
		size{width, height},
	}
	var mask uint16 = xproto.ConfigWindowX |
		xproto.ConfigWindowY |
		xproto.ConfigWindowWidth |
		xproto.ConfigWindowHeight
	values := []uint32{c.Position.X, c.Position.Y, c.Size.Width, c.Size.Height}
	err = xproto.ConfigureWindowChecked(conn, window, mask, values).Check()
	return c, err
}

func (c *client) Reconfigure(conn *xgb.Conn) (err error) {
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

func (c *client) Focus(conn *xgb.Conn) (err error) {
	return xproto.ConfigureWindowChecked(
		conn,
		c.Window,
		xproto.ConfigWindowStackMode,
		[]uint32{xproto.StackModeAbove}).Check()
}

func (c *client) PutBelow(window xproto.Window, conn *xgb.Conn) (err error) {
	var mask uint16 = xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode
	values := []uint32{uint32(window), xproto.StackModeBelow}
	return xproto.ConfigureWindowChecked(conn, c.Window, mask, values).Check()
}
