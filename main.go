package main

import (
  "net/http"
  "strings"
  "encoding/json"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/login.html")
}
func waitingRoomHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/waiting.html")
}
func gameRoomHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/game.html")
}

func waitingHubMsgHandler(text string, info *map[string]any) []byte {
  split := strings.Split(text, ",")
  username := split[0]
  topicVote := split[1]
  (*info)[username] = topicVote
  textjson, _ := json.Marshal(info)
  return textjson
}
func gameHubMsgHandler(text string, info *map[string]any) []byte {
  /*
  split := strings.Split(text, ",")
  username := split[0]
  key := split[1]
  value := split[2]
  (*info)[key] = topicVote
  textjson, _ := json.Marshal(info)
  return string(textjson)
  */
  return []byte("test")
}

func main() {
  waitingHub := newHub(waitingHubMsgHandler)
  gameHub := newHub(gameHubMsgHandler)
  go waitingHub.run()
  go gameHub.run()
  http.HandleFunc("/", loginHandler)
  http.HandleFunc("/waiting-room", waitingRoomHandler)
  http.HandleFunc("/game-room", gameRoomHandler)

  http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
    serveWs(waitingHub, w, r)
  })
  http.HandleFunc("/ws2", func(w http.ResponseWriter, r *http.Request) {
    serveWs(gameHub, w, r)
  })
  http.ListenAndServe(":8080", nil)
}
