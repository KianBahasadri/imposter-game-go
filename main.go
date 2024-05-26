package main
import (
  "strings"
	"fmt"
	"net/http"
	"os"
  "io"
  "encoding/json"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("username")
  if errors.Is(err, http.ErrNoCookie) {
    http.Redirect(w, r, "/login", http.StatusSeeOther)
  } else
    http.Redirect(w, r, "/waiting-room", http.StatusSeeOther)
  }
}
func getLogin(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // TODO serve static file with form submission for username
  case http.MethodPut:
    usernameCookie, _ := r.Cookie("username")
    usernameCookie.HttpOnly = true
    http.setCookie(w, usernameCookie)
    http.Redirect(w, r, "/waiting-room", http.StatusSeeOther)
  }
}
func getWaitingRoom(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // TODO serve static file for when users are inside waiting room
    // make sure they can vote for a topic by number
  case http.MethodPut:
    voteforstr := r.FormValue("votefor")
    voteforint, _ := strconv.Atoi(voteforstr)
    wordlists[voteforint] += 1
  }
}
func getGameRoom(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // TODO serve static file where users can see the current game state
  case http.MethodPut:
    username := r.Cookie("username")
    votefor := r.FormValue("votefor")
    word := r.FormValue("word")
    if word != "" {
      gameroom.players.username.word = word
    else if votefor != "" {
      gameroom.players.votefor.votes += 1
    }
  }
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
}

type WaitingRoom struct {
  usernames []string
  wordlists wordlists []Wordlist
}
// TODO globally instantiate waitingroom

type Player struct {
  username string
  word string
  votes int
  imposter bool
}
type GameRoom struct {
  wordlist []string
  players = map[string]Player{}
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


