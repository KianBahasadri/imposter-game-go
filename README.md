# imposter-game-go


this is my first time using go
please dont cry

i based the logic off of the hub example from the gorillas library

/:
"please set your username"
button -> go to waiting room

/waiting-room:
"please vote"
form -> lock in vote (sends username and topic vote)
websocket -> live update votes

/game-room:
stateful game cycle until end screen
websocket -> live update etc.
