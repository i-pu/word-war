# TODO: 辞書を探して、辞書の追加
# TODO: k8sに移行したい
# TODO: CI/CDを完全に自動化したい
# FROM iputus/go-mecab-neologd
FROM golang:1.13.4-alpine

WORKDIR /go/src/app

ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# https://studygolang.com/topics/5954
# enable cgo
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/app

# [TODO] Hot reload
# RUN go get github.com/oxequa/realize
# CMD ["realize start"]
EXPOSE 50051
ENTRYPOINT ["/go/bin/app"]