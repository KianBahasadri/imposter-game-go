package main
import (
  "net/http"
  "encoding/json"
  "sort"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/login.html")
}

func queryRooms(w http.ResponseWriter, r *http.Request) {
  rooms := make([]string, 0)
  for name := range activeHubs {
    rooms = append(rooms, name)
  }
  sort.Strings(rooms)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(rooms)
}

func createHub(w http.ResponseWriter, r *http.Request) {
  roomName := r.FormValue("roomName")
  activeHubs[roomName] = newHub()
  go activeHubs[roomName].run()
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/game.html")
}

func websocketRouter(w http.ResponseWriter, r *http.Request) {
  roomName := r.PathValue("roomName")
  serveWs(activeHubs[roomName], w, r)
}
  
var activeHubs = make(map[string]*Hub)

func main() {
  http.HandleFunc("/{$}", loginHandler) // only matches "/"
  http.HandleFunc("/queryRooms", queryRooms)
  http.HandleFunc("POST /createHub", createHub)
  http.HandleFunc("/game", gameHandler)
  http.HandleFunc("/ws/{roomName}", websocketRouter)

  http.ListenAndServe(":8080", nil)
}
