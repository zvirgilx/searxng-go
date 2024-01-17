Searxng-go is a metasearch engine written in Go, inspired by searxng.

--- 

## Features

Searxng-go implements most of the impressive features of searxng, including but not limited to the following:

- [ ] protect user privacy from online tracking
  - [ ] removing private data from requests
  - [ ] avoid forwarding anything from third-party through search
- [ ] enable multiple search engines, and aggregate their result
- [ ] multiple categories support
- [ ] multilingual support
- [ ] plugin builtin
- [ ] setup own proxy pool
- [ ] ...

## Usage

### API server
To start a api server in develop environment:
```

go run main.go -l debug api -m debug -a :9999
```

Search keyword in command line:
```
go run main.go -l debug search superman
```

### Web server

[WebServer Starter](web/README.md)

### All in one

Use `docker-compose` to start both the `API server` and `Web server`.
```
docker compose up --build
```

## License