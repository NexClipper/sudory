server:
	go build -o ./bin/server/sudory-server ./cmd/server

client:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/client/sudory-client ./cmd/client

swagger:
	cd pkg/server/route;go generate

prep:
	go install github.com/swaggo/swag/cmd/swag@latest

docker:
ifeq ($(target),server)
	docker build -t p8s.me/nexclipper/sudory-server:$(version) -f Dockerfile.server .
else
	docker build -t p8s.me/nexclipper/sudory-client:$(version) -f Dockerfile.client .
endif

clean:
	rm ./bin/server
	rm ./bin/client
