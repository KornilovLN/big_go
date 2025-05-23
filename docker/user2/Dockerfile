FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /user2 ./cmd/user2/main.go

EXPOSE 8083

CMD ["/user2"]