lindyn: main.go linode myip
	CGO_ENABLED=0 go build -ldflags="-s -w"

lindyn-mips64: main.go linode myip
	GOOS=linux GOARCH=mips64 CGO_ENABLED=0 go build -ldflags="-s -w" -o lindyn-mips64

clean:
	rm lindyn lindyn-mips
