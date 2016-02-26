
ALL: bin/fomo

bin/fomo: $(shell find src)
	GOPATH=`pwd` go build -o $@ fomo

clean:
	rm -f bin/fomo

subl:
	GOPATH=`pwd` subl .

atom:
	GOPATH=`pwd` atom .
