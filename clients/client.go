package clients

import (
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type position struct{ X, Y uint32 }
type size struct{ Width, Height uint32 }
type Client struct {
	Window          xproto.Window
	Position        position
	Size            size
	currentPosition position
	currentSize     size
	wasConfigured   bool
}

func NewClient(window xproto.Window, posX, posY, width, height uint32) *Client {
	return &Client{
		window,
		position{posX, posY},
		size{width, height},
		position{posX, posY},
		size{width, height},
		false,
	}
}

func (c *Client) Focus(conn *xgb.Conn) (err error) {
	return xproto.ConfigureWindowChecked(
		conn,
		c.Window,
		xproto.ConfigWindowStackMode,
		[]uint32{xproto.StackModeAbove}).Check()
}

func (c *Client) Kill(conn *xgb.Conn, atomWMProtocols, atomWMDeleteWindow xproto.Atom) (err error) {
	prop, err := xproto.GetProperty(
		conn, false, c.Window, atomWMProtocols, xproto.GetPropertyTypeAny, 0, 64).Reply()
	if err != nil {
		return err
	}
	if prop == nil {
		if err = xproto.SetCloseDownModeChecked(conn, xproto.CloseDownDestroyAll).Check(); err != nil {
			return err
		}
		return xproto.DestroyWindowChecked(conn, c.Window).Check()
	}
	for propValue := prop.Value; len(propValue) >= 4; propValue = propValue[4:] {
		val := xproto.Atom(uint32(propValue[0]) | uint32(propValue[1])<<8 | uint32(propValue[2])<<16 | uint32(propValue[3])<<24)
		if val == atomWMDeleteWindow {
			t := time.Now().Unix()
			return xproto.SendEventChecked(
				conn,
				false,
				c.Window,
				xproto.EventMaskNoEvent,
				string(xproto.ClientMessageEvent{
					Format: 32,
					Window: c.Window,
					Type:   atomWMProtocols,
					Data: xproto.ClientMessageDataUnionData32New([]uint32{
						uint32(atomWMDeleteWindow),
						uint32(t),
						0,
						0,
						0,
					}),
				}.Bytes())).Check()
		}
	}
	return nil
}

func (c *Client) PutBelow(window xproto.Window, conn *xgb.Conn) (err error) {
	var mask uint16 = xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode
	values := []uint32{uint32(window), xproto.StackModeBelow}
	return xproto.ConfigureWindowChecked(conn, c.Window, mask, values).Check()
}

func (c *Client) Reconfigure(conn *xgb.Conn) (err error) {
	var mask uint16
	var values []uint32
	if c.currentPosition.X != c.Position.X || !c.wasConfigured {
		mask = mask | xproto.ConfigWindowX
		values = append(values, c.Position.X)
	}
	if c.currentPosition.Y != c.Position.Y || !c.wasConfigured {
		mask = mask | xproto.ConfigWindowY
		values = append(values, c.Position.Y)
	}
	if c.currentSize.Width != c.Size.Width || !c.wasConfigured {
		mask = mask | xproto.ConfigWindowWidth
		values = append(values, c.Size.Width)
	}
	if c.currentSize.Height != c.Size.Height || !c.wasConfigured {
		mask = mask | xproto.ConfigWindowHeight
		values = append(values, c.Size.Height)
	}

	if len(values) > 0 {
		err := xproto.ConfigureWindowChecked(conn, c.Window, mask, values).Check()
		if err == nil {
			c.currentSize = c.Size
			c.currentPosition = c.Position
			if !c.wasConfigured {
				c.wasConfigured = true
			}
		}
	}
	return err
}
