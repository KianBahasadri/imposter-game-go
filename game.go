package main

import (
  "strings"
  "encoding/json"
  "fmt"
  "os"
)
func initializeInfo() *map[string]any {
  info := make(map[string]any)
  dir, _ := os.Open("wordlists/")
  topics, _ := dir.Readdirnames(0)
  info["topiclist"] = topics
  info["usernames"] = make([]string, 0)
  info["topicvotes"] = make(map[string]string)
  info["playervotes"] = make(map[string]string)
  info["round"] = 0
  info["secret"] = nil
  return &info
}

func hubMsgHandler(text string, info *map[string]any) []byte {
  fmt.Println(text)
  split := strings.Split(text, "::")
  action := split[0]
  username := split[1]
  switch action {
  case "setUsername":
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
  case "voteTopic":
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
      topicvotes, _ := (*info)["topicvotes"].(map[string]string)
      topicvotes[username] = topicVote
    }
  case "votePlayer":
    playerVote := split[1]
    usernames, _ := (*info)["usernames"].([]string)
    validVote := false
    for _, name := range usernames {
      if name == playerVote {
        validVote = true
        break
      }
    }
    if validVote {
      playerVotes, _ := (*info)["playervotes"].(map[string]string)
      playerVotes[username] = playerVote
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

