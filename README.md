### `To create Go Application`

1. create main.go
2. run go mod init

### `To import remote library`

1. go get -u github.com/gorilla/mux

### `To Deploy Go app to Heroku`

1. Make Dockerfile
    FROM golang:1.15.6      // golang version
    ENV GO111MODULE=on
    RUN mkdir /app          // have docker environment run these commands
    ADD . /app/
    WORKDIR /app
    RUN go mod download
    RUN go build -o main .
    CMD ["/app/main"]
2. Make heroku.yml
    build:
      docker:
        web: Dockerfile
3. Follow this tutorial [https://www.youtube.com/watch?v=4axmcEZTE7M&t=420s]

### `To deploy Docker container to local machine`

1. docker build -t {new name for docker container} .
2. docker run -t -p 8080:8081 {new name for docker container}

### `go run main.go`

Runs the app in the development mode.\
Open [http://localhost:8081](http://localhost:8081) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.
