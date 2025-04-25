test testname:
    go run test/main.go {{testname}}

build-linux-x86:
    GOOS="linux" \
    GOARCH="386" \
    go build -o local/test-linux-x86 test/main.go

build-linux-x64:
    GOOS="linux" \
    GOARCH="amd64" \
    go build -o local/test-linux-x64 test/main.go

build-linux-arm64:
    GOOS="linux" \
    GOARCH="arm64" \
    go build -o local/test-linux-arm64 test/main.go

build-linux-arm:
    GOOS="linux" \
    GOARCH="arm" \
    go build -o local/test-linux-arm test/main.go

build-windows-x64:
    GOOS="windows" \
    GOARCH="amd64" \
    go build -o local/test-windows-x64.exe test/main.go

build-windows-arm64:
    GOOS="windows" \
    GOARCH="arm64" \
    go build -o local/test-windows-arm64.exe test/main.go

build: build-linux-x86 \
    build-linux-x64 \
    build-linux-arm64 \
    build-linux-arm \
    build-windows-x64 \
    build-windows-arm64
    ls -l local
