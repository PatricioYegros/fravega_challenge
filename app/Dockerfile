FROM golang:alpine

WORKDIR /challenge_pyegros/app/
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/src ./src

CMD ["/challenge_pyegros/app/bin/src"]
EXPOSE 8080