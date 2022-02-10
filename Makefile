
prep:
	go install github.com/swaggo/swag/cmd/swag@v1.7.8

swagger:
	cd pkg/server/route;go generate

docker-login:
	docker login ${register} -u ${user}

go-build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${target}/sudory-${target} ./cmd/${target}

docker-build:
	docker build -t ${image}-${target}:$(version) -f Dockerfile.${target} .

docker-push:
	docker push ${image}-${target}:${version}

clean:
	rm ./bin/server
	rm ./bin/client
