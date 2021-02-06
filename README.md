# Features
## Multiple brains
Ingestion and output generation are separate domain concepts. One can generate an arbitrary number of brains at your discretion. Presently, they're stored as files but in theory, database backed brains would also support "brain isolation". Also important to preclude re-ingesting inputs that don't change or do so infrequently.

## Planned:
[ ] Multiple backend support

# Usage:
## Learn
```
$ go run cmd/learn/main.go -help
Usage of /tmp/go-build868195560/b001/exe/main:
  -brainpath string
        path to file containing a brain (default "default.brain")
  -corpus string
        path to directory containing corpus data (default ".")
  -order int
        Ordinality of Markov chains. (default 2)

$ go run cmd/chat/main.go -brain ai.brain -length 45
```
## Chatting
```
$ go run cmd/chat/main.go -help
Usage of /tmp/go-build275738441/b001/exe/main:
  -brain string
        path to brain file (default "default.brain")
  -length int
        desired message length in words (default 25)
$ go run cmd/learn/main.go -corpus path/to/textFiles/ -brainpath name.brain
```
