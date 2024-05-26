package main
import (
  "strings"
	"fmt"
	"net/http"
	"os"
  "io"
  "encoding/json"
  "html/templates"
  "github.com/gorilla/websocket"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/login.html")
}
func waitingRoomHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/waiting.html")
  upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
  }
  conn, _ := upgrader.Upgrade(w, r, nil)
}

func getGameRoom(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    gamePage := templates.ParseFile("game.html")
    gamePage.Execute(w, gameroom)
  case http.MethodPost:
    username := r.Cookie("username")
    votefor := r.FormValue("votefor")
    word := r.FormValue("word")
    if word != "" {
      gameroom.players.username.word = word
    } else if votefor != "" {
      gameroom.players.votefor.votes += 1
    }
  }
}
func getGameInfo(w http.ResponseWriter, r *http.Request) {
  gameinfo := Response{
    Message: Marshal(gameroom)
  }
  json.NewEncoder(w).Encode(gameinfo)
}


type Wordlist struct {
  name string
  votes int
  words []string
}
var wordlists []Wordlist = func() {
  lists, err := os.ReadDir("wordlists")
  var wordlists [len(lists)]Wordlist
  for i:=0;i<len(lists);i++ {
    wordlists[i].name = strings.TrimSuffix(lists[i].Name, ".txt")
    contents, _ = os.ReadFile(lists[i].Name)
    wordlists[i].words = strings.split(contents, "\n")
  }
  return wordlists
}

type WaitingRoom struct {
  sockets map[websocket.Conn]string
  wordlists []Wordlist
  votes int
}
var waitingroom WaitingRoom
waitingroom.wordlists = wordlists

type Player struct {
  username string
  word string
  votes int
  imposter bool
}
type GameRoom struct {
  wordlist []string
  players map[string]Player
  round int
}
// TODO globally instantiate gameroom


func main() {
	http.HandleFunc("/", getRoot)
  http.HandleFunc("/login", getLogin)
  http.HandleFunc("/waiting-room", getWaitingRoom)
  http.HandleFunc("/game-room", getGameRoom)
  fmt.Println("Launching server now")
	http.ListenAndServe(":3333", nil)
}


