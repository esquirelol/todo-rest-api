FROM golang:1.25-alpine

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go mod tidy
RUN go build -o /app/exe ./cmd/todo_list/main.go

CMD ["/app/exe"]