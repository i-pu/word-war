# word-war
> Developping a word-chain game service for studying

## ✔ Status
### Testing
[![Go Report Card](https://goreportcard.com/badge/github.com/i-pu/word-war)](https://goreportcard.com/report/github.com/i-pu/word-war)
[![codecov](https://codecov.io/gh/i-pu/word-war/branch/develop%2F1.0/graph/badge.svg)](https://codecov.io/gh/i-pu/word-war)
### CI/CD
[![Actions Status](https://github.com/i-pu/word-war/workflows/Client%20Firebase%20Hosting/badge.svg)](https://github.com/i-pu/word-war/actions)
[![Actions Status](https://github.com/i-pu/word-war/workflows/Server%20Docker%20Build/badge.svg)](https://github.com/i-pu/word-war/actions)
[![Actions Status](https://github.com/i-pu/word-war/workflows/Server%20Test/badge.svg)](https://github.com/i-pu/word-war/actions)

## 💪 Skill Stacks
#### Infra/Server
![go](https://img.shields.io/badge/-Go-76E1FE.svg?logo=go&style=flat-square)
![firebase](https://img.shields.io/badge/-Firebase-000000.svg?logo=firebase&style=flat-square)
![redis](https://img.shields.io/badge/-Redis-D82C20.svg?logo=redis&style=flat-square)
![grpc](https://img.shields.io/badge/-gRPC-47848F.svg?logo=&style=flat-square)
![docker](https://img.shields.io/badge/-Docker-48bcfa.svg?logo=docker&style=flat-square)
#### Frontend
![js](https://img.shields.io/badge/-Javascript-FAEB7F.svg?logo=javascript&style=flat-square)
![bulma](https://img.shields.io/badge/-Bulma-FCEA0.svg?logo=&style=flat-square)

## 🚧 Roadmap
> want to practice scrum development

- v1.0
  - [x] minimum viable product
- v1.1
  - [x] siritori
- v1.2
  - [x] rating
  - [x] Health Check
  - [x] routingのガード
- v1.3
  - [x] redisの初期化といらないキーの掃除
  - [x] ランダムマッチングシステムの実装
- v1.4
  - [x] リファクタリング大会
    - [x] redigo じゃなくて、https://github.com/go-redis/redisを使いたい
- v1.5
  - [ ] ルーム指定のマッチングシステムの実装
    - [ ] マッチングのストリーム
    - [ ] クライアントの待機画面の作成
  - [ ] ワーカーの導入
    - [ ] 部屋を作るのは１つの部屋に対して１回のみなのでそうゆうこと
	  - [ ] 時間制にする
	  - [ ] Critical: 部屋に100人いれば100回UpdateRatingが呼ばれるので部屋に固有のgoroutineを作成し、1回だけ呼ばれるようにしたい
    - see sample
  - [ ] Game画面ブラッシュアップ
    - [ ] 今のスコア表示をしてもいいかも
- v1.6
  - [ ] twitterでツイートをさせる(OGP)
    - <https://qiita.com/yuneco/items/5e526464939082862f5d>
  - [ ] マイページの作成(履歴)
  - [ ] 爽快なエフェクトで表示(UI強化)
    - [ ] 無効な文字列も受け取ってふわっと表示するようしたい
- v?.?
  - [ ] 空文字送られるとぶっ壊れるのでサーバ側でバリデーションをする(client済み)
  - [ ] ちゃんとAPI
  - [ ] 監視、メトリクス
  - [ ] k8s上げたい
  - [ ] CI/CDを完全に自動化したい
  - [ ] OAuthとかの認証を真面目にやってみたい
  - [ ] テスト
    - [ ] goleakでgoroutineの数を計測する
	- [ ] testでfirebase余り叩きに行きたくないので、モックを作成してテストしたい
	- [ ] rpcのテストからではなくmemoryとかusecaseのテストから始めるべきかもしれない

## 📖 Document
<https://i-pu.github.io/word-war/>

## 💻 Environment
### init
```
./protoc-gen.sh
```

### develop setup
```
docker-compose up --build
```

### server
```
cd server && go run main.go
```

### client
```
cd client && yarn && yarn start
```

## ❓ 参考文献
<https://github.com/improbable-eng/grpc-web/tree/master/go/grpcweb>

<https://techblog.ap-com.co.jp/entry/2019/07/31/165309>
