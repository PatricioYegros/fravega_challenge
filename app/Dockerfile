FROM golang:alpine AS builder

WORKDIR /challenge_pyegros/app/
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/src ./src/main.go

CMD ["./bin/src"]
EXPOSE 8080