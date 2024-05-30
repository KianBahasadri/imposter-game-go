package main

import (
  "strings"
  "encoding/json"
  "fmt"
  "os"
  "math/rand"
)

type Info struct {
  Usernames   []string
  Topiclist   []string
  Topicvotes  map[string]string
  Playervotes map[string]string
  Words       map[int]map[string]string
  Round       int
  Secret      string
  Imposter    string
  Topic       string
  Voted       string
}

func initializeInfo() *Info {
  var info Info
  dir, _ := os.Open("wordlists/")
  topics, _ := dir.Readdirnames(0)
  info.Topiclist = topics
  info.Usernames = make([]string, 0)
  info.Topicvotes = make(map[string]string)
  info.Playervotes = make(map[string]string)
  info.Words = make(map[int]map[string]string)
  info.Round = 0
  info.Secret = ""
  info.Imposter = ""
  info.Topic = ""
  info.Voted = "Abstain"
  return &info
}

func hubMsgHandler(text string, info *Info) []byte {
  fmt.Println(text)
  fmt.Println(info.Usernames)
  split := strings.Split(text, "::")
  action := split[0]
  username := split[1]

  switch action {
  case "setUsername":
    usernameExists := false
    for _, existingName := range info.Usernames {
      if existingName == username {
        usernameExists = true
        break
      }
    }
    if !usernameExists {
      info.Usernames = append(info.Usernames, username)
    }
  case "Topicvotes":
    vote := split[2]
    info.Topicvotes[username] = vote

    // if everyone voted, set the secret word
    if len(info.Usernames) >= 3 && len(info.Topicvotes) == len(info.Usernames) {
      if info.Secret == "" {
        count := make(map[string]int)
        for _, topic := range info.Topicvotes {
          count[topic] = count[topic] + 1
        }
        _maxTopic := ""
        _maxVotes := 0
        for topic, votes := range count {
          if _maxVotes < votes {
            _maxTopic = topic
            _maxVotes = votes
          }
        }
        info.Topic = _maxTopic
        wordlistBytes, _ := os.ReadFile("wordlists/" + info.Topic)
        wordlist := strings.Split(string(wordlistBytes), "\n")
        randInt := rand.Int() % len(wordlist)
        info.Secret = wordlist[randInt]
        info.Round = 1
        info.Words[1] = make(map[string]string)
        randInt = rand.Int() % len(info.Usernames)
        info.Imposter = info.Usernames[randInt]
      }
    }
  case "Playervotes":
    playerVote := split[2]
    info.Playervotes[username] = playerVote
    // if everyone voted, check if gameEnd and increment Round
    if len(info.Playervotes) == len(info.Usernames) {
      count := make(map[string]int)
      for _, name := range info.Playervotes {
        count[name] = count[name] + 1
      }
      _maxName := ""
      _maxVotes := 0
      for name, votes := range count {
        if _maxVotes < votes {
          _maxName = name
          _maxVotes = votes
        }
      }
      info.Voted = _maxName
      if info.Voted != "Abstain" {
        defer fmt.Println("close this hub")
      }
      info.Round += 1
      info.Words[info.Round] = make(map[string]string)
    }
 
  case "Submitword":
    word := split[2]
    info.Words[info.Round][username] = word
  }
  textjson, _ := json.Marshal(info)
  return textjson
}
