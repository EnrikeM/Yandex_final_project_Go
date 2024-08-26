FROM golang:1.23

ENV TODO_PORT=7540
ENV TODO_PASSWORD=12345
ENV TODO_DBFILE=scheduler.db

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY app/ ./app/
COPY docs/ ./docs/  
COPY web/ ./web

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app ./app/cmd/main.go

EXPOSE ${TODO_PORT}

CMD ["/my_app"]