package clients

import (
	"github.com/jezek/xgb/xproto"
)

var clients []*Client

func Add(clientToAdd *Client) {
	clients = append(clients, clientToAdd)
}

func FindByTag(tag int) *Client {
	if len(clients) < tag {
		return nil
	}
	return clients[tag-1]
}

func FindByWindow(window xproto.Window) *Client {
	return findClient(func(c *Client) bool { return c.Window == window })
}

func FindMany(predicate func(*Client) bool) []*Client {
	var result []*Client
	for _, c := range clients {
		if predicate(c) {
			result = append(result, c)
		}
	}
	return result
}

func MoveOneTagDown(clientToMove *Client) {
  for i, c := range clients {
    if c == clientToMove {
      if i-1 >= 0 {
        clients[i-1], clients[i] = clients[i], clients[i-1]
      }
    }
  }
}

func MoveOneTagUp(clientToMove *Client) {
  for i, c := range clients {
    if c == clientToMove {
      if i+1 < len(clients) {
        clients[i+1], clients[i] = clients[i], clients[i+1]
      }
    }
  }
}

func Remove(clientToRemove *Client) {
	var newClients []*Client
	for _, c := range clients {
		if c != clientToRemove {
			newClients = append(newClients, c)
		}
	}
	clients = newClients
}

func findClient(predicate func(*Client) bool) *Client {
	for _, c := range clients {
		if predicate(c) {
			return c
		}
	}
	return nil
}
