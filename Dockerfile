# syntax=docker/dockerfile:1

FROM golang:1.16-alpine
RUN apk add  --no-cache ffmpeg

WORKDIR /

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /

RUN go build -o /youtubez

EXPOSE 8080

CMD [ "/youtubez" ]