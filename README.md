### `To create Go Application`

1. create main.go
2. run `go mod init`

### `To import remote library`

1. `go get -u github.com/gorilla/mux`

### `To Deploy Go app to Heroku`

1. Make Dockerfile
    ```
    FROM golang:1.15.6      // golang version
    ENV GO111MODULE=on
    RUN mkdir /app          // have docker environment run these commands
    ADD . /app/
    WORKDIR /app
    RUN go mod download
    RUN go build -o main .
    CMD ["/app/main"]
    ```
2. Make heroku.yml
    ```
    build:
      docker:
        web: Dockerfile
    ```
3. change ports in go file to access heroku ports
    ```
    // grab the port from heroku's environment variables else default to 5000
    port := os.Getenv("PORT")
    if port == "" {
      port = "5000"
    }

    // serve to localhost
    http.ListenAndServe(":8081", nil)
    ```
4. heroku app -> settings -> stack : change from 'heroku-20' needs to 'container'
    `heroku stack:set container --app chat-room-app-go`
5. heroku app -> deploy -> automatic deployment

5. Follow this tutorial [https://www.youtube.com/watch?v=4axmcEZTE7M&t=420s]

### `To deploy Docker container to local machine`

1. `docker build -t {new name for docker container} .`
2. `docker run -t -p 8080:8081 {new name for docker container}`

### `go run main.go`

Runs the app in the development mode.\
Open [http://localhost:8081](http://localhost:8081) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.
