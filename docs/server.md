# サーバーサイド

## mock test
[github gomock](https://github.com/golang/mock)

`go get github.com/golang/mock/mockgen`

create mock for counter repo interface:
```bash
mockgen -source domain/repository/counter.go -destination domain/repository/mock/counter.go
```

## フォーマットとコードに異常がないかの確認
```bash
gofmt -w .
go vet
```

## go test
### -race
並列処理に関する問題を色々勝手に見てくれるみたい
### -covermode=atomic
複数のパッケージにまたがってるテストのカバレッジを1つにまとめてくれるみたい


## TIPS

### `Clean Architecture` っぽい何か
[これ](https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1)
を参考にそれっぽく作りました。

#### ディレクトリ構成
```bash
$ tree -L 2
.
├── domain
│   ├── entity
│   ├── repository
│   └── service
├── interface
│   ├── memory
│   └── rpc
├── usecase
└── main.go
```
`domain` が一番内側の部分になります。

`interface/memory` は `domain/repository` の実装のつもりです。

`interface/rpc` にgRPCの `rpc` を記述します。
`rpc` は `usecase` で定義してある関数を呼ぶことになります。

grpcにおけるサービスとCAにおけるサービスとで2つサービスがあるので文脈に注意です。

いまの構成は不完全な部分があるので少しずつ直すことになります。
