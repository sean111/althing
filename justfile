build:
    CGO_ENABLED=0 go build -o dist/althing -ldflags="-s -w" -trimpath main.go
test:
    go test ./...