FROM golang:1.13.1 as builder

WORKDIR /go/src/app

# for mecab
RUN apt update
RUN apt -y install sudo
RUN sudo apt install -y mecab libmecab-dev mecab-ipadic-utf8 git

ENV GO111MODULE on
# for mecab
ENV CGO_LDFLAGS -L/usr/lib/x86_64-linux-gnu -lmecab -lstdc++
ENV CGO_CFLAGS -I/usr/include

RUN groupadd -g 10001 myapp &&\
    useradd -u 10001 -g myapp myapp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app

FROM alpine:3.9


COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 50051

USER myapp

ENTRYPOINT ["/go/bin/app"]