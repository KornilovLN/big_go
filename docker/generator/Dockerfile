FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go get github.com/streadway/amqp


COPY . .

RUN go build -o /generator ./cmd/generator/main.go

EXPOSE 8080

CMD ["/generator"]