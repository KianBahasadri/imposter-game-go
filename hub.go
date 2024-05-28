package main
import "fmt"
type Hub struct {
  msgHandler func(text string, info *map[string]any) []byte
  info       *map[string]any
  clients    map[*Client]bool
  broadcast  chan []byte
  register   chan *Client
  unregister chan *Client
}

func newHub(msgHandler func(text string, info *map[string]any) []byte, info *map[string]any) *Hub {
  return &Hub{
    msgHandler: msgHandler,
    info:       info,
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
  }
}

func (h *Hub) run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case message := <-h.broadcast:
      fmt.Println(string(message))
      resp := h.msgHandler(string(message), h.info)
      for client := range h.clients {
        select {
        case client.send <- resp:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}
