FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /user1 ./cmd/user1/main.go

EXPOSE 8082

CMD ["/user1"]