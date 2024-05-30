package main

import (
  "strings"
  "encoding/json"
  "fmt"
  "os"
  "math/rand"
  "strconv"
)

type Info struct {
  Usernames   []string
  Topiclist   []string
  Topicvotes  map[string]string
  Playervotes map[string]string
  Round       int
  Secret      string
}

func initializeInfo() Info {
  var info Info
  dir, _ := os.Open("wordlists/")
  topics, _ := dir.Readdirnames(0)
  info.Topiclist = topics
  info.Usernames = make([]string, 0)
  info.Topicvotes = make(map[string]string)
  info.Playervotes = make(map[string]string)
  info.Round = 0
  info.Secret = nil
  return &info
}

func hubMsgHandler(text string, info Info) []byte {
  fmt.Println(text)
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
  case "voteTopic":
    topicVote := split[1]
    validVote := false
    for _, topic := range info.Topiclist {
      if topic == topicVote {
        validVote = true
        break
      }
    }
    if validVote {
      info.Topicvotes[username] = topicVote

      // if everyone voted, set the secret word
      if len(Info.Topicvotes) == len(info.Usernames) {
        if info.Secret != nil {
          count := make(map[string]int)
          for _, topic := range info.Topicvotes {
            count[topic] = count[topic] + 1
          }
          _maxTopic := ""
          _maxVotes := 0
          for topic, votesStr := range topicvotes {
            votes, _ := strconv.Atoi(votesStr)
            if _maxVotes < votes {
              _maxTopic = topic
              _maxVotes = votes
            }
          }
          wordlistBytes, _ := os.ReadFile("wordlists/" + _maxTopic)
          wordlist := strings.Split(string(wordlistBytes), "\n")
          randInt := rand.Int() % len(wordlist)
          info.Secret = wordlist[randInt]
        }
      }
    }
  case "votePlayer":
    playerVote := split[1]
    validVote := false
    for _, name := range info.Usernames {
      if name == info.PlayerVote {
        validVote = true
        break
      }
    }
    if validVote {
      info.PlayerVotes[username] = playerVote
    }
  }
  textjson, _ := json.Marshal(info)
  return textjson
}