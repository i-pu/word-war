# TODO: k8sに移行したい
# TODO: CI/CDを完全に自動化したい
# FROM iputus/go-mecab-neologd
FROM iputus/kagome-neologd-dic AS csv


FROM python:latest AS dic

COPY --from=csv /go/src/github.com/ikawaha/kagome/cmd/_dictool/mecab-user-dict-seed.20191212.csv mecab.csv

RUN curl -sOL https://github.com/BurntSushi/xsv/releases/download/0.13.0/xsv-0.13.0-x86_64-unknown-linux-musl.tar.gz \
    && tar -zxvf xsv-0.13.0-x86_64-unknown-linux-musl.tar.gz \
    && mv xsv /usr/local/bin

COPY scripts/make_words.py .

RUN pip install jaconv && python make_words.py mecab.csv /word.txt


FROM golang:1.13.4-alpine

WORKDIR /go/src/app

ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=dic /word.txt /

# https://studygolang.com/topics/5954
# enable cgo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app

# [TODO] Hot reload
# RUN go get github.com/oxequa/realize
# CMD ["realize start"]
EXPOSE 50051

WORKDIR /
ENTRYPOINT ["/go/bin/app"]