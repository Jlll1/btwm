package clients

import (
	"github.com/jezek/xgb/xproto"
)

var clients []*Client

func Add(clientToAdd *Client) {
	clients = append(clients, clientToAdd)
}

func FindByTag(tag int) *Client {
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
