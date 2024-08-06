# go-guardian-liveblog

Go package for watching a variety of live blog events and reading them aloud using the operating system's text-to-speech APIs.

## Important

This is MacOS specific right now.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/follow cmd/follow/main.go
```

