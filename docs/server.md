# サーバーサイド

## 決め事

### Redis `Pub/Sub`
[`redigo`](https://github.com/gomodule/redigo)を使います。

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
