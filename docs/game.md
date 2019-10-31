# ゲーム仕様
## レーティングシステム
基本的には2人用ゲームで用いられる [イロレーティング](https://ja.wikipedia.org/wiki/%E3%82%A4%E3%83%AD%E3%83%AC%E3%83%BC%E3%83%86%E3%82%A3%E3%83%B3%E3%82%B0) を拡張したものを用いる.  対戦相手それぞれ2人での勝ち, 負け, 引き分けでレーティングを順番に更新する.

### 使用パッケージ
- [elo-go](https://github.com/kortemy/elo-go)

### 擬似コード
```go
// Use 1 in case A wins, 0 in case B wins, 0.5 in case of a draw
isWin(a, b) {
  return a.score == b.score ? 0.5 : a.score > b.score ? 1 : 0
}

// Calc ratings
deltas := [4]int64{0, 0, 0, 0}

for i, a := range players {
  for j, b := range players {
    if i == j { 
      continue
    }
    delta := elo.RatingDelta(a.rating, b.rating, isWin(a, b))
    deltas[i] += delta
    deltas[j] -= delta
  }
}

// Update rating
for i, d := range deltas {
  players[i].rating += d
}
```

## しりとりシステム
形態素解析する  

### 使用予定パッケージ
- MeCab
- [go-mecab](https://github.com/shogo82148/go-mecab)
- [mecab-ipadic-NEologd](https://github.com/neologd/mecab-ipadic-neologd/blob/master/README.ja.md)

### 導入記事
<https://blog.nownabe.com/2017/12/25/1228.html>

### 判定方法
- 頭文字がお題と一致
- 文字数がお題と一致
- 「ん」で終わらない
- 1語
- 品詞が名詞
