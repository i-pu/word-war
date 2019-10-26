FROM wakametech/go-mecab-neologd

WORKDIR /go/src/app

ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# https://studygolang.com/topics/5954
# enable cgo
RUN CGO_ENABLED=1\
  CGO_LDFLAGS="`mecab-config --libs`"\
  CGO_CFLAGS="-I`mecab-config --inc-dir`"\
  GOOS=linux\
  GOARCH=amd64\
  go build -o /go/bin/app

# [TODO] Hot reload
# RUN go get github.com/oxequa/realize
# CMD ["realize start"]
EXPOSE 50051
ENTRYPOINT ["/go/bin/app"]