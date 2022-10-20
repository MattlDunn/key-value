FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /kv-database

EXPOSE 8080

CMD ["/kv-database"]
