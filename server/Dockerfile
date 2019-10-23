FROM golang:1.13.1-alpine as build

# ref https://hub.docker.com/r/smizy/mecab/dockerfile
RUN apk add --update --no-cache build-base

ENV MECAB_VERSION 0.996
ENV IPADIC_VERSION 2.7.0-20070801
ENV mecab_url https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7cENtOXlicTFaRUE
ENV ipadic_url https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7MWVlSDBCSXZMTXM
ENV build_deps 'curl git bash file sudo openssh'
ENV dependencies 'openssl'

RUN apk add --update --no-cache ${build_deps} \
  # Install dependencies
  && apk add --update --no-cache ${dependencies} \
  # Install MeCab
  && curl -SL -o mecab-${MECAB_VERSION}.tar.gz ${mecab_url} \
  && tar zxf mecab-${MECAB_VERSION}.tar.gz \
  && cd mecab-${MECAB_VERSION} \
  && ./configure --enable-utf8-only --with-charset=utf8 \
  && make \
  && make install \
  && cd \
  # Install IPA dic
  && curl -SL -o mecab-ipadic-${IPADIC_VERSION}.tar.gz ${ipadic_url} \
  && tar zxf mecab-ipadic-${IPADIC_VERSION}.tar.gz \
  && cd mecab-ipadic-${IPADIC_VERSION} \
  && ./configure --with-charset=utf8 \
  && make \
  && make install \
  && cd \
  # Install Neologd
  && git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git \
  && mecab-ipadic-neologd/bin/install-mecab-ipadic-neologd -n -y \
  # Clean up
  && apk del ${build_deps} \
  && rm -rf \
    mecab-${MECAB_VERSION}* \
    mecab-${IPADIC_VERSION}* \
    mecab-ipadic-neologd

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

# prd
FROM alpine:latest AS prd
COPY --from=build /go/bin/app /go/bin/app
WORKDIR /go/src/app
COPY . .
EXPOSE 50051
ENTRYPOINT ["/go/bin/app"]

# dev
FROM golang:1.13.1 AS dev
EXPOSE 50051
WORKDIR /go/src/app
ENV GO111MODULE=on
COPY . .
# [TODO] Hot reload
# RUN go get github.com/oxequa/realize
# CMD ["realize start"]
ENTRYPOINT ["/go/bin/app"]