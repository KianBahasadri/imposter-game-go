package main

import (
  "net/http"
  "strings"
  "encoding/json"
  "os"
  "fmt"
  "math/rand/v2"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
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
  info["usernames"] = make([]string, 0)
  info["uservotes"] = make(map[string]string)
  dir, _ := os.Open("wordlists/")
  topics, _ := dir.Readdirnames(0)
  info["topiclist"] = topics
  return &info
}
func initializeGameInfo() *map[string]any {
  info := make(map[string]any)
  info["usernames"] = make([]string, 0)
  info["uservotes"] = make(map[string]string)
  info["round"] = 0
  info["secret"] = nil
  return &info
}

func waitingHubMsgHandler(text string, info *map[string]any) []byte {
  split := strings.Split(text, ",")
  username := split[0]
  fmt.Println(split)
  if len(split) == 1 {
    fmt.Println("setting username")
    usernames, _ := (*info)["usernames"].([]string)
    usernameExists := false
    for _, existingName := range usernames {
      if existingName == username {
        usernameExists = true
        break
      }
    }
    if !usernameExists {
      usernames = append(usernames, username)
      (*info)["usernames"] = usernames
    }
  } else if len(split) == 2 {
    topicVote := split[1]
    topics, _ := (*info)["topiclist"].([]string)
    validVote := false
    for _, topic := range topics {
      if topic == topicVote {
        validVote = true
        break
      }
    }
    if validVote {
      uservotes, _ := (*info)["uservotes"].(map[string]string)
      uservotes[username] = topicVote
    }
  }
  textjson, _ := json.Marshal(info)
  return textjson
}
func gameHubMsgHandler(text string, info *map[string]any) []byte {
  split := strings.Split(text, ",")
  username := split[0]
  fmt.Println(split)
  if len(split) == 1 {
    fmt.Println("setting username")
    usernames, _ := (*info)["usernames"].([]string)
    usernameExists := false
    for _, existingName := range usernames {
      if existingName == username {
        usernameExists = true
        break
      }
    }
    if !usernameExists {
      usernames = append(usernames, username)
      (*info)["usernames"] = usernames
    }
  } else if len(split) == 2 {
    nameVote := split[1]
    usernames, _ := (*info)["usernames"].([]string)
    validVote := false
    for _, name := range usernames {
      if name == nameVote {
        validVote = true
        break
      }
    }
    if validVote {
      uservotes, _ := (*info)["uservotes"].(map[string]string)
      uservotes[username] = nameVote
    }
  }
  textjson, _ := json.Marshal(info)
  return textjson
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
