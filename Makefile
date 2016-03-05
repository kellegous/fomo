
ALL: bin/fomo

src/golang.org/x/crypto/ssh:
	GOPATH=`pwd` go get golang.org/x/crypto/ssh

src/github.com/BurntSushi/toml:
	GOPATH=`pwd` go get github.com/BurntSushi/toml

bin/go-bindata:
	GOPATH=`pwd` go get github.com/jteeuwen/go-bindata/...

src/fomo/local/bindata.go:bin/go-bindata $(shell find src/fomo/local/data)
	$< -o $@ -pkg local -prefix src/fomo/local/data -nomemcopy src/fomo/local/data/...

src/fomo/remote/bindata.go:bin/go-bindata $(shell find src/fomo/remote/data)
	$< -o $@ -pkg remote -prefix src/fomo/remote/data -nomemcopy src/fomo/remote/data/...

bin/fomo: src/golang.org/x/crypto/ssh src/github.com/BurntSushi/toml src/fomo/local/bindata.go src/fomo/remote/bindata.go $(shell find src -name '*.go')
	GOPATH=`pwd` go build -o $@ fomo

test:
	GOPATH=`pwd` go test fomo/...
clean:
	rm -f bin/fomo

subl:
	GOPATH=`pwd` subl .

atom:
	GOPATH=`pwd` atom .
