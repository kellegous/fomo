
ALL: bin/fomo

src/golang.org/x/crypto/ssh:
	GOPATH=`pwd` go get golang.org/x/crypto/ssh

bin/fomo: src/golang.org/x/crypto/ssh $(shell find src)
	GOPATH=`pwd` go build -o $@ fomo

clean:
	rm -f bin/fomo

subl:
	GOPATH=`pwd` subl .

atom:
	GOPATH=`pwd` atom .
