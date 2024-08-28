FROM golang:1.22.0

ENV TODO_PORT=7540
ENV TODO_PASSWORD=12345
ENV TODO_DBFILE=scheduler.db
ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY app/ ./app/
COPY docs/ ./docs/  
COPY web/ ./web

RUN go build -o /todo_list_app ./app/cmd/main.go

EXPOSE ${TODO_PORT}

CMD ["/todo_list_app"]