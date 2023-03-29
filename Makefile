PACKAGE=github.com/NexClipper/sudory/pkg
VERSION=$(shell sed -n 's/VERSION=//p' properties.${target})
COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS=-X $(PACKAGE)/version.Version=$(VERSION) -X $(PACKAGE)/version.Commit=$(COMMIT) -X $(PACKAGE)/version.BuildDate=$(BUILD_DATE)

prep:
	go install github.com/swaggo/swag/cmd/swag@v1.8.7

swagger:
	cd pkg/server/route;go generate

docker-login:
	docker login ${register} -u ${user}

go-build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/${target}/sudory-${target} ./cmd/${target}

docker-build:
	docker build -t ${image}-${target}:latest -f Dockerfile.${target} .
	docker tag ${image}-${target}:latest ${image}-${target}:$(VERSION)

docker-push:
	docker push ${image}-${target}:$(VERSION)
	docker push ${image}-${target}:latest

clean:
	rm ./bin/server/sudory-server
	rm ./bin/client/sudory-client
