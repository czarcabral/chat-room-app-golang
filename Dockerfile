FROM golang:1.15.6
ENV GO111MODULE=on
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -o main .
CMD ["/app/main"]
