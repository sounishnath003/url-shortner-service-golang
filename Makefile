AppName=url_shortner_service

build:
	ls -GFlash
	go build -ldflags "-s -w" -o bin/$(AppName) cmd/*.go

run: build
	./bin/$(AppName)

docker-build:
	docker images
	docker rmi -f $(AppName)
	docker build -t $(AppName) -f Dockerfile .

docker-run:
	docker images
	docker run -ti $(AppName)