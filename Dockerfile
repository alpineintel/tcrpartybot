FROM golang:1.11.1-stretch

WORKDIR /go/src/gitlab.com/alpinefresh/tcrpartybot
COPY . .
WORKDIR tcrpartybot
RUN go get
RUN go build -o tcrpartybot *.go
RUN cp tcrpartybot /usr/bin/tcrpartybot

CMD tcrpartybot
