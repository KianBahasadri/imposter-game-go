package main

type Hub struct {
  msgHandler func(text string, info *Info) []byte
  info       *Info
  clients    map[*Client]bool
  broadcast  chan []byte
  register   chan *Client
  unregister chan *Client
}

func newHub(roomName string) *Hub {
  return &Hub{
    msgHandler: hubMsgHandler,
    info:       initializeInfo(roomName),
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
