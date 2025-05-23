FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go get github.com/streadway/amqp


COPY . .

RUN go build -o /collector ./cmd/collector/main.go

EXPOSE 8081

CMD ["/collector"]