FROM golang:1.11.1-stretch

WORKDIR /go/src/gitlab.com/alpinefresh/tcrpartybot
COPY . .

# Fix bug described here: https://github.com/ethereum/go-ethereum/issues/2738#issuecomment-365239248
RUN go get github.com/ethereum/go-ethereum
RUN cp -r \
  "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" \
  "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

WORKDIR scanner
RUN go get

RUN go build -o tcrpscanner *.go
RUN cp tcrpscanner /usr/bin/tcrpscanner

CMD tcrpscanner
