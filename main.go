package main
import (
  "net/http"
  "strings"
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
  uri := r.URL.Path
  split := strings.Split(uri, "-")
  roomName := split[1]
  serveWs(activeHubs[roomName], w, r)
}
  
var activeHubs = make(map[string]*Hub)

func main() {
  http.HandleFunc("/", loginHandler)
  http.HandleFunc("/queryRooms", queryRooms)
  http.HandleFunc("/createHub", createHub)
  http.HandleFunc("/game", gameHandler)
  http.HandleFunc("/ws-", websocketRouter)

  http.ListenAndServe(":8080", nil)
}
