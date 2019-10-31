package infra

import (
	mecab "github.com/shogo82148/go-mecab"
)

const ipadic = "/usr/local/lib/mecab/dic/mecab-ipadic-neologd"

var Tagger *mecab.MeCab

func init() {
	t, err := mecab.New(map[string]string{"dicdir": ipadic})
	if err != nil {
		panic(err)
	}
	Tagger = &t
}
