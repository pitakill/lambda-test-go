.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/*.go
	cp -v ./assets/* bin/

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose --aws-profile yalo

remove:
	sls remove --verbose --aws-profile yalo

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
