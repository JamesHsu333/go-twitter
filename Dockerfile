# Initial stage: download modules
FROM golang:1.17-alpine as builder

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go build -o bin/ cmd/api/main.go

EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

VOLUME [ "./assets" ]
ENTRYPOINT ["./bin/main"]