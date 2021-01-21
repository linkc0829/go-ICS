build_linux_amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o release/linux/amd64/icsharing

build_linux_i386:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o release/linux/i386/icsharing

test:
	go test -v .