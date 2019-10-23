FROM intimatemerger/mecab-ipadic-neologd as builder
FROM golang:1.13.1

WORKDIR /go/src/app

# https://blog.csdn.net/zxy131072/article/details/90514978
RUN apt update && apt -y install sudo gcc-multilib libmecab-dev mecab-ipadic-utf8

ENV GO111MODULE on

RUN groupadd -g 10001 myapp &&\
  useradd -u 10001 -g myapp myapp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# https://studygolang.com/topics/5954
# enable cgo
RUN CGO_ENABLED=1\
  CGO_LDFLAGS="`mecab-config --libs`"\
  CGO_CFLAGS="-I`mecab-config --inc-dir`"\
  GOOS=linux\
  GOARCH=386\
  go build -o /go/bin/app

# [Error message]
# runtime/cgo
# /usr/bin/ld: cannot find -lmecab
# /usr/bin/ld: cannot find -lstdc++
# collect2: error: ld returned 1 exit status

FROM alpine:3.9

COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 50051

USER myapp

ENTRYPOINT ["/go/bin/app"]