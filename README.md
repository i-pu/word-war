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
![elm](https://img.shields.io/badge/-Elm-3DBB3D.svg?logo=&style=flat-square)
![js](https://img.shields.io/badge/-Javascript-FAEB7F.svg?logo=javascript&style=flat-square)
![bulma](https://img.shields.io/badge/-Bulma-FCEA0.svg?logo=&style=flat-square)

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

## 🚧 Roadmap
> want to practice scrum development

- v1.0 : minimum viable product
- v1.1 : siritori
- v1.2 : rating, ホーム画面からserverのHealth状態を知りたい。userが見るためのHealth Checkをしたい。
- v1.3 : マッチングアルゴリズムの作成。待機画面の作成。
- v1.4 : 爽快なエフェクトで表示(UI強化), twitterでツイートをさせる, tbd ...
  - しりとりの表示の仕方。
  - userがだれがいるのか
- v9.9 : マイページの作成。レーティングの履歴をAtCoder風に。k8s上げたい。

## ❓ 参考文献
<https://github.com/improbable-eng/grpc-web/tree/master/go/grpcweb>

<https://techblog.ap-com.co.jp/entry/2019/07/31/165309>
