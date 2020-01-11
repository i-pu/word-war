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
  - [ ] ジョブワーカーみたいなのでルーム毎に
    - [ ] redisの初期化といらないキーの掃除のワーカの作成
  - [ ] k8s上げたい
  - [ ] マッチングシステムの実装
- v1.4
  - [ ] twitterでツイートをさせる(OGP)
  - [ ] Game画面ブラッシュアップ
- v1.5
  - [ ] マイページの作成(履歴)
  - [ ] 爽快なエフェクトで表示(UI強化)
- v?.?
  - [ ] リファクタリング大会
  - [ ] ちゃんとAPI
  - [ ] テスト
  - [ ] 監視、メトリクス

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
