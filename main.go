package main
import (
  "net/http"
  "encoding/json"
  "sort"
  "github.com/stripe/stripe-go/v78"
  "github.com/stripe/stripe-go/v78/checkout/session"
  "fmt"
  "log"
  "os"
  "time"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "templates/login.html")
}

func serveLogo (w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "logo.webp")
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
  activeHubs[roomName] = newHub(roomName)
  go activeHubs[roomName].run()
  time.AfterFunc(30 * time.Minute, func() {delete(activeHubs, roomName)})
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("roomName")
  if err != nil {
    if err == http.ErrNoCookie {
      http.Redirect(w, r, "/", http.StatusSeeOther) 
      return 
    } else {
      log.Printf("Error retrieving cookie: %v", err)
      return
    }
  }
  // check if room actually exists
  for name := range activeHubs {
    if name == cookie.Value {
      http.ServeFile(w, r, "templates/game.html")
      return
    }
  }
  http.Redirect(w, r, "/", http.StatusSeeOther) 
}

func websocketRouter(w http.ResponseWriter, r *http.Request) {
  roomName := r.PathValue("roomName")
  serveWs(activeHubs[roomName], w, r)
}
  
func createCheckoutSession(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("roomName")
  if err != nil {
    log.Printf("Error retrieving cookie: %v", err)
    return
  }
  imposter := activeHubs[cookie.Value].info.Imposter
  domain := "https://imposter.bahasadri.com"
  params := &stripe.CheckoutSessionParams{
    LineItems: []*stripe.CheckoutSessionLineItemParams{
      &stripe.CheckoutSessionLineItemParams{
        // Provide the exact Price ID (for example, pr_1234) of the product you want to sell
        Price: stripe.String("price_1PMI0ZAM8ZGd4LYrBs3ucPtp"),
        Quantity: stripe.Int64(1),
      },
    },
    Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
    SuccessURL: stripe.String(domain + "/game" + "?imposter=" + imposter),
    CancelURL: stripe.String(domain + "/game" + "?imposter= Pay up and find out"),
  }
  s, err := session.New(params)
  if err != nil {
    fmt.Printf("session.New: %v", err)
  }
  http.Redirect(w, r, s.URL, http.StatusSeeOther)
}


var activeHubs = make(map[string]*Hub)

func main() {
  apikey, err := os.ReadFile("keys/stripe.key")
  if err != nil {
    fmt.Printf("Failed to read API key: %v", err)
  }
  stripe.Key = string(apikey[:len(apikey)-1])
  http.HandleFunc("/{$}", loginHandler) // only matches "/"
  http.HandleFunc("/logo.webp", serveLogo)
  http.HandleFunc("/queryRooms", queryRooms)
  http.HandleFunc("POST /createHub", createHub)
  http.HandleFunc("/game", gameHandler)
  http.HandleFunc("/ws/{roomName}", websocketRouter)
  http.HandleFunc("/create-checkout-session", createCheckoutSession)
  
  //log.Fatal(http.ListenAndServe(":8080", nil))
  log.Fatal(http.ListenAndServeTLS(":443", "keys/imposter.bahasadri.com.pem", "keys/imposter.bahasadri.com.key", nil))
}



