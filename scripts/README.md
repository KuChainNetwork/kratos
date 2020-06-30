# Kuchain Scripts

## 1. Gen proto

install protoc-gen-gocosmos:

```bash
git clone https://github.com/regen-network/cosmos-proto.git
cd cosmos-proto
go mod vendor
go install ./protoc-gen-gocosmos
```

the tool will in `$GOPATH/bin`