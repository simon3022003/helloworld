FROM golang:1.16-alpine

RUN apk add git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY handlers ./handlers/

RUN CGO_ENABLED=0 GOOS=linux go build -o /helloworld

EXPOSE 3000

ENV MESSAGE="Simon Shi"

VOLUME [ "/tmp" ]

CMD [ "go", "run", "/app/main.go" ]