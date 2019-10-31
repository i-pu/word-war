# Server

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

