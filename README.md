# sudory

## Make Example

source build
```
$ make go-build target=server
```

docker login
```
$ make docker-login register=p8s.me user=blah
```

image build(server / client)  
```
$ make docker-build image=p8s.me/nexclipper/sudory target=server version=0.1.0

or

$ make docker-build image=p8s.me/nexclipper/sudory target=client version=0.1.0
```

image push
```
$ make docker-push image=p8s.me/nexclipper/sudory target=server version=0.1.0
```