package main

import (
  "net/http"
  "strings"
  "encoding/json"
  "os"
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

func initializeWaitingInfo() *map[string]any {
  info := make(map[string]any)
  info["uservotes"] = make(map[string]string)
  dir, _ := os.Open("wordlists/")
  topics, _ := dir.Readdirnames(0)
  info["topiclist"] = topics
  return &info
}
func initializeGameInfo() *map[string]any {
  info := make(map[string]any)
  return &info
}

func waitingHubMsgHandler(text string, info *map[string]any) []byte {
  split := strings.Split(text, ",")
  username := split[0]
  topicVote := split[1]
  uservotes, _ := (*info)["uservotes"].(map[string]string)
  uservotes[username] = topicVote
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
  waitingInfo := initializeWaitingInfo()
  gameInfo := initializeGameInfo()
  waitingHub := newHub(waitingHubMsgHandler, waitingInfo)
  gameHub := newHub(gameHubMsgHandler, gameInfo)
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
